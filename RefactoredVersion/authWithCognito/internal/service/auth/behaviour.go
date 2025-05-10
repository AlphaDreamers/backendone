package auth

import (
	"github.com/SwanHtetAungPhyo/authCognito/internal/model"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"go.uber.org/fx"
)

var ServiceLayerModule = fx.Module("service_layer_module",
	fx.Provide(
		NewAuthConcrete,
	),
)

type Behaviour interface {
	SignUp(req *model.UserSignUpRequest) (err error)
	SignIn(req *model.UserSignInReq) (data *model.UserData, input *cognitoidentityprovider.InitiateAuthOutput, err error)
	Confirm(req *model.EmailVerificationRequest) (err error)
	ResetPasswordConfirm(email, confirmationCode, newPassword string) (err error)
	Logout(accessToken string) error
}
