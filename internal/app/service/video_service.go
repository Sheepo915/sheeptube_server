package service

import (
	"context"
	video_dto "sheeptube/internal/app/dto/video"
	"sheeptube/internal/db/repository"
)

type VideoService struct {
	videoRepo   repository.VideoRepository
	channelRepo repository.ChannelRepository
}

func NewVideoService(repo *repository.Repository) *VideoService {
	return &VideoService{
		videoRepo:   repo.VideoRepository,
		channelRepo: repo.ChannelRepository,
	}
}

func (vs *VideoService) GetAllVideo(ctx context.Context) ([]video_dto.GetVideoResponse, error) {
	videos, err := vs.videoRepo.GetAllVideo(ctx)
	if err != nil {
		return nil, err
	}

	var out []video_dto.GetVideoResponse
	for _, video := range videos {
		channel, err := vs.channelRepo.GetChannelNameByID(ctx, video.PostedBy)
		if err != nil {
			return nil, err
		}

		out = append(out, video_dto.GetVideoResponse{
			VideoID: video.VideoID.String(),
			Name:    video.Name,
			Source:  video.Source,
			Views:   video.Views,
			PostedBy: video_dto.PostedByData{
				ChannelID: channel.ChannelID.String(),
				Name:      channel.Name,
			},
			Likes: video.Likes,
		})
	}

	return out, err
}

func (vs *VideoService) GetVideoByID(ctx context.Context, request video_dto.GetViewRequestDTO) (*video_dto.GetVideoResponse, error) {
	video, err := vs.videoRepo.GetVideoByID(ctx, request.VideoID)
	if err != nil {
		return nil, err
	}

	channel, err := vs.channelRepo.GetChannelNameByID(ctx, video.PostedBy)
	if err != nil {
		return nil, err
	}

	out := video_dto.GetVideoResponse{
		VideoID: video.VideoID.String(),
		Name:    video.Name,
		Source:  video.Source,
		Views:   video.Views,
		PostedBy: video_dto.PostedByData{
			ChannelID: channel.ChannelID.String(),
			Name:      channel.Name,
		},
		Likes: video.Likes,
	}

	return &out, err
}
