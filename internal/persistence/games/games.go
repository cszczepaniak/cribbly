package games

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/formatter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/table"
	"github.com/google/uuid"
)

type Score struct {
	GameID string
	TeamID string
	Score  int
}

type Game struct {
	ID     string
	Scores [2]Score
}

type Service struct {
	db *sql.DB
	b  *sqlbuilder.Builder
}

func NewService(db *sql.DB) Service {
	return Service{
		db: db,
		b:  sqlbuilder.New(formatter.Sqlite{}),
	}
}

func (s Service) Init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS Scores (
			GameID VARCHAR(36),
			TeamID VARCHAR(36),
			Score SMALLINT,
			
			PRIMARY KEY (GameID, TeamID)
		)`)
	return err
}

func (s Service) Create(ctx context.Context, teamID1, teamID2 string) (_ string, err error) {
	id := uuid.NewString()

	_, err = s.b.InsertIntoTable("Scores").
		Fields("GameID", "TeamID", "Score").
		Values(id, teamID1, 0).
		Values(id, teamID2, 0).
		ExecContext(ctx, s.db)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s Service) UpdateScore(ctx context.Context, gameID, teamID string, score int) error {
	_, err := s.b.UpdateTable("Scores").SetFieldTo("Score", score).WhereAll(
		filter.Equals("GameID", gameID),
		filter.Equals("TeamID", teamID),
	).ExecContext(ctx, s.db)
	return err
}

func (s Service) GetScore(ctx context.Context, gameID, teamID string) (int, error) {
	row, err := s.b.SelectFrom(table.Named("Scores")).Columns("Score").WhereAll(
		filter.Equals("GameID", gameID),
		filter.Equals("TeamID", teamID),
	).QueryRowContext(ctx, s.db)
	if err != nil {
		return 0, err
	}

	var score int
	err = row.Scan(&score)
	if err != nil {
		return 0, err
	}

	return score, nil
}

func (s Service) GetAll(ctx context.Context) ([]Score, error) {
	rows, err := s.b.SelectFrom(table.Named("Scores")).Columns("GameID", "TeamID", "Score").QueryContext(ctx, s.db)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []Score
	for rows.Next() {
		var s Score
		err := rows.Scan(&s.GameID, &s.TeamID, &s.Score)
		if err != nil {
			return nil, err
		}

		scores = append(scores, s)
	}

	return scores, nil
}

func (s Service) DeleteAll(ctx context.Context) error {
	_, err := s.b.DeleteFromTable("Scores").ExecContext(ctx, s.db)
	return err
}

func (s Service) GetForTeam(ctx context.Context, teamID string) (map[string][2]Score, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT GameID, TeamID, Score FROM Scores WHERE GameID IN (
			SELECT GameID FROM Scores WHERE TeamID = ?
		) ORDER BY GameID, TeamID
	`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(map[string][2]Score)
	for rows.Next() {
		var gameID, teamID string
		var score int
		err := rows.Scan(&gameID, &teamID, &score)
		if err != nil {
			return nil, err
		}

		scores, ok := res[gameID]
		if !ok {
			// Doesn't exist yet; fill in the first element.
			res[gameID] = [2]Score{{
				GameID: gameID,
				TeamID: teamID,
				Score:  score,
			}}
			continue
		}

		// Sanity check that the second score hasn't already been filled in.
		if scores[1].GameID != "" {
			return nil, errors.New("more than two scores for a game")
		}

		scores[1] = Score{
			GameID: gameID,
			TeamID: teamID,
			Score:  score,
		}
		res[gameID] = scores
	}

	return res, nil
}
