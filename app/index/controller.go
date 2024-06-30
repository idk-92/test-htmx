package index

import (
	"main/app/video_player"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(e *chi.Mux) {

	e.Get("/", templ.Handler(video_player.VideoPlayerPage()).ServeHTTP)
	// e.GET("/index/your-name-is", IndexRouteYourNameIs)
	// e.POST("/index/new-global-name", IndexRouteYourNameIsGlobal)
}
