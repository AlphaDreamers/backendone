package user

import (
	"context"
	"fmt"
	"github.com/SwanHtetAungPhyo/authCognito/internal/repo/user"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"mime/multipart"
)

var UserServiceModule = fx.Module("user-service", fx.Provide(NewUserService))

type UserService struct {
	log      *logrus.Logger
	repo     *user.UserRepository
	s3Client *s3.Client
}

func NewUserService(log *logrus.Logger, repo *user.UserRepository, s3Client *s3.Client) *UserService {
	return &UserService{log: log, repo: repo, s3Client: s3Client}

}

func (us *UserService) UpdateAvatar(cognito_username string, multipartFile multipart.FileHeader) (string, error) {
	openFile, err := multipartFile.Open()
	if err != nil {
		return "", err
	}
	defer func(openFile multipart.File) {
		err := openFile.Close()
		if err != nil {
			us.log.Errorf("failed to close multipart file: %v", err)
			return
		}
	}(openFile)

	bucketName := "wolftagon-swan-htet"
	key := bucketName + "/" + cognito_username
	_, err = us.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(cognito_username),
		Body:        openFile,
		ContentType: aws.String(multipartFile.Header.Get("Content-Type")),
		//ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", err
	}
	urlAva := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, key)

	if err := us.repo.UpdateAvatar(cognito_username, urlAva); err != nil {
		return "", err
	}
	return urlAva, nil
}
