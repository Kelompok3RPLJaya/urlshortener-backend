package service

import (
	"context"
	"errors"
	"url-shortener-backend/dto"
	"url-shortener-backend/entity"
	"url-shortener-backend/repository"

	"github.com/google/uuid"
	"github.com/mashingan/smapping"
)

type UrlShortenerService interface {
	CreateUrlShortener(ctx context.Context, urlShortenerDTO dto.UrlShortenerCreateDTO) (entity.UrlShortener, error)
	ValidateUrlShortenerUser(ctx context.Context, userID string, urlShortenerID string) bool
	ValidateShortUrl(ctx context.Context, urlShortenerID string) (entity.UrlShortener, error)
	GetUrlShortenerByUserID(ctx context.Context, UserID string, search string, pagination entity.Pagination) (dto.PaginationResponse, error)
	GetAllUrlShortener(ctx context.Context) ([]entity.UrlShortener, error)
	UpdateUrlShortener(ctx context.Context, urlShortenerDTO dto.UrlShortenerUpdateDTO, urlShortenerID string) error
	UpdatePrivate(ctx context.Context, urlShortenerID string, privateDTO dto.PrivateUpdateDTO) error
	UpdatePublic(ctx context.Context, urlShortenerID string) error
	GetUrlShortenerByID(ctx context.Context, urlShortenerID string) (entity.UrlShortener, error)
	DeleteUrlShortener(ctx context.Context, urlShortenerID string) error
	GetUrlShortenerByShortUrl(ctx context.Context, shortUrl string, private dto.PrivateUpdateDTO) (entity.UrlShortener, error)
}

type urlShortenerService struct {
	urlShortenerRepository repository.UrlShortenerRepository
	privateRepository      repository.PrivateRepository
	userRepository         repository.UserRepository
}

func NewUrlShortenerService(ur repository.UrlShortenerRepository, pr repository.PrivateRepository, usr repository.UserRepository) UrlShortenerService {
	return &urlShortenerService{
		urlShortenerRepository: ur,
		privateRepository:      pr,
		userRepository:         usr,
	}
}

func (us *urlShortenerService) CreateUrlShortener(ctx context.Context, urlShortenerDTO dto.UrlShortenerCreateDTO) (entity.UrlShortener, error) {
	urlShortener := entity.UrlShortener{}
	err := smapping.FillStruct(&urlShortener, smapping.MapFields(urlShortenerDTO))
	if err != nil {
		return urlShortener, err
	}
	if *urlShortener.UserID == uuid.Nil {
		urlShortener.UserID = nil
	}
	res, err := us.urlShortenerRepository.CreateUrlShortener(ctx, urlShortener)
	if err != nil {
		return urlShortener, err
	}
	if *urlShortenerDTO.IsPrivate {
		private := entity.Private{
			Password:       urlShortenerDTO.Password,
			UrlShortenerID: res.ID,
		}
		_, err = us.privateRepository.CreatePrivate(ctx, private)
		if err != nil {
			return urlShortener, err
		}
	}
	return res, err
}

func (us *urlShortenerService) ValidateUrlShortenerUser(ctx context.Context, userID string, urlShortenerID string) bool {
	urlShortenerUUID, err := uuid.Parse(urlShortenerID)
	if err != nil {
		return false
	}
	urlShortener, err := us.urlShortenerRepository.GetUrlShortenerByID(ctx, urlShortenerUUID)
	if err != nil {
		return false
	}
	if userID == urlShortener.UserID.String() {
		return true
	}
	return false
}

func (us *urlShortenerService) ValidateShortUrl(ctx context.Context, urlShortenerID string) (entity.UrlShortener, error) {
	return us.urlShortenerRepository.GetUrlShortenerByShortUrl(ctx, urlShortenerID)
}

func (us *urlShortenerService) GetUrlShortenerByUserID(ctx context.Context, UserID string, search string, pagination entity.Pagination) (dto.PaginationResponse, error) {
	userUUID, err := uuid.Parse(UserID)
	if err != nil {
		return dto.PaginationResponse{}, err
	}
	if search != "" {
		resPagination, resUrlShortener, err := us.urlShortenerRepository.GetUrlShortenerByUserIDWithSearch(ctx, userUUID, search, pagination)
		if err != nil {
			return dto.PaginationResponse{}, err
		}	
		resUser, err := us.userRepository.FindUserByID(ctx, userUUID)
		if err != nil {
			return dto.PaginationResponse{}, err
		}
		var userDTOResponse = []dto.UrlShortenerResponseDTO{}
		var userDTO = dto.UrlShortenerResponseDTO{}
		for _, v := range resUrlShortener {
			userDTO.ID = v.ID
			userDTO.Title = v.Title
			userDTO.LongUrl = v.LongUrl
			userDTO.ShortUrl = v.ShortUrl
			userDTO.Views = v.Views
			userDTO.IsPrivate = v.IsPrivate
			userDTO.IsFeeds = v.IsFeeds
			userDTO.UserID = *v.UserID
			userDTO.Username = resUser.Name
			userDTO.CreatedAt = v.CreatedAt
			userDTO.UpdatedAt = v.UpdatedAt
			userDTO.DeletedAt = v.DeletedAt
			userDTOResponse = append(userDTOResponse, userDTO)
		}
		resPagination.DataPerPage = userDTOResponse
		return resPagination, err
	} else {
		resPagination, resUrlShortener, err := us.urlShortenerRepository.GetUrlShortenerByUserID(ctx, userUUID, pagination)
		if err != nil {
			return dto.PaginationResponse{}, err
		}
		resUser, err := us.userRepository.FindUserByID(ctx, userUUID)
		if err != nil {
			return dto.PaginationResponse{}, err
		}
		var userDTOResponse = []dto.UrlShortenerResponseDTO{}
		var userDTO = dto.UrlShortenerResponseDTO{}
		for _, v := range resUrlShortener {
			userDTO.ID = v.ID
			userDTO.Title = v.Title
			userDTO.LongUrl = v.LongUrl
			userDTO.ShortUrl = v.ShortUrl
			userDTO.Views = v.Views
			userDTO.IsPrivate = v.IsPrivate
			userDTO.IsFeeds = v.IsFeeds
			userDTO.UserID = *v.UserID
			userDTO.Username = resUser.Name
			userDTO.CreatedAt = v.CreatedAt
			userDTO.UpdatedAt = v.UpdatedAt
			userDTO.DeletedAt = v.DeletedAt
			userDTOResponse = append(userDTOResponse, userDTO)
		}
		resPagination.DataPerPage = userDTOResponse
		return resPagination, err
	}
}

func (us *urlShortenerService) GetAllUrlShortener(ctx context.Context) ([]entity.UrlShortener, error) {
	return us.urlShortenerRepository.GetAllUrlShortener(ctx)
}

func (us *urlShortenerService) UpdateUrlShortener(ctx context.Context, urlShortenerDTO dto.UrlShortenerUpdateDTO, urlShortenerID string) error {
	urlShortener := entity.UrlShortener{}
	err := smapping.FillStruct(&urlShortener, smapping.MapFields(urlShortenerDTO))
	if err != nil {
		return err
	}
	urlShortenerUUID, err := uuid.Parse(urlShortenerID)
	if err != nil {
		return err
	}
	urlShortener.ID = urlShortenerUUID
	return us.urlShortenerRepository.UpdateUrlShortener(ctx, urlShortener)
}

func (us *urlShortenerService) UpdatePrivate(ctx context.Context, urlShortenerID string, privateDTO dto.PrivateUpdateDTO) error {
	urlShortenerUUID, err := uuid.Parse(urlShortenerID)
	if err != nil {
		return err
	}
	urlShortener := entity.UrlShortener{
		ID:        urlShortenerUUID,
		IsPrivate: dto.BoolPointer(true),
	}
	err = us.urlShortenerRepository.UpdateUrlShortener(ctx, urlShortener)
	if err != nil {
		return err
	}
	private := entity.Private{
		ID:             uuid.New(),
		Password:       privateDTO.Password,
		UrlShortenerID: urlShortenerUUID,
	}
	_, err = us.privateRepository.CreatePrivate(ctx, private)
	if err != nil {
		return err
	}
	return err
}

func (us *urlShortenerService) UpdatePublic(ctx context.Context, urlShortenerID string) error {
	urlShortenerUUID, err := uuid.Parse(urlShortenerID)
	if err != nil {
		return err
	}
	urlShortener := entity.UrlShortener{
		ID:        urlShortenerUUID,
		IsPrivate: dto.BoolPointer(false),
	}
	err = us.urlShortenerRepository.UpdateUrlShortener(ctx, urlShortener)
	if err != nil {
		return err
	}
	private, err := us.privateRepository.GetPrivateByUrlShortenerID(ctx, urlShortenerUUID)
	if err != nil {
		return err
	}
	return us.privateRepository.DeletePrivate(ctx, private.ID)
}

func (us *urlShortenerService) GetUrlShortenerByID(ctx context.Context, urlShortenerID string) (entity.UrlShortener, error) {
	urlShortenerUUID, err := uuid.Parse(urlShortenerID)
	if err != nil {
		return entity.UrlShortener{}, err
	}
	return us.urlShortenerRepository.GetUrlShortenerByID(ctx, urlShortenerUUID)
}

func (us *urlShortenerService) DeleteUrlShortener(ctx context.Context, urlShortenerID string) error {
	urlShortenerUUID, err := uuid.Parse(urlShortenerID)
	if err != nil {
		return err
	}

	return us.urlShortenerRepository.DeleteUrlShortener(ctx, urlShortenerUUID)
}

func(us *urlShortenerService) GetUrlShortenerByShortUrl(ctx context.Context, shortUrl string, private dto.PrivateUpdateDTO) (entity.UrlShortener, error) {
	var urlShortenerPrivate = entity.UrlShortener{}
	res, err := us.urlShortenerRepository.GetUrlShortenerByShortUrl(ctx, shortUrl)
	if err != nil {
		return entity.UrlShortener{}, err
	}
	if *res.IsPrivate {
		resPrivate, err := us.privateRepository.GetPrivateByUrlShortenerID(ctx, res.ID)
		if err != nil {
			return entity.UrlShortener{}, err
		}
		if resPrivate.Password != private.Password {
			urlShortenerPrivate.IsPrivate = dto.BoolPointer(true)
			return urlShortenerPrivate, errors.New("Password Url Shortener Salah")
		}
	}
	return us.urlShortenerRepository.IncreaseViewsCount(ctx, res)
}