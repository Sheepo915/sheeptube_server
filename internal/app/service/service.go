package service

import "sheeptube/internal/db/repository"

type Service struct {
	repo *repository.Repository

	VideoService *VideoService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo:         repo,
		VideoService: NewVideoService(repo),
	}
}
