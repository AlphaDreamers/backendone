package server

import (
	"context"
	"github.com/SwanHtetAungPhyo/order_service/providers"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"strings"
	"sync"
)

type AppState struct {
	log        *logrus.Logger
	app        *fiber.App
	nats       *nats.Conn
	clientPool sync.Map
	orderChan  chan string
}

func NewAppState(
	log *logrus.Logger,
	app *fiber.App,
	nats *nats.Conn) *AppState {

	state := &AppState{
		log:       log,
		app:       app,
		nats:      nats,
		orderChan: make(chan string),
	}
	return state
}

func (stat *AppState) Start() error {
	var err error
	go func() {
		stat.app = providers.FiberAppProvider(&sync.Once{})
		if err = stat.app.Listen(":8089"); err != nil {
			stat.log.Fatal(err.Error())
			return
		}
	}()
	return nil
}

func (stat *AppState) wsHandler() {
	stat.app.Get("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			userId := c.Get("userId")
			if userId == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "user_id is required",
				})
			}
			c.Locals("user_id", userId)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	stat.app.Get("/ws/:userId", websocket.New(func(c *websocket.Conn) {
		userId := c.Params("userId")
		if userId == "" {
			_ = c.WriteMessage(websocket.TextMessage, []byte("userId is required"))
			err := c.Close()
			if err != nil {
				return
			}
			return
		}

		stat.clientPool.Store(userId, c)

		defer func() {
			err := c.Close()
			if err != nil {
				stat.log.Error(err.Error())
				return
			}
			stat.clientPool.Delete(userId)
		}()

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
		}

		_, err := stat.nats.Subscribe("orders", func(m *nats.Msg) {
			stat.orderChan <- string(m.Data)
		})
		if err != nil {
			stat.log.Error(err.Error())
		}

	}))

}
func (stat *AppState) SubscribeToOrders() {
	_, err := stat.nats.Subscribe("orders", func(m *nats.Msg) {
		parts := strings.SplitN(string(m.Data), "|", 2)
		if len(parts) != 2 {
			stat.log.Error("Invalid message format")
			return
		}
		userId, payload := parts[0], parts[1]
		stat.OrderHandler(userId, []byte(payload))
	})
	if err != nil {
		stat.log.Error("NATS subscribe error: ", err.Error())
	}
}

func (stat *AppState) OrderHandler(userId string, message []byte) {
	conn, exists := stat.clientPool.Load(userId)
	if !exists {
		//Store in redis
	}
	ws, ok := conn.(*websocket.Conn)
	if !ok {
		// Store in redis
	}
	err := ws.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return
	}
}
func (state *AppState) Stop() error {
	if err := state.app.ShutdownWithContext(context.Background()); err != nil {
		return err
	}
	return nil
}

func AppLifeCycle(lc fx.Lifecycle, state *AppState) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			state.log.Info("Starting")
			state.SubscribeToOrders()
			return state.Start()
		},
		OnStop: func(ctx context.Context) error {
			return state.Stop()
		},
	})
}
