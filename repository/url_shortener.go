package repository

import (
	"context"
	"url-shortener-backend/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UrlShortenerRepository interface {
	CreateUrlShortener(ctx context.Context, urlShortener entity.UrlShortener) (entity.UrlShortener, error)
	GetAllUrlShortener(ctx context.Context) ([]entity.UrlShortener, error)
	GetUrlShortenerByID(ctx context.Context, urlShortenerID uuid.UUID) (entity.UrlShortener, error)
	GetUrlShortenerByUserID(ctx context.Context, UserID uuid.UUID) ([]entity.UrlShortener, error)
	GetUrlShortenerByShortUrl(ctx context.Context, shortUrl string) (entity.UrlShortener, error)
	UpdateUrlShortener(ctx context.Context, urlShortener entity.UrlShortener) (error)
	DeleteUrlShortener(ctx context.Context, urlShortenerID uuid.UUID) (error)
	IncreaseViewsCount(ctx context.Context, urlShortener entity.UrlShortener) (entity.UrlShortener, error)
}

type urlShortenerConnection struct {
	connection *gorm.DB
}

func NewUrlShortenerRepository(db *gorm.DB) UrlShortenerRepository {
	return &urlShortenerConnection{
		connection: db,
	}
}

func(db *urlShortenerConnection) CreateUrlShortener(ctx context.Context, urlShortener entity.UrlShortener) (entity.UrlShortener, error) {
	urlShortener.ID = uuid.New()
	tx := db.connection.Create(&urlShortener)
	if tx.Error != nil {
		return entity.UrlShortener{}, tx.Error
	}
	return urlShortener, nil
}

func(db *urlShortenerConnection) GetAllUrlShortener(ctx context.Context) ([]entity.UrlShortener, error) {
	var urlShortenerList []entity.UrlShortener
	tx := db.connection.Find(&urlShortenerList)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return urlShortenerList, nil
}

func(db *urlShortenerConnection) GetUrlShortenerByID(ctx context.Context, urlShortenerID uuid.UUID) (entity.UrlShortener, error) {
	var urlShortener entity.UrlShortener
	tx := db.connection.Where("id = ?", urlShortenerID).Take(&urlShortener)
	if tx.Error != nil {
		return entity.UrlShortener{}, tx.Error
	}
	return urlShortener, nil
}

func(db *urlShortenerConnection) GetUrlShortenerByUserID(ctx context.Context, UserID uuid.UUID) ([]entity.UrlShortener, error) {
	var urlShortener []entity.UrlShortener
	tx := db.connection.Where("user_id = ?", UserID).Find(&urlShortener)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return urlShortener, nil
}

func(db *urlShortenerConnection) GetUrlShortenerByShortUrl(ctx context.Context, shortUrl string) (entity.UrlShortener, error) {
	var urlShortener entity.UrlShortener
	tx := db.connection.Where("short_url = ?", shortUrl).Take(&urlShortener)
	if tx.Error != nil {
		return entity.UrlShortener{}, tx.Error
	}
	return urlShortener, nil
}

func(db *urlShortenerConnection) UpdateUrlShortener(ctx context.Context, urlShortener entity.UrlShortener) (error) {
	tx := db.connection.Updates(&urlShortener)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func(db *urlShortenerConnection) DeleteUrlShortener(ctx context.Context, urlShortenerID uuid.UUID) (error) {
	tx := db.connection.Delete(&entity.UrlShortener{}, &urlShortenerID)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func(db *urlShortenerConnection) IncreaseViewsCount(ctx context.Context, urlShortener entity.UrlShortener) (entity.UrlShortener, error) {
	urlShortener.Views = urlShortener.Views + 1
	tx := db.connection.Updates(&urlShortener)
	if tx.Error != nil {
		return entity.UrlShortener{}, tx.Error
	}
	return urlShortener, nil
}