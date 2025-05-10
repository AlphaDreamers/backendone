package auth

import (
	"context"
	"github.com/SwanHtetAungPhyo/authCognito/internal/repo/auth"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AuthConcrete struct {
	log                *logrus.Logger
	cognitoClient      *cognitoidentityprovider.Client
	ctx                context.Context
	clientId           string
	clientSecret       string
	textractClient     *textract.Client
	rekognitiionClient *rekognition.Client
	repo               *auth.AuthRepositry
	v                  *viper.Viper
}

func NewAuthConcrete(
	log *logrus.Logger,
	cognitoClient *cognitoidentityprovider.Client,
	textractClient *textract.Client,
	rekognitiionClient *rekognition.Client,
	v *viper.Viper,
	repo *auth.AuthRepositry,
) *AuthConcrete {
	return &AuthConcrete{
		log:                log,
		cognitoClient:      cognitoClient,
		ctx:                context.Background(),
		clientId:           "7qllcjjcq7p506kq88vkfiu92g",
		clientSecret:       "1ipuga7399127snjbbgletfpr25lk6hleucb5fptn6nvrefn40ri",
		textractClient:     textractClient,
		rekognitiionClient: rekognitiionClient,
		v:                  v,
		repo:               repo,
	}
}
