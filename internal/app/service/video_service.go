package service

import (
	"context"
	"fmt"
	"mime/multipart"
	video_dto "sheeptube/internal/app/dto/video"
	"sheeptube/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
)

type VideoService struct {
	queries     *db.Queries
	minioClient *minio.Client
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

func (vs *VideoService) GetVideoByID(ctx context.Context, request video_dto.GetViewRequest) (*video_dto.GetVideoResponse, error) {
	id := pgtype.UUID{}
	if err := id.Scan(request.VideoID); err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
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

func (vs *VideoService) NewVideo(ctx context.Context, request video_dto.NewVideoRequest, file *multipart.FileHeader) error {
	desc := pgtype.Text{}
	if err := desc.Scan(request.Description); err != nil {
		return fmt.Errorf("invalid description: %w", err)
	}

	exists, err := vs.minioClient.BucketExists(ctx, "video")
	if err != nil {
		return err
	}
	if !exists {
		err := vs.minioClient.MakeBucket(ctx, "video", minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	fileData, err := file.Open()
	if err != nil {
		return err
	}

	info, err := vs.minioClient.PutObject(ctx, "video", request.Name, fileData, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("content-type"),
	})
	if err != nil {
		return err
	}

	err = vs.queries.NewVideo(ctx, db.NewVideoParams{
		Name:        request.Name,
		Description: desc,
		Source:      info.Location,
		Poster:      request.Poster,
		PostedBy:    1,
	})
	if err != nil {
		return err
	}

	return nil
}
