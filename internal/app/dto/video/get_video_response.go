package video_dto

type PostedByData struct {
	ChannelID string `json:"channel_id"`
	Name      string `json:"name"`
}

type GetVideoResponse struct {
	VideoID  string       `json:"video_id"`
	Name     string       `json:"name"`
	Source   string       `json:"source"`
	PostedBy PostedByData `json:"posted_by"`
	Likes    int32        `json:"likes"`
	Views    int32        `json:"views"`
}
