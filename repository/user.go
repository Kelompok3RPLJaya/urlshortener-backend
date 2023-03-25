package repository

import (
	"context"
	"url-shortener-backend/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, user entity.User) (entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) error
	FindUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) RegisterUser(ctx context.Context, user entity.User) (entity.User, error) {
	user.ID = uuid.New()
	tx := db.connection.Create(&user)
	if tx.Error != nil {
		return entity.User{}, tx.Error
	}
	return user, nil
}

func (db *userConnection) FindUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	tx := db.connection.Where("email = ?", email).Take(&user)
	if tx.Error != nil {
		return user, tx.Error
	}
	return user, nil
}

func (db *userConnection) UpdateUser(ctx context.Context, user entity.User) error {
	tx := db.connection.Updates(&user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db *userConnection) FindUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	var userDetail entity.User
	tx := db.connection.Where("id = ?", userID).Take(&userDetail)
	if tx.Error != nil {
		return userDetail, tx.Error
	}
	userDetail.Password = ""
	return userDetail, nil
}
