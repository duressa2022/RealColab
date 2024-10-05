package usecase

import (
	"context"
	"errors"
	"time"
	"working/super_task/internal/domain"
	"working/super_task/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	UserRepository *repository.UserRepository
	Timeout        time.Duration
}

func NewUserUseCase(userrepository *repository.UserRepository, timout time.Duration) *UserUseCase {
	return &UserUseCase{
		UserRepository: userrepository,
		Timeout:        timout,
	}
}

// method for regestering user into the system
func (uu *UserUseCase) RegisterUser(cxt context.Context, user *domain.UserRegistrationRequest) (*domain.UserResponse, error) {
	_, err := uu.UserRepository.GetUserByEmail(cxt, user.Email)
	if err == nil {
		return nil, errors.New("already exsting user")
	}

	userId := primitive.NewObjectID()
	var createdUser domain.UserInformation
	createdUser.UserID = userId
	createdUser.FirstName = user.FirstName
	createdUser.LastName = user.LastName
	createdUser.Email = user.Email
	createdUser.PhoneNumber = user.PhoneNumber
	createdUser.DateOfBirth = user.DateOfBirth

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	createdUser.Password = string(hashed)
	createdUser.CreatedAt = time.Now()

	return uu.UserRepository.InsertUser(cxt, &createdUser)
}
