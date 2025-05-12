package repository

import (
	"context"
	"github.com/SwanHtetAungPhyo/chat-order/internal/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ChatRepository struct {
	log *logrus.Logger
	v   *viper.Viper
	db  *gorm.DB
}

func NewChatRepository(log *logrus.Logger, v *viper.Viper, db *gorm.DB) *ChatRepository {
	return &ChatRepository{
		log: log,
		v:   v,
		db:  db,
	}
}
func (r ChatRepository) GetAllChatRoomByUserId(ctx context.Context, userId uuid.UUID) ([]*model.ChatRoom, error) {
	var chatRooms []*model.ChatRoom

	err := r.db.
		WithContext(ctx).
		Model(&model.ChatRoom{}).
		Where("participant_one = ? OR participant_two = ?", userId, userId).
		Find(&chatRooms).Error

	if err != nil {
		r.log.WithField("user_id", userId).Errorf("Failed to fetch chat rooms: %v", err)
		return nil, err
	}

	if len(chatRooms) == 0 {
		r.log.WithField("user_id", userId).Info("No chat rooms found for user")
	}

	return chatRooms, nil
}
