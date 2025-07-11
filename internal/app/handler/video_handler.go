package handler

import (
	"net/http"
	video_dto "sheeptube/internal/app/dto/video"

	"github.com/gin-gonic/gin"
)

type VideoHandler interface {
	GetAllVideo(ctx *gin.Context)
	GetVideo(ctx *gin.Context)
}

func (h *Handler) GetAllVideo(ctx *gin.Context) {
	video, err := h.service.VideoService.GetAllVideosForHome(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusFound, video)
}

func (h *Handler) GetVideo(ctx *gin.Context) {
	var request video_dto.GetViewRequestDTO

	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	video, err := h.service.VideoService.GetVideoByID(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusFound, video)
}
