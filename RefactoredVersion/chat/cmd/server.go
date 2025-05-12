package cmd

import (
	"context"
	"fmt"
	"github.com/SwanHtetAungPhyo/chat-order/internal/handler/chat"
	"github.com/SwanHtetAungPhyo/chat-order/internal/handler/placeOrder"
	"time"

	"github.com/SwanHtetAungPhyo/chat-order/internal/handler/ws"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var AppStateModule = fx.Module("app_state_module", fx.Provide(
	NewAppState))

type AppState struct {
	log          *logrus.Logger
	app          *fiber.App
	v            *viper.Viper
	wsHandler    *ws.WSHandler
	orderHandler *placeOrder.OrderHandler
	chatHandler  *chat.ChatRestHanlder
}

func NewAppState(
	log *logrus.Logger,
	app *fiber.App,
	v *viper.Viper,
	wsH *ws.WSHandler,
	orderH *placeOrder.OrderHandler,
	chatH *chat.ChatRestHanlder,
) *AppState {
	return &AppState{log: log, app: app, v: v, wsHandler: wsH,
		orderHandler: orderH,
		chatHandler:  chatH}
}

func (a *AppState) routeSetUp() {
	a.app.Get("/ws/chat", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}, websocket.New(a.wsHandler.ChatHandle))
	orderRest := a.app.Group("/orders")
	orderRest.Post("/orders", a.orderHandler.PlaceHandler)
	orderRest.Get("/:userId", a.orderHandler.GetAllOrdersByUserId)
	orderRest.Get("/:orderId/chat/", a.chatHandler.GetChatRoomByOrderId)
	a.app.Get("/:userId/chat/", a.chatHandler.GetAllChatRoomByUserId)
	a.app.Get("/sse/seller/:sellerId", a.orderHandler.NotificationHandler)
}

func (a *AppState) Start() error {
	a.routeSetUp()
	port := a.v.GetString("fiber.port")
	return a.app.Listen(fmt.Sprintf(":%s", port))
}

func (a *AppState) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return a.app.ShutdownWithContext(ctx)
}

func RegisterLifeCycle(lc fx.Lifecycle, state *AppState) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := state.Start()
				if err != nil {
					logrus.Error(err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return state.Stop()
		},
	})
}
