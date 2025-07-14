package route

import (
	"sheeptube/internal/app/handler"

	"github.com/gin-gonic/gin"
)

const (
	version1 = "/v1"

	videoPrefix = "/video"
	videos      = "/"
	video       = "/:id"
	newVideo    = "/"
)

func SetupRouter(r *gin.Engine, h *handler.Handler) {
	v1 := r.Group(version1)
	{
		videoGroup := v1.Group(videoPrefix)
		{
			videoGroup.GET(videos, h.GetAllVideo) // /v1/video/
			videoGroup.GET(video, h.GetVideo)     // /v1/video/{id}
			videoGroup.POST(newVideo, h.NewVideo)
		}
	}
}
