package index

import (
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(e *chi.Mux) {

	e.Get("/", templ.Handler(IndexPage("asd")).ServeHTTP)
	// e.GET("/index/your-name-is", IndexRouteYourNameIs)
	// e.POST("/index/new-global-name", IndexRouteYourNameIsGlobal)
}
