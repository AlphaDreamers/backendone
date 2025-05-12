package ws

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/SwanHtetAungPhyo/chat-order/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"strings"
	"sync"
	"time"
)

var WsModule = fx.Module("ws_module", fx.Provide(NewWSHandler))

type Client struct {
	UserID     uuid.UUID
	ChatRoomID uuid.UUID
	Conn       *websocket.Conn
}

type ConnectionRequest struct {
	UserID     uuid.UUID `json:"userId"`
	ChatRoomID uuid.UUID `json:"chatRoomId"`
}

type Message struct {
	From       uuid.UUID `json:"from"`
	To         uuid.UUID `json:"to"`
	ChatRoomID uuid.UUID `json:"chat_room_id"`
	Body       string    `json:"body"`
	File       string    `json:"file,omitempty"` // Will be replaced with S3 URL after upload
}

type WSHandler struct {
	log      *logrus.Logger
	Hub      sync.Map
	dynamodb *dynamodb.Client
	redis    *redis.Client
	s3       *s3.Client
	s3Bucket string
}

func NewWSHandler(log *logrus.Logger, ddb *dynamodb.Client, rdb *redis.Client, s3Client *s3.Client) *WSHandler {
	return &WSHandler{
		log:      log,
		dynamodb: ddb,
		redis:    rdb,
		s3:       s3Client,
		s3Bucket: "wolftagon-swan-htet",
	}
}

func (wc *WSHandler) ChatHandle(conn *websocket.Conn) {
	var init ConnectionRequest
	if _, raw, err := conn.ReadMessage(); err != nil {
		wc.log.Error("init read:", err)
		err := conn.Close()
		if err != nil {
			wc.log.Error("close connection:", err.Error())
			return
		}
		return
	} else if err := json.Unmarshal(raw, &init); err != nil {
		wc.log.Error("invalid init JSON:", err)
		err := conn.Close()
		if err != nil {
			wc.log.Error("close connection:", err.Error())
			return
		}
		return
	}

	if init.ChatRoomID == uuid.Nil || init.UserID == uuid.Nil {
		err := conn.WriteJSON(map[string]string{"error": "invalid connection parameters"})
		if err != nil {
			wc.log.Error("write:", err.Error())
			return
		}
		err = conn.Close()
		if err != nil {
			wc.log.Error("close connection:", err.Error())
			return
		}
		return
	}

	// Send welcome message
	welcome := Message{
		From:       uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		To:         init.UserID,
		ChatRoomID: init.ChatRoomID,
		Body:       "Welcome to the chat!",
	}
	err := conn.WriteJSON(welcome)
	if err != nil {
		wc.log.Error("write:", err.Error())
		return
	}

	// Check and send unread messages
	if unread, err := wc.getAndClearUnreadMessages(init.ChatRoomID, init.UserID); err == nil {
		for _, msg := range unread {
			err := conn.WriteJSON(msg)
			if err != nil {
				wc.log.Error("write:", err.Error())
				return
			}
		}
	} else {
		wc.log.Error("failed to retrieve unread messages:", err)
	}

	// Register client
	key := wc.HubKey(init.UserID, init.ChatRoomID)
	client := &Client{UserID: init.UserID, ChatRoomID: init.ChatRoomID, Conn: conn}
	wc.Hub.Store(key, client)

	defer func() {
		wc.Hub.Delete(key)
		err := conn.Close()
		if err != nil {
			wc.log.Error("close connection:", err.Error())
			return
		}
	}()

	for {
		var in Message
		if err := conn.ReadJSON(&in); err != nil {
			wc.log.Error("read:", err)
			break
		}

		// Handle file upload
		if in.File != "" {
			fileURL, err := wc.uploadFileToS3(in.File, init.UserID)
			if err != nil {
				wc.log.Error("file upload failed:", err)
				err := conn.WriteJSON(map[string]string{"error": "file upload failed"})
				if err != nil {
					wc.log.Error("write:", err.Error())
					return
				}
				continue
			}
			in.File = fileURL
		}

		// Store in DynamoDB
		dynamoMsg := &model.Message{
			ChatRoomId: in.ChatRoomID,
			To:         in.To,
			From:       in.From,
			Body:       in.Body,
			ImageUrl:   in.File, // Store S3 URL
			Timestamp:  time.Now().UTC().Unix(),
		}

		if err := wc.putMessage("ChatMessages", dynamoMsg); err != nil {
			wc.log.Error("dynamodb put:", err)
			continue
		}

		in.From = init.UserID
		if in.ChatRoomID == uuid.Nil {
			in.ChatRoomID = init.ChatRoomID
		}

		// Distribute message
		if in.To == uuid.Nil {
			wc.broadcastToRoom(in, in.ChatRoomID)
		} else {
			wc.sendToUser(in, in.To, in.ChatRoomID)
		}
	}
}

func (wc *WSHandler) uploadFileToS3(base64Data string, userID uuid.UUID) (string, error) {
	parts := strings.SplitN(base64Data, ",", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid base64 data")
	}

	fileType := strings.TrimSuffix(strings.Split(parts[0], ";")[0], "data:")
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}

	fileName := fmt.Sprintf("%s/%d.%s", userID.String(), time.Now().UnixNano(), strings.TrimPrefix(fileType, "image/"))

	_, err = wc.s3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(wc.s3Bucket),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(fileType),
	})
	if err != nil {
		return "", fmt.Errorf("s3 upload failed: %w", err)
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", wc.s3Bucket, fileName), nil
}

func (wc *WSHandler) sendToUser(msg Message, userID uuid.UUID, roomID uuid.UUID) {
	key := wc.HubKey(userID, roomID)
	if client, ok := wc.Hub.Load(key); ok {
		if err := client.(*Client).Conn.WriteJSON(msg); err != nil {
			wc.log.Errorf("Error sending to user %s: %v", userID, err)
		}
	} else {
		if err := wc.storeUnreadMessage(roomID, userID, msg); err != nil {
			wc.log.Errorf("Failed to store unread message: %v", err)
		}
	}
}

// Redis operations
func (wc *WSHandler) storeUnreadMessage(chatRoomID, userID uuid.UUID, msg Message) error {
	ctx := context.Background()
	key := fmt.Sprintf("unread:%s:%s", chatRoomID.String(), userID.String())
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	return wc.redis.RPush(ctx, key, data).Err()
}

func (wc *WSHandler) getAndClearUnreadMessages(chatRoomID, userID uuid.UUID) ([]Message, error) {
	ctx := context.Background()
	key := fmt.Sprintf("unread:%s:%s", chatRoomID.String(), userID.String())

	pipe := wc.redis.TxPipeline()
	getCmd := pipe.LRange(ctx, key, 0, -1)
	pipe.Del(ctx, key)
	_, err := pipe.Exec(ctx)

	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("redis transaction failed: %w", err)
	}

	var messages []Message
	for _, data := range getCmd.Val() {
		var msg Message
		if err := json.Unmarshal([]byte(data), &msg); err != nil {
			wc.log.Warnf("Failed to unmarshal message: %v", err)
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (wc *WSHandler) broadcastToRoom(msg Message, roomID uuid.UUID) {
	wc.Hub.Range(func(k, v interface{}) bool {
		client := v.(*Client)
		if client.ChatRoomID == roomID {
			if err := client.Conn.WriteJSON(msg); err != nil {
				wc.log.Errorf("Error broadcasting message to client %s: %v", client.UserID, err)
			}
		}
		return true
	})
}

func (wc *WSHandler) HubKey(userID, chatRoomID uuid.UUID) string {
	return fmt.Sprintf("%s|%s", userID.String(), chatRoomID.String())
}

func (wc *WSHandler) putMessage(tableName string, msg *model.Message) error {
	item := map[string]types.AttributeValue{
		"chat_room_id": &types.AttributeValueMemberS{Value: msg.ChatRoomId.String()},
		"time_stamp":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", msg.Timestamp)},
		"to":           &types.AttributeValueMemberS{Value: msg.To.String()},
		"from":         &types.AttributeValueMemberS{Value: msg.From.String()},
		"body":         &types.AttributeValueMemberS{Value: msg.Body},
	}

	if msg.ImageUrl != "" {
		item["image_url"] = &types.AttributeValueMemberS{Value: msg.ImageUrl}
	}

	_, err := wc.dynamodb.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put message into DynamoDB: %w", err)
	}
	return nil
}
