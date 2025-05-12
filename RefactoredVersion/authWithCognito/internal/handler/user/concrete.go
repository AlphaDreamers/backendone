package user

import (
	"github.com/SwanHtetAungPhyo/authCognito/internal/model"
	"github.com/SwanHtetAungPhyo/authCognito/internal/service/user"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var UserHandlerModule = fx.Module("user", fx.Provide(
	NewUserHandler))

type UserHandler struct {
	log *logrus.Logger
	srv *user.UserService
}

func NewUserHandler(log *logrus.Logger, svc *user.UserService) *UserHandler {
	return &UserHandler{
		log: log,
		srv: svc,
	}
}

func (uh *UserHandler) AvatarUploadHandler(c *fiber.Ctx) error {
	//cognitoUserName := c.Locals("cognitoUserName").(string)
	//if cognitoUserName == "" {
	//	return c.Status(fiber.StatusBadRequest).JSON(model.Response{
	//		Message: "cognito user name is empty",
	//	})
	//}
	cognitoUserName := c.Params("cognito_user_name")
	if cognitoUserName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "user name is empty",
		})
	}
	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: err.Error(),
		})
	}
	resp, err := uh.srv.UpdateAvatar(cognitoUserName, *file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.Response{
		Data: resp,
	})
}
