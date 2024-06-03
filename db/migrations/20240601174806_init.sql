-- +goose Up
CREATE TABLE
  leaderboard (id SERIAL PRIMARY KEY, name VARCHAR(255), score INT);


-- +goose Down
DROP TABLE
  leaderboard;
