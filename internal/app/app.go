package app

import (
	"sheeptube/internal/app/handler"
	"sheeptube/internal/app/route"

	"github.com/gin-gonic/gin"
)

type App struct {
	router  *gin.Engine
	handler *handler.Handler
}

func NewApp(handler *handler.Handler) *App {
	r := gin.Default()

	route.SetupRouter(r, handler)

	return &App{
		router:  r,
		handler: handler,
	}
}

func (a *App) Run() {
	a.router.Run()
}
