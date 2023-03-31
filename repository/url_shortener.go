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
	UpdateUrlShortener(ctx context.Context, urlShortener entity.UrlShortener) error
	DeleteUrlShortener(ctx context.Context, urlShortenerID uuid.UUID) error
	IncreaseViewsCount(ctx context.Context, urlShortener entity.UrlShortener) (entity.UrlShortener, error)
	GetUrlShortenerByIDUnscopped(ctx context.Context, urlShortenerID uuid.UUID) (entity.UrlShortener, error)
	GetUrlShortenerByUserIDWithSearch(ctx context.Context, userUUID uuid.UUID, search string) ([]entity.UrlShortener, error)
}

type urlShortenerConnection struct {
	connection      *gorm.DB
	feedsRepository FeedsRepository
}

func NewUrlShortenerRepository(db *gorm.DB, fs FeedsRepository) UrlShortenerRepository {
	return &urlShortenerConnection{
		connection:      db,
		feedsRepository: fs,
	}
}

func (db *urlShortenerConnection) CreateUrlShortener(ctx context.Context, urlShortener entity.UrlShortener) (entity.UrlShortener, error) {
	urlShortener.ID = uuid.New()
	tx := db.connection.Create(&urlShortener)
	if tx.Error != nil {
		return entity.UrlShortener{}, tx.Error
	}
	var feeds = entity.Feeds{
		Data:           urlShortener.ShortUrl,
		Method:         "Create",
		UrlShortenerID: urlShortener.ID,
	}
	_, err := db.feedsRepository.CreateFeeds(ctx, feeds)
	if err != nil {
		return entity.UrlShortener{}, err
	}
	return urlShortener, nil
}

func (db *urlShortenerConnection) GetAllUrlShortener(ctx context.Context) ([]entity.UrlShortener, error) {
	var urlShortenerList []entity.UrlShortener
	tx := db.connection.Find(&urlShortenerList)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return urlShortenerList, nil
}

func (db *urlShortenerConnection) GetUrlShortenerByID(ctx context.Context, urlShortenerID uuid.UUID) (entity.UrlShortener, error) {
	var urlShortener entity.UrlShortener
	tx := db.connection.Where("id = ?", urlShortenerID).Take(&urlShortener)
	if tx.Error != nil {
		return entity.UrlShortener{}, tx.Error
	}
	return urlShortener, nil
}

func (db *urlShortenerConnection) GetUrlShortenerByUserID(ctx context.Context, UserID uuid.UUID) ([]entity.UrlShortener, error) {
	var urlShortener []entity.UrlShortener
	tx := db.connection.Where("user_id = ?", UserID).Find(&urlShortener)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return urlShortener, nil
}

func (db *urlShortenerConnection) GetUrlShortenerByShortUrl(ctx context.Context, shortUrl string) (entity.UrlShortener, error) {
	var urlShortener entity.UrlShortener
	tx := db.connection.Where("short_url = ?", shortUrl).Take(&urlShortener)
	if tx.Error != nil {
		return entity.UrlShortener{}, tx.Error
	}
	return urlShortener, nil
}

func (db *urlShortenerConnection) UpdateUrlShortener(ctx context.Context, urlShortener entity.UrlShortener) error {
	urlShortenerFeeds, err := db.GetUrlShortenerByID(ctx, urlShortener.ID)
	if err != nil {
		return err
	}
	tx := db.connection.Updates(&urlShortener)
	if tx.Error != nil {
		return tx.Error
	}
	data := urlShortenerFeeds.ShortUrl + "|||" + urlShortener.ShortUrl
	var feeds = entity.Feeds{
		Data:           data,
		Method:         "Update",
		UrlShortenerID: urlShortener.ID,
	}
	_, errFeeds := db.feedsRepository.CreateFeeds(ctx, feeds)
	if err != nil {
		return errFeeds
	}
	return nil
}

func (db *urlShortenerConnection) DeleteUrlShortener(ctx context.Context, urlShortenerID uuid.UUID) error {
	urlShortenerFeeds, err := db.GetUrlShortenerByID(ctx, urlShortenerID)
	if err != nil {
		return err
	}
	tx := db.connection.Delete(&entity.UrlShortener{}, &urlShortenerID)
	if tx.Error != nil {
		return tx.Error
	}
	var feeds = entity.Feeds{
		Data:           urlShortenerFeeds.ShortUrl,
		Method:         "Delete",
		UrlShortenerID: urlShortenerFeeds.ID,
	}
	_, errFeeds := db.feedsRepository.CreateFeeds(ctx, feeds)
	if err != nil {
		return errFeeds
	}
	return nil
}

func (db *urlShortenerConnection) IncreaseViewsCount(ctx context.Context, urlShortener entity.UrlShortener) (entity.UrlShortener, error) {
	urlShortener.Views = urlShortener.Views + 1
	tx := db.connection.Updates(&urlShortener)
	if tx.Error != nil {
		return entity.UrlShortener{}, tx.Error
	}
	return urlShortener, nil
}

func (db *urlShortenerConnection) GetUrlShortenerByIDUnscopped(ctx context.Context, urlShortenerID uuid.UUID) (entity.UrlShortener, error) {
	var urlShortener entity.UrlShortener
	tx := db.connection.Unscoped().Where("id = ?", urlShortenerID).Take(&urlShortener)
	if tx.Error != nil {
		return entity.UrlShortener{}, tx.Error
	}
	return urlShortener, nil
}

func (db *urlShortenerConnection) GetUrlShortenerByUserIDWithSearch(ctx context.Context, userUUID uuid.UUID, search string) ([]entity.UrlShortener, error) {
	var urlShortener []entity.UrlShortener
	tx := db.connection.Where("user_id = ? and (short_url LIKE ? or long_url LIKE ? or title LIKE ?)", userUUID, "%"+search+"%", "%"+search+"%", "%"+search+"%").Find(&urlShortener)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return urlShortener, nil
}
