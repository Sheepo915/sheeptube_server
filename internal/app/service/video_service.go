package service

import (
	"context"
	video_dto "sheeptube/internal/app/dto/video"
	"sheeptube/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

type VideoService struct {
	queries *db.Queries
}

func NewVideoService(queries *db.Queries) *VideoService {
	return &VideoService{
		queries: queries,
	}
}

func (vs *VideoService) GetAllVideosForHome(ctx context.Context) ([]video_dto.GetVideoResponse, error) {
	data, err := vs.queries.GetVideosForHome(ctx, db.GetVideosForHomeParams{
		Limit:  20,
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}

	var response []video_dto.GetVideoResponse
	for _, v := range data {
		response = append(response, video_dto.GetVideoResponse{
			Name:   v.Name,
			Source: v.Source,
			Likes:  int32(v.Likes.Int64),
			PostedBy: video_dto.PostedByData{
				ChannelID: v.ChannelID.String(),
				Name:      v.Name_2,
				Pic:       v.Pic.String,
			},
			CreatedAt: v.CreatedAt.Time,
		})
	}

	return response, nil
}

func (vs *VideoService) GetVideoByID(ctx context.Context, request video_dto.GetViewRequestDTO) (*video_dto.GetVideoResponse, error) {
	id := pgtype.UUID{}
	if err := id.Scan(request.VideoID); err != nil {
		return nil, err
	}

	data, err := vs.queries.GetVideoByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &video_dto.GetVideoResponse{
		Name:   data.Name,
		Source: data.Source,
		Likes:  int32(data.Likes.Int64),
		Views:  int32(data.Views.Int64),
		Shares: int32(data.Shares.Int64),
		PostedBy: video_dto.PostedByData{
			ChannelID: data.ChannelID.String(),
			Name:      data.ChannelName,
			Pic:       data.ChannelPic.String,
		},
		Categories: data.Categories,
		Tag:        data.Tags,
		Actors:     data.Actors,
	}, nil
}
