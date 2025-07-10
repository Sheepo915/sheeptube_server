package route

import (
	"sheeptube/internal/app/handler"

	"github.com/gin-gonic/gin"
)

const (
	version1 = "/v1"

	videoPrefix = "/video"
	videos      = "/"
	video       = "/s"
)

func SetupRouter(r *gin.Engine, h *handler.Handler) {
	v1 := r.Group(version1)
	{
		videoGroup := v1.Group(videoPrefix)
		{
			videoGroup.GET(videos, h.VideoHandler.GetAllVideo) // /v1/video/
			videoGroup.GET(video, h.VideoHandler.GetVideo)     // /v1/video/s{query}
		}
	}
}
