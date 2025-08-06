package route

import (
	"sheeptube/internal/app/handler"

	"github.com/gin-gonic/gin"
)

const (
	version1 = "/v1"

	videoPrefix    = "/video"
	videos         = "/"
	video          = "/:id"
	newVideo       = "/upload"
	updateMetadata = "/metadata/:id"
)

func SetupRouter(r *gin.Engine, h *handler.Handler) {
	v1 := r.Group(version1)
	{
		videoGroup := v1.Group(videoPrefix)
		{
			videoGroup.GET(videos, h.GetAllVideo)                 // /v1/video/
			videoGroup.GET(video, h.GetVideo)                     // /v1/video/{id}
			videoGroup.POST(newVideo, h.NewVideo)                 // /v1/video/upload
			videoGroup.PUT(updateMetadata, h.UpdateVideoMetadata) // /v1/video/metadata/{id}
			videoGroup.GET("test", h.Test)
		}
	}
}
