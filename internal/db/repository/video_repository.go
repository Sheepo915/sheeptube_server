package repository

import (
	"context"
	"sheeptube/internal/db/model"

	"github.com/jackc/pgx/v5"
)

type VideoRepository interface {
	GetAllVideo(ctx context.Context) ([]model.Video, error)
	GetVideoByID(ctx context.Context, videoID string) (*model.Video, error)
}

func (r *Repository) GetAllVideo(ctx context.Context) ([]model.Video, error) {
	tx, err := r.GetTx()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT * FROM video")
	if err != nil {
		return nil, err
	}

	videos, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Video, error) {
		var video model.Video
		err := row.Scan(&video.VideoID, &video.Name, &video.Source, &video.PostedBy, &video.Likes, &video.Views)
		return video, err
	})
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (r *Repository) GetVideoByID(ctx context.Context, videoID string) (*model.Video, error) {
	tx, err := r.GetTx()
	if err != nil {
		return nil, err
	}

	rows := tx.QueryRow(ctx, "SELECT * FROM video WHERE video_id = $1", videoID)
	if err != nil {
		return nil, err
	}

	var video model.Video
	err = rows.Scan(&video.VideoID, &video.Name, &video.Source, &video.PostedBy, &video.Likes, &video.Views)
	if err != nil {
		return nil, err
	}

	return &video, nil
}
