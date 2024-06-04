package leaderboard

import (
	"context"
	"log"
	"main/app/leaderboard/components"
	"main/db"
	"main/dbgen"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func LeaderboardService(e *chi.Mux) templ.Component {

	// template.NewTemplateRenderer(c.Echo())

	// Get params
	// dateParam := chi.URLParam(r, "date")

	// title := chi.URLParam(r, "date")

	// Get the updated list of items from the database
	items, err := getLeaderboardTableData(context.Background())

	if err != nil {
		log.Println("getLeaderboardTableData err: ", err)
		return LeaderboardPage("test title", items)
	}

	// Render the component
	return LeaderboardPage("test title", items)
}

func addLeaderboardItem(w http.ResponseWriter, r *http.Request) templ.Component {
	// template.NewTemplateRenderer(c.Echo())

	// Construct the model and bind the form data to it
	// leaderboardItem := &dbgen.CreateLeaderboardValueParams{}
	// if err := c.Bind(leaderboardItem); err != nil {
	// 	return nil
	// }

	pool, err := db.Connect(r.Context())
	if err != nil {
		log.Println("LeaderboardAddService: pool connect err", err)
		return nil
	}
	defer pool.Close()

	queries := dbgen.New(pool)
	// Save to database using GORM

	r.ParseForm()
	scoreString := r.Form.Get("score")

	score, err := strconv.ParseInt(scoreString, 10, 32)
	if err != nil {
		log.Println(err)
		return nil
	}

	scoreValue := pgtype.Int4{Int32: int32(score), Valid: true}
	err = queries.CreateLeaderboardValue(r.Context(), dbgen.CreateLeaderboardValueParams{
		Name:  pgtype.Text{String: r.Form.Get("name"), Valid: r.Form.Get("name") != ""},
		Score: scoreValue,
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	// Get the updated list of items from the database
	items, err := getLeaderboardTableData(r.Context())

	if err != nil {
		log.Println(err)
		return nil
	}

	// Create the component with the leaderboard data
	component := components.LeaderboardTable(items)

	// Render the component
	return component
}
