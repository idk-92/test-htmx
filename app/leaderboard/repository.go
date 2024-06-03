package leaderboard

import (
	"context"
	"fmt"
	"main/db"
	"main/dbgen"
	"sort"
)

// Get the list of items from the database
func getLeaderboardTableData(ctx context.Context) ([]dbgen.GetLeaderboardValuesRow, error) {

	pool, err := db.Connect(ctx)

	if err != nil {
		return nil, err
	}
	defer pool.Close()
	queries := dbgen.New(pool)
	items, err := queries.GetLeaderboardValues(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetLeaderboardValues: %w", err)
	}

	// Sort the items
	sort.Slice(items, func(i, j int) bool {
		return items[i].Score.Int32 > items[j].Score.Int32
	})

	return items, nil
}
