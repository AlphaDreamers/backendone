package services

import (
	"github.com/SwanHtetAungPhyo/service-one/auth/internal/model"
	"github.com/google/uuid"
	"time"
)

func RegistrationConvertor(useReq *model.RegisterRequest) *model.User {
	return &model.User{
		ID:                  uuid.New(),
		FullName:            useReq.FullName,
		Email:               useReq.Email,
		Password:            useReq.Password,
		BioMetricHash:       useReq.BioMetricHash,
		IsVerified:          false,
		WalletCreated:       false,
		WalletPublicAddress: "",
		CreatedAt:           time.Now(),
	}
}
