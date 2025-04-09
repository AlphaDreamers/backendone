package services

import (
	"errors"
	"github.com/SwanHtetAungPhyo/auth/internal/models"
	"github.com/SwanHtetAungPhyo/auth/internal/repo"
	"github.com/SwanHtetAungPhyo/common/pkg/logutil"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthService interface {
		Login(req *models.LoginRequest) error
		Register(req *models.UserRegisterRequest) error
		Me() *models.UserInDB
		PasswordHash(plainPassword string) []byte
		ComparePassword(plainPassword string, hashedPassword []byte) error
		Convertor(req *models.UserRegisterRequest) *models.UserInDB
		UpdateStatus(email string) error
	}
	Impl struct {
		log  *logrus.Logger
		repo *repo.Impl
	}
)

func NewServiceImpl() *Impl {
	return &Impl{
		log:  logutil.GetLogger(),
		repo: repo.NewImpl(),
	}
}

func (i *Impl) Login(req *models.LoginRequest) error {
	fromDB, err := i.repo.Login(req)
	if fromDB == nil || err != nil {
		return errors.New("Login error: " + err.Error())
	}
	if fromDB.Verified == false {
		return errors.New("login error: " + "User email is not verified" + "you need to verify first")
	}
	err = i.ComparePassword(req.Password, []byte(fromDB.Password))
	if err != nil {
		return errors.New("invalid password")
	}
	return nil
}
func (i *Impl) Register(req *models.UserRegisterRequest) error {
	var userToSave = new(models.UserInDB)
	userToSave.Email = req.Email
	userToSave.FullName = req.FullName
	hashPassword := i.PasswordHash(req.Password)
	if hashPassword != nil {
		userToSave.Password = string(hashPassword)
	}
	if err := i.repo.Create(userToSave, req.BioMetricHash); err != nil {
		i.log.Error(err.Error())
		return err
	}
	i.log.Info("User registered successfully")
	return nil
}
func (i *Impl) Me() *models.UserInDB {
	return nil
}
func (i *Impl) PasswordHash(password string) []byte {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil
	}
	return fromPassword
}
func (i *Impl) ComparePassword(password string, hashedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
func (i *Impl) UpdateStatus(email string) error {
	err := i.repo.UpdateUserStatus(email)
	if err != nil {
		return err
	}
	return nil
}
func (i *Impl) Convertor(req *models.UserRegisterRequest) *models.UserInDB {
	return &models.UserInDB{
		Email:    req.Email,
		FullName: req.FullName,
	}
}

func (i *Impl) GetUserByEmail(email string) (*models.UserInDB, error) {
	byEmail, err := i.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return byEmail, nil
}

func (i *Impl) UpdateWalletStatus(userid string) error {
	err := i.repo.UpdateWalletStatus(userid)
	if err != nil {
		return err
	}
	return nil
}
