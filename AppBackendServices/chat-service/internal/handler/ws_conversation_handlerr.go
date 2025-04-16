package handler

import (
	"github.com/gofiber/websocket/v2"
	"github.com/sirupsen/logrus"
	"githubc.com/SwanHtetAungPhyo/chat-service/internal/model"
	"strings"
	"sync"
)

var connectionMap = make(map[string]*model.Client)
var mutex sync.Mutex

type WebSocketConversationHandler struct {
	logger         *logrus.Logger
	registerChan   chan *model.Client
	unregisterChan chan *model.Client
	broadcastChan  chan *model.ChatMessage
}

func NewWebSocketConversationHandler(logger *logrus.Logger) *WebSocketConversationHandler {
	return &WebSocketConversationHandler{logger: logger, registerChan: make(chan *model.Client), unregisterChan: make(chan *model.Client),
		broadcastChan: make(chan *model.ChatMessage)}
}
func (w *WebSocketConversationHandler) MainHandler(c *websocket.Conn) {
	userIdRaw := c.Locals("userId")
	userId, ok := userIdRaw.(string)
	if !ok || userId == "" {
		w.logger.Warn("Invalid or missing userId in WebSocket connection")
		_ = c.WriteMessage(websocket.TextMessage, []byte("Unauthorized"))
		_ = c.Close()
		return
	}

	client := &model.Client{
		UserID: userId,
		Conn:   c,
	}

	w.registerChan <- client

	defer func() {
		w.unregisterChan <- client
	}()

	var message model.ChatMessage
	for {
		err := c.ReadJSON(&message)
		if err != nil {
			w.logger.Warn("Error reading WebSocket message: ", err)
			break
		}
		w.logger.Warn("Received message from client: ", client.UserID)
		w.broadcastChan <- &message
	}
}

func (w *WebSocketConversationHandler) StartHub() {
	go func() {
		for {
			select {
			case client := <-w.registerChan:
				mutex.Lock()
				connectionMap[client.UserID] = client
				mutex.Unlock()
				w.logger.Infof("Client %s connected", client.UserID)
				go FetchOfflineOnConnect(client.UserID, client)
			case client := <-w.unregisterChan:
				mutex.Lock()
				if _, ok := connectionMap[client.UserID]; ok {
					delete(connectionMap, client.UserID)
					err := client.Conn.Close()
					if err != nil {
						w.logger.Warn("Error closing connection: ", err)
						return
					}
				}
				mutex.Unlock()
				w.logger.Infof("Client %s disconnected", client.UserID)
			case message := <-w.broadcastChan:
				w.logger.Infof("Trying to broadcast Message: with senderID %s and receipent ID %s", message.SenderID, message.RecipientID)
				go w.Hub(message)
			}
		}
	}()
}

func (w *WebSocketConversationHandler) Hub(message *model.ChatMessage) {

	mutex.Lock()
	defer mutex.Unlock()
	w.logger.Infof("Received message from client: %s In HUB", message.SenderID)
	if _, found := connectionMap[message.RecipientID.String()]; !found {
		PublishingToQueue(message, w.logger)
		return
	}
	for _, target := range connectionMap {
		if message.SenderID.String() == target.UserID {
			continue
		}
		if strings.EqualFold(message.RecipientID.String(), target.UserID) {
			err := target.Conn.WriteJSON(message)
			if err != nil {
				w.logger.Error("Error sending message to target:", err)
				return
			}
		}
	}
}
