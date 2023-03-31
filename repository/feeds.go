package repository

import (
	"context"
	"url-shortener-backend/common"
	"url-shortener-backend/dto"
	"url-shortener-backend/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeedsRepository interface {
	CreateFeeds(ctx context.Context, feeds entity.Feeds) (entity.Feeds, error)
	GetAllFeeds(ctx context.Context, pagination entity.Pagination) (dto.PaginationResponse, []entity.Feeds, error)
	GetTotalData(ctx context.Context) (int64, error)
}

type feedsConnection struct {
	connection *gorm.DB
}

func NewFeedsRepository(db *gorm.DB) FeedsRepository {
	return &feedsConnection{
		connection: db,
	}
}

func(db *feedsConnection) CreateFeeds(ctx context.Context, feeds entity.Feeds) (entity.Feeds, error) {
	feeds.ID = uuid.New()
	tx := db.connection.Create(&feeds)
	if tx.Error != nil {
		return entity.Feeds{}, tx.Error
	}
	return feeds, nil
}

func(db *feedsConnection) GetAllFeeds(ctx context.Context, pagination entity.Pagination) (dto.PaginationResponse, []entity.Feeds, error) {
	var paginationResponse dto.PaginationResponse
	var feedsList []entity.Feeds

	totalData, _ := db.GetTotalData(ctx)

	tx := db.connection.Debug().Scopes(common.Pagination(&pagination, totalData)).Order("created_at desc").Find(&feedsList)
	if tx.Error != nil {
		return dto.PaginationResponse{}, nil, tx.Error
	}
	// paginationResponse.DataPerPage = feedsList
	paginationResponse.Meta.MaxPage = pagination.MaxPage
	paginationResponse.Meta.Page = pagination.Page
	paginationResponse.Meta.TotalData = pagination.TotalData
	return paginationResponse, feedsList, nil
}

func(db *feedsConnection) GetTotalData(ctx context.Context) (int64, error) {
	var totalData int64
	bc := db.connection.Model(&entity.Feeds{}).Count(&totalData)
	if bc.Error != nil {
		return 0, bc.Error
	}
	return totalData, nil
}