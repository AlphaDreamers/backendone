package server

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type AppState struct {
	log      *logrus.Logger
	db       *gorm.DB
	app      *fiber.App
	mongo    *mongo.Client
	v        *viper.Viper
	s3Client *s3.Bucket
}

func NewAppState(
	log *logrus.Logger,
	db *gorm.DB,
	mongo *mongo.Client,
	v *viper.Viper,

) *AppState {
	return &AppState{}
}

func (state *AppState) Start() error {

}
func (state *AppState) Stop() error {

}

func (state *AppState) routeSetup() {

}

func (state *AppState) middleWareSetup() {

}

func (state *AppState) AppLifeCycle(lc *fx.Lifecycle, server *AppState) {

}
