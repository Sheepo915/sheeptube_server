package handler

import (
	"net/http"
	video_dto "sheeptube/internal/app/dto/video"
	"sheeptube/internal/app/service"

	"github.com/gin-gonic/gin"
)

type VideoHandler interface {
	GetAllVideo(ctx *gin.Context)
	GetVideo(ctx *gin.Context)
	NewVideo(ctx *gin.Context)
}

type videoHandlerImpl struct {
	videoService *service.VideoService
}

func newVideoHandler(vs *service.VideoService) *videoHandlerImpl {
	return &videoHandlerImpl{
		videoService: vs,
	}
}

func (h *videoHandlerImpl) GetAllVideo(ctx *gin.Context) {
	video, err := h.videoService.GetAllVideosForHome(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusFound, video)
}

func (h *videoHandlerImpl) GetVideo(ctx *gin.Context) {
	var request video_dto.GetViewRequest

	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	video, err := h.videoService.GetVideoByID(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusFound, video)
}

func (h *videoHandlerImpl) NewVideo(ctx *gin.Context) {
	var request video_dto.NewVideoRequest

	err := ctx.ShouldBindBodyWithJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	file, err := ctx.FormFile("video")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.videoService.NewVideo(ctx, request, file)
}
