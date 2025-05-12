package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"strings"

	"github.com/SwanHtetAungPhyo/authCognito/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	rtype "github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	textraTyp "github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/sirupsen/logrus"
)

var _ Behaviour = (*AuthConcrete)(nil)

func (a AuthConcrete) SignUp(req *model.UserSignUpRequest) error {
	a.log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("Starting user sign up")

	userName := req.Email

	userId, err := a.repo.CheckUserExistence(req.Email)
	if userId != nil && errors.Is(err, errors.New("user exists")) {
		a.log.WithFields(logrus.Fields{
			"email": req.Email,
		}).Debug("User already exists")
		return errors.New("user already exists, try to login")
	}
	hash := a.genSecretHash(userName)
	input := &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(a.clientId),
		Username:   aws.String(req.Email),
		Password:   aws.String(req.Password),
		SecretHash: aws.String(hash),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(req.Email),
			},
			{
				Name:  aws.String("given_name"),
				Value: aws.String(req.FirstName),
			},
			{
				Name:  aws.String("family_name"),
				Value: aws.String(req.LastName),
			},
			{
				Name:  aws.String("custom:country"),
				Value: aws.String(req.Country),
			},
			{
				Name:  aws.String("custom:bio_metric_hash"),
				Value: aws.String(req.BioMetricHash),
			},
		},
	}

	up, err := a.cognitoClient.SignUp(a.ctx, input)
	if err != nil {
		a.log.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err,
		}).Error("Sign up failed")
		return fmt.Errorf("sign up failed: %w", err)
	}
	a.log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Infoln("Sign up succeeded", up)
	modelInDB := &model.User{
		ID:                uuid.New(),
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Email:             req.Email,
		CognitoUsername:   *up.UserSub,
		TwoFactorVerified: false,
		Username:          req.Username,
		Avatar:            nil,
		Country:           req.Country,
		WalletCreated:     false,
	}
	err = a.repo.SignUp(modelInDB)
	if err != nil {
		a.log.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err.Error(),
		}).Error("Sign up failed", err.Error())
		return err
	}
	userBiometric := model.Biometrics{
		Value:           req.BioMetricHash,
		CognitoUsername: up.UserSub,
		UserID:          modelInDB.ID,
	}
	err = a.repo.SaveBiometrics(userBiometric)
	if err != nil {
		a.log.WithFields(logrus.Fields{
			"email": req.Email,
		}).Info("Save biometrics failed")
	}
	a.log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("User sign up successful")
	a.log.WithField("response", up).Debug("Sign up response")
	return nil
}

func (a AuthConcrete) extractCognitoUsername(idToken string) (string, error) {
	parser := new(jwt.Parser)
	token, _, err := parser.ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if username, exists := claims["cognito:username"]; exists {
			return fmt.Sprintf("%v", username), nil
		}
		return "", fmt.Errorf("cognito:username claim not found")
	}

	return "", fmt.Errorf("invalid token claims")
}
func (a AuthConcrete) SignIn(req *model.UserSignInReq) (*model.UserData, *cognitoidentityprovider.InitiateAuthOutput, error) {
	a.log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("Starting user sign in")

	secretHash := a.genSecretHash(req.Email)
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(a.clientId),
		AuthParameters: map[string]string{
			"USERNAME":    req.Email,
			"PASSWORD":    req.Password,
			"SECRET_HASH": secretHash,
		},
	}

	resp, err := a.cognitoClient.InitiateAuth(a.ctx, input)
	if err != nil {
		a.log.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err,
		}).Error("Sign in failed")
		return nil, nil, fmt.Errorf("sign in failed: %w", err)
	}
	userdataDecode, err := a.DecodeIdToken(*resp.AuthenticationResult.IdToken)
	if err != nil {
		a.log.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err,
		}).Error("Failed to decode ID token")
		return nil, nil, fmt.Errorf("token decode failed: %w", err)
	}

	//status, err := a.repo.GetKYCVerifiedStatus(email)
	//if err != nil {
	//	return nil, nil, err
	//}

	a.log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("User sign in successful")
	return userdataDecode, resp, nil
}

func (a AuthConcrete) Confirm(req *model.EmailVerificationRequest) error {
	a.log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("Starting email confirmation")

	hash := a.genSecretHash(req.Email)
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(a.clientId),
		Username:         aws.String(req.Email),
		ConfirmationCode: aws.String(req.Code),
		SecretHash:       aws.String(hash),
	}

	_, err := a.cognitoClient.ConfirmSignUp(a.ctx, input)
	if err != nil {
		a.log.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err,
		}).Error("Email confirmation failed")
		return fmt.Errorf("confirmation failed: %w", err)
	}
	if err := a.repo.UpdateAccountVerificationStatus(req.Email); err != nil {
		a.log.WithFields(logrus.Fields{
			"email": req.Email,
		}).Error(err.Error())
		return fmt.Errorf("email verification update in local failed: %w", err)
	}
	a.log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("Email confirmed successfully")
	return nil
}

func (a AuthConcrete) ResendConfirmation(email string) error {
	a.log.WithFields(logrus.Fields{
		"email": email,
	}).Info("Resending confirmation code")

	hash := a.genSecretHash(email)
	input := &cognitoidentityprovider.ResendConfirmationCodeInput{
		Username:   aws.String(email),
		ClientId:   aws.String(a.clientId),
		SecretHash: aws.String(hash),
	}

	_, err := a.cognitoClient.ResendConfirmationCode(a.ctx, input)
	if err != nil {
		a.log.WithFields(logrus.Fields{
			"email": email,
			"error": err,
		}).Error("Failed to resend confirmation code")
		return fmt.Errorf("resend failed: %w", err)
	}

	a.log.WithFields(logrus.Fields{
		"email": email,
	}).Info("Confirmation code resent successfully")
	return nil
}

func (a AuthConcrete) ForgotPassword(email string) error {
	a.log.WithFields(logrus.Fields{
		"email": email,
	}).Info("Starting password reset process")

	hash := a.genSecretHash(email)
	input := &cognitoidentityprovider.ForgotPasswordInput{
		Username:   aws.String(email),
		ClientId:   aws.String(a.clientId),
		SecretHash: aws.String(hash),
	}

	_, err := a.cognitoClient.ForgotPassword(a.ctx, input)
	if err != nil {
		a.log.WithFields(logrus.Fields{
			"email": email,
			"error": err,
		}).Error("Password reset request failed")
		return fmt.Errorf("password reset failed: %w", err)
	}

	a.log.WithFields(logrus.Fields{
		"email": email,
	}).Info("Password reset request sent successfully")
	return nil
}

func (a AuthConcrete) ResetPasswordConfirm(email, confirmationCode, newPassword string) error {
	a.log.WithFields(logrus.Fields{
		"email": email,
	}).Info("Confirming password reset")

	hash := a.genSecretHash(email)
	input := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		Username:         aws.String(email),
		ClientId:         aws.String(a.clientId),
		ConfirmationCode: aws.String(confirmationCode),
		Password:         aws.String(newPassword),
		SecretHash:       aws.String(hash),
	}

	_, err := a.cognitoClient.ConfirmForgotPassword(a.ctx, input)
	if err != nil {
		a.log.WithFields(logrus.Fields{
			"email": email,
			"error": err,
		}).Error("Password reset confirmation failed")
		return fmt.Errorf("password reset confirmation failed: %w", err)
	}

	a.log.WithFields(logrus.Fields{
		"email": email,
	}).Info("Password reset successfully")
	return nil
}

func (a AuthConcrete) genSecretHash(username string) string {
	key := []byte("1ipuga7399127snjbbgletfpr25lk6hleucb5fptn6nvrefn40ri")
	message := []byte(username + a.clientId)
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (a AuthConcrete) DecodeIdToken(idToken string) (*model.UserData, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		a.log.Error("Invalid token format")
		return nil, errors.New("invalid token format")
	}

	b, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		a.log.WithError(err).Error("Failed to decode token payload")
		return nil, fmt.Errorf("token decode failed: %w", err)
	}

	var userData model.UserData
	if err := json.Unmarshal(b, &userData); err != nil {
		a.log.WithError(err).Error("Failed to unmarshal token payload")
		return nil, fmt.Errorf("token unmarshal failed: %w", err)
	}

	return &userData, nil
}

func (a AuthConcrete) Logout(accessToken string) error {
	if accessToken == "" {
		a.log.Error("Empty access token provided")
		return errors.New("invalid access token")
	}

	input := &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(accessToken),
	}

	_, err := a.cognitoClient.GlobalSignOut(a.ctx, input)
	if err != nil {
		a.log.WithError(err).Error("Logout failed")
		return fmt.Errorf("logout failed: %w", err)
	}

	a.log.Info("User logged out successfully")
	return nil
}

func (a AuthConcrete) KYCVerification(files []*multipart.FileHeader, email string) (bool, error) {
	if len(files) != 2 {
		a.log.WithField("file_count", len(files)).Error("Invalid number of files for KYC")
		return false, errors.New("exactly 2 files required (ID and selfie)")
	}

	// Process ID file
	idFile, err := files[0].Open()
	if err != nil {
		a.log.WithError(err).Error("Failed to open ID file")
		return false, fmt.Errorf("failed to open ID file: %w", err)
	}
	defer idFile.Close()

	idBytes, err := io.ReadAll(idFile)
	if err != nil {
		a.log.WithError(err).Error("Failed to read ID file")
		return false, fmt.Errorf("failed to read ID file: %w", err)
	}

	selfieFile, err := files[1].Open()
	if err != nil {
		a.log.WithError(err).Error("Failed to open selfie file")
		return false, fmt.Errorf("failed to open selfie file: %w", err)
	}
	defer selfieFile.Close()

	selfieBytes, err := io.ReadAll(selfieFile)
	if err != nil {
		a.log.WithError(err).Error("Failed to read selfie file")
		return false, fmt.Errorf("failed to read selfie file: %w", err)
	}

	_, err = a.TextractClientProcessor(idBytes)
	if err != nil {
		a.log.WithError(err).Error("ID analysis failed")
		return false, fmt.Errorf("ID analysis failed: %w", err)
	}

	// Detect faces in selfie
	faces, err := a.DetectFaces(selfieBytes)
	if err != nil {
		a.log.WithError(err).Error("Face detection failed")
		return false, fmt.Errorf("face detection failed: %w", err)
	}

	// Validate face quality
	if err := a.Processing(faces, &textract.AnalyzeIDInput{
		DocumentPages: []textraTyp.Document{{Bytes: idBytes}},
	}); err != nil {
		a.log.WithError(err).Error("Face validation failed")
		return false, fmt.Errorf("face validation failed: %w", err)
	}

	// Compare faces between ID and selfie
	compareResult, err := a.CompareFaces(idBytes, selfieBytes)
	if err != nil {
		a.log.WithError(err).Error("Face comparison failed")
		return false, fmt.Errorf("face comparison failed: %w", err)
	}

	if len(compareResult.FaceMatches) == 0 {
		a.log.Error("No face matches found")
		return false, errors.New("no face matches found")
	}

	similarity := compareResult.FaceMatches[0].Similarity
	a.log.WithField("similarity", *similarity).Info("Face comparison result")
	if *similarity >= 70 {
		err := a.repo.UpdateTheKYCVerificationStatus(email)
		if err != nil {
			a.log.WithError(err).Error("Failed to update KYC verification status")
		}
		err = a.repo.UpdateBioMetricsVerification(email)
		if err != nil {
			a.log.WithError(err).Error("Failed to update BIO metrics verification status")
		}
	}

	return *similarity >= 70, nil
}

func (a AuthConcrete) Processing(faces *rekognition.DetectFacesOutput, src *textract.AnalyzeIDInput) error {
	if len(faces.FaceDetails) != 1 {
		a.log.WithField("face_count", len(faces.FaceDetails)).Error("Invalid number of faces detected")
		return errors.New("exactly one face should be detected")
	}

	face := faces.FaceDetails[0]
	minConfidence := float32(90)

	if face.Confidence == nil || *face.Confidence < minConfidence {
		a.log.WithField("confidence", *face.Confidence).Error("Low face detection confidence")
		return fmt.Errorf("low face detection confidence: %.2f", *face.Confidence)
	}

	if face.Quality == nil {
		a.log.Error("Missing face quality data")
		return errors.New("missing face quality data")
	}

	if face.Quality.Brightness == nil || face.Quality.Sharpness == nil {
		a.log.Error("Incomplete face quality metrics")
		return errors.New("incomplete face quality metrics")
	}

	minBrightness := float32(50)
	minSharpness := float32(50)
	if *face.Quality.Brightness < minBrightness || *face.Quality.Sharpness < minSharpness {
		a.log.WithFields(logrus.Fields{
			"brightness": *face.Quality.Brightness,
			"sharpness":  *face.Quality.Sharpness,
		}).Error("Poor image quality")
		return fmt.Errorf("poor image quality (brightness: %.2f, sharpness: %.2f)",
			*face.Quality.Brightness, *face.Quality.Sharpness)
	}

	a.log.WithFields(logrus.Fields{
		"confidence": *face.Confidence,
		"brightness": *face.Quality.Brightness,
		"sharpness":  *face.Quality.Sharpness,
	}).Info("Face validation passed")
	return nil
}

func (a AuthConcrete) TextractClientProcessor(idFile []byte) (*textract.AnalyzeIDOutput, error) {
	input := &textract.AnalyzeIDInput{
		DocumentPages: []textraTyp.Document{
			{Bytes: idFile},
		},
	}

	result, err := a.textractClient.AnalyzeID(a.ctx, input)
	if err != nil {
		a.log.WithError(err).Error("Textract analysis failed")
		return nil, fmt.Errorf("textract analysis failed: %w", err)
	}

	a.log.WithField("result", result).Debug("Textract analysis completed")
	return result, nil
}

func (a AuthConcrete) DetectFaces(imageBytes []byte) (*rekognition.DetectFacesOutput, error) {
	input := &rekognition.DetectFacesInput{
		Image: &rtype.Image{
			Bytes: imageBytes,
		},
		Attributes: []rtype.Attribute{
			rtype.AttributeDefault,
			rtype.AttributeAll,
		},
	}

	result, err := a.rekognitiionClient.DetectFaces(a.ctx, input)
	if err != nil {
		a.log.WithError(err).Error("Face detection failed")
		return nil, fmt.Errorf("face detection failed: %w", err)
	}

	a.log.WithField("face_count", len(result.FaceDetails)).Debug("Face detection completed")
	return result, nil
}

func (a AuthConcrete) CompareFaces(src, target []byte) (*rekognition.CompareFacesOutput, error) {
	input := &rekognition.CompareFacesInput{
		SourceImage: &rtype.Image{
			Bytes: src,
		},
		TargetImage: &rtype.Image{
			Bytes: target,
		},
		SimilarityThreshold: aws.Float32(70.0),
	}

	result, err := a.rekognitiionClient.CompareFaces(a.ctx, input)
	if err != nil {
		a.log.WithError(err).Error("Face comparison failed")
		return nil, fmt.Errorf("face comparison failed: %w", err)
	}

	a.log.WithField("matches", len(result.FaceMatches)).Debug("Face comparison completed")
	return result, nil
}
func (a AuthConcrete) RefreshAccessToken(client cognitoidentityprovider.Client, refreshToken string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeRefreshTokenAuth,
		ClientId: aws.String(a.v.GetString("client_id")),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
		},
	}

	if clientSecret := a.v.GetString("client_secret"); clientSecret != "" {
		secretHash := computeSecretHash(
			a.v.GetString("client_id"),
			a.v.GetString("client_secret"),
			a.v.GetString("username"),
		)
		input.AuthParameters["SECRET_HASH"] = secretHash
	}

	result, err := client.InitiateAuth(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return result, nil
}

func computeSecretHash(clientID, clientSecret, username string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
