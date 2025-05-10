package model

type UserSignUpResp struct {
	Message string `json:"message"`
}

type UserSignInResp struct {
	Message  string          `json:"message"`
	MetaData *SignUpMetaData `json:"metaData"`
}

type SignUpMetaData struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	UserId      string `json:"userId"`
	AccessToken string `json:"accessToken"`
}
