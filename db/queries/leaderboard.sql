-- name: CreateLeaderboardValue :exec
INSERT INTO
  leaderboard (name, score)
VALUES
  ($1, $2);


-- name: GetLeaderboardValues :many
SELECT
  name,
  score
FROM
  leaderboard
ORDER BY
  score DESC
LIMIT
  10;
