package handler

import (
	"fmt"
	"net/http"
	video_dto "sheeptube/internal/app/dto/video"
	"sheeptube/internal/app/service"
	"time"

	"github.com/gin-gonic/gin"
)

type VideoHandler interface {
	GetAllVideo(ctx *gin.Context)
	GetVideo(ctx *gin.Context)
	NewVideo(ctx *gin.Context)
	UpdateVideoMetadata(ctx *gin.Context)
	Test(ctx *gin.Context)
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
	video, err := h.videoService.GetAllVideo(ctx.Request.Context())
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
	file, err := ctx.FormFile("video")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.videoService.NewVideo(ctx, file)
}

func (h *videoHandlerImpl) UpdateVideoMetadata(ctx *gin.Context) {

}

func (h *videoHandlerImpl) Test(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	disconnect := ctx.Request.Context().Done()
	t := time.NewTicker(2 * time.Second)

	const MAX = 100
	count := 0

	defer t.Stop()
	for {
		select {
		case <-disconnect:
			return
		case <-t.C:
			if count <= MAX {
				fmt.Fprintf(ctx.Writer, "data: Loading: %d%%\n\n", count)
				ctx.Writer.Flush()
				count += 10
			} else {
				return
			}
		}
	}
}
