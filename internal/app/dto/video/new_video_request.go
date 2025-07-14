package video_dto

type NewVideoRequest struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description" binding:"required"`
	Poster      string `form:"poster" binding:"required"`
	PostedBy    string `form:"posted_by" binding:"required"`
}
