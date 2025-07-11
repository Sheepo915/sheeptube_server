package video_dto

import "time"

type PostedByData struct {
	ChannelID string `json:"channel_id"`
	Name      string `json:"name"`
	Pic       string `json:"pic"`
}

type GetVideoResponse struct {
	VideoID    string       `json:"video_id,omitempty"`
	Name       string       `json:"name"`
	Source     string       `json:"source"`
	Likes      int32        `json:"likes,omitempty"`
	Views      int32        `json:"views,omitempty"`
	Shares     int32        `json:"shares,omitempty"`
	PostedBy   PostedByData `json:"posted_by"`
	Categories []string     `json:"categories"`
	Tag        []string     `json:"tags"`
	Actors     []string     `json:"actors"`
	CreatedAt  time.Time    `json:"created_at,omitempty"`
}
