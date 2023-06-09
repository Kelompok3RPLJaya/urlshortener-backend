package service

import (
	"context"
	"url-shortener-backend/dto"
	"url-shortener-backend/entity"
	"url-shortener-backend/helpers"
	"url-shortener-backend/repository"

	"github.com/google/uuid"
	"github.com/mashingan/smapping"
)

type UserService interface {
	RegisterUser(ctx context.Context, userDTO dto.UserCreateDto) (entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (entity.User, error)
	CheckUser(ctx context.Context, email string) (bool, error)
	Verify(ctx context.Context, email string, password string) (bool, error)
	UpdateUser(ctx context.Context, userUpdateDto dto.UserUpdateDto) error
	GetDetailAccount(ctx context.Context, userID uuid.UUID) (entity.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		userRepository: ur,
	}
}

func (us *userService) RegisterUser(ctx context.Context, userDTO dto.UserCreateDto) (entity.User, error) {
	user := entity.User{}
	err := smapping.FillStruct(&user, smapping.MapFields(userDTO))
	user.Role = "user"
	if err != nil {
		return user, err
	}
	return us.userRepository.RegisterUser(ctx, user)
}

func (us *userService) FindUserByEmail(ctx context.Context, email string) (entity.User, error) {
	return us.userRepository.FindUserByEmail(ctx, email)
}

func (us *userService) CheckUser(ctx context.Context, email string) (bool, error) {
	result, err := us.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if result.Email == "" {
		return false, nil
	}
	return true, nil
}

func (us *userService) Verify(ctx context.Context, email string, password string) (bool, error) {
	res, err := us.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	CheckPassword, err := helpers.CheckPassword(res.Password, []byte(password))
	if err != nil {
		return false, err
	}
	if res.Email == email && CheckPassword {
		return true, nil
	}
	return false, nil
}

func (us *userService) UpdateUser(ctx context.Context, userUpdateDto dto.UserUpdateDto) error {
	user := entity.User{}
	err := smapping.FillStruct(&user, smapping.MapFields(userUpdateDto))
	if err != nil {
		return err
	}
	return us.userRepository.UpdateUser(ctx, user)
}

func (us *userService) GetDetailAccount(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	detailUser, err := us.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return entity.User{}, err
	}

	return detailUser, nil
}
