package leaderboard

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(e *chi.Mux) {

	// e.Get("/leaderboard", templ.Handler(LeaderboardPage("asd")).ServeHTTP)
	e.Get("/leaderboard", templ.Handler(LeaderboardService(e)).ServeHTTP)
	e.Post("/leaderboard/add", func(w http.ResponseWriter, r *http.Request) {

		templ.Handler(addLeaderboardItem(w, r)).ServeHTTP(w, r)
	})
	e.Post("/leaderboard/add", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(addLeaderboardItem(w, r)).ServeHTTP(w, r)
	})
}
