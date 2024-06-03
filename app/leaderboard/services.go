package leaderboard

import (
	"log"
	"main/app/leaderboard/components"
	"main/db"
	"main/dbgen"
	"main/template"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func LeaderboardService(c echo.Context) error {

	template.NewTemplateRenderer(c.Echo())

	// Get params
	title := c.QueryParam("title")

	// Get the updated list of items from the database
	items, err := getLeaderboardTableData(c.Request().Context())

	if err != nil {
		log.Println("getLeaderboardTableData err: ", err)
		return err
	}

	// Get the component
	component := Leaderboard(title, items)

	// Render the component
	return template.AssertRender(c, http.StatusOK, component)
}

func LeaderboardAddService(c echo.Context) error {
	template.NewTemplateRenderer(c.Echo())

	// Construct the model and bind the form data to it
	leaderboardItem := &dbgen.CreateLeaderboardValueParams{}
	if err := c.Bind(leaderboardItem); err != nil {
		return err
	}

	pool, err := db.Connect(c.Request().Context())
	if err != nil {
		log.Println("LeaderboardAddService: pool connect err", err)
		return err
	}
	defer pool.Close()

	queries := dbgen.New(pool)
	// Save to database using GORM
	scoreString := c.Request().FormValue("score")
	score, err := strconv.ParseInt(scoreString, 10, 32)
	if err != nil {
		log.Println(err)
		return err
	}
	scoreValue := pgtype.Int4{Int32: int32(score), Valid: true}
	queries.CreateLeaderboardValue(c.Request().Context(), dbgen.CreateLeaderboardValueParams{
		Name:  pgtype.Text{String: c.Request().FormValue("name"), Valid: c.Request().FormValue("name") != ""},
		Score: scoreValue,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	// Get the updated list of items from the database
	items, err := getLeaderboardTableData(c.Request().Context())

	if err != nil {
		log.Println(err)
		return err
	}

	// Create the component with the leaderboard data
	component := components.LeaderboardTable(items)

	// Render the component
	return template.AssertRender(c, http.StatusOK, component)
}
