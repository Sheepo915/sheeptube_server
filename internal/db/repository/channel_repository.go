package repository

import (
	"context"
	"sheeptube/internal/db/model"
)

type ChannelRepository interface {
	GetChannelNameByID(ctx context.Context, channelID int32) (*model.Channel, error)
}

func (r *Repository) GetChannelNameByID(ctx context.Context, channelID int32) (*model.Channel, error) {
	tx, err := r.GetTx()
	if err != nil {
		return nil, err
	}

	rows := tx.QueryRow(ctx, "SELECT * FROM channels WHERE id = $1", channelID)
	if err != nil {
		return nil, err
	}

	var channel model.Channel
	err = rows.Scan(&channel.ChannelID, &channel.Name)
	if err != nil {
		return nil, err
	}

	return &channel, nil
}
