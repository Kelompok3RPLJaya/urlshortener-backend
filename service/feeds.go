package service

import (
	"context"
	"strings"
	"url-shortener-backend/dto"
	"url-shortener-backend/entity"
	"url-shortener-backend/repository"
)

type FeedsService interface {
	GetAllFeeds(ctx context.Context, pagination entity.Pagination) (dto.PaginationResponse, error)
}

type feedsService struct {
	feedsRepository repository.FeedsRepository
	urlShortenerRepository repository.UrlShortenerRepository
	userRepository repository.UserRepository
}

func NewFeedsService(fr repository.FeedsRepository, ur repository.UrlShortenerRepository, usr repository.UserRepository) FeedsService {
	return &feedsService{
		feedsRepository: fr,
		urlShortenerRepository: ur,
		userRepository: usr,
	}
}

func(fs *feedsService) GetAllFeeds(ctx context.Context, pagination entity.Pagination) (dto.PaginationResponse, error) {
	resPagination, resFeeds, err := fs.feedsRepository.GetAllFeeds(ctx, pagination)
	if err != nil {
		return dto.PaginationResponse{}, err
	}
	var feedsDTOArray []dto.FeedsResponseDTO
	var feedsDTO dto.FeedsResponseDTO
	for _, v := range resFeeds {
		urlShortener, err := fs.urlShortenerRepository.GetUrlShortenerByIDUnscopped(ctx, v.UrlShortenerID)
		if urlShortener.ShortUrl == "" {
			continue
		}
		user, err := fs.userRepository.FindUserByID(ctx, *urlShortener.UserID)
		if err != nil {
			return dto.PaginationResponse{}, err
		}
		feedsDTO.ID = v.ID
		feedsDTO.Title = urlShortener.Title
		feedsDTO.Username = user.Name
		feedsDTO.Method = v.Method
		feedsDTO.UrlShortenerID = v.UrlShortenerID
		feedsDTO.CreatedAt = v.Timestamp.CreatedAt
		if v.Method == "Create" || v.Method == "Delete" {
			feedsDTO.Data.Before = ""
			feedsDTO.Data.After = v.Data
		} else {
			stringSplit := strings.Split(v.Data, "|||")
			feedsDTO.Data.Before = stringSplit[0]
			feedsDTO.Data.After = stringSplit[1]
		}
		feedsDTOArray = append(feedsDTOArray, feedsDTO)
	}
	resPagination.DataPerPage = feedsDTOArray
	return resPagination, nil
}