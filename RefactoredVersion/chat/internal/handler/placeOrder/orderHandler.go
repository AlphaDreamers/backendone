package placeOrder

import (
	"bufio"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/SwanHtetAungPhyo/chat-order/internal/model"
	"github.com/SwanHtetAungPhyo/chat-order/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var OrdHandlerModule = fx.Module("order_handler_module", fx.Provide(
	NewOrderHandler))

type OrderHandler struct {
	log    *logrus.Logger
	srv    *service.OrderService
	redis  *redis.Client
	gormDB *gorm.DB
}

func NewOrderHandler(log *logrus.Logger,
	redisClient *redis.Client,
	srv *service.OrderService,
	db *gorm.DB,
) *OrderHandler {
	return &OrderHandler{log: log,
		srv:    srv,
		redis:  redisClient,
		gormDB: db,
	}
}

func (o *OrderHandler) PlaceHandler(c *fiber.Ctx) error {
	var req model.OrderPlaceRequest
	if err := c.BodyParser(&req); err != nil {
		o.log.Error("Error parsing request:", err.Error())
		return c.SendStatus(fiber.StatusBadRequest)
	}

	OrderId, err := o.srv.FindOrCreate(&req)
	if err != nil {
		o.log.Error("Error in FindOrCreate:", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	masterKey := o.MasterKey(
		req.SellerId,
		req.BuyerId,
		req.ServiceId,
		OrderId.OrderNumber,
	)

	chatRoom := model.ChatRoom{
		ParticipantOne: req.BuyerId,
		ParticipantTwo: req.SellerId,
		ServiceId:      req.ServiceId,
		MasterKey:      masterKey,
		OrderId:        OrderId.ID,
	}

	if err := o.gormDB.FirstOrCreate(&chatRoom, "master_key = ?", masterKey).Error; err != nil {
		o.log.Error("Error creating chat room:", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	publish := fmt.Sprintf("New Order for %s , Status (%s), Buyer (%s)", OrderId.OrderNumber, "PENDING", req.BuyerId)
	o.log.Infof("Publishing to Redis channel '%s': %s", "order_notification:"+req.SellerId.String(), publish)

	if err := o.redis.Publish(context.TODO(), "order_notification:"+req.SellerId.String(), publish).Err(); err != nil {
		o.log.Error("Failed to publish to Redis:", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"chat_room": chatRoom.ChatRoomID,
	})
}
func (o *OrderHandler) NotificationHandler(c *fiber.Ctx) error {
	sellerId := c.Params("sellerId")
	if sellerId == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Status(fiber.StatusOK)

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		channel := "order_notification:" + sellerId
		pubSub := o.redis.Subscribe(context.Background(), channel)
		if pubSub == nil {
			o.log.Error("Redis subscription is nil")
			return
		}
		defer func() {
			if err := pubSub.Close(); err != nil {
				o.log.Error("Error closing Redis subscription:", err)
			}
		}()

		o.log.Infof("Subscribed to Redis channel: %s", channel)

		for msg := range pubSub.Channel() {
			o.log.Infof("Received message from Redis channel '%s': %s", channel, msg.Payload)

			if _, err := fmt.Fprintf(w, "data: %s\n\n", msg.Payload); err != nil {
				o.log.Error("Error writing to stream:", err)
				return
			}

			if err := w.Flush(); err != nil {
				o.log.Error("Error flushing buffer:", err)
				return
			}
		}
	}))

	return nil
}
func (o *OrderHandler) MasterKey(seller, buyer, serviceId uuid.UUID, orderId string) string {
	hash := sha256.New()
	hash.Write([]byte(seller.String() + buyer.String() + serviceId.String() + orderId))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
