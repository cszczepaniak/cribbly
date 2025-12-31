package games

import (
	"context"
	"database/sql"

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

func (s Service) Create(ctx context.Context, teamID1, teamID2 string) (string, error) {
	id := uuid.NewString()

	_, err := s.b.InsertIntoTable("Scores").
		Fields("GameID", "TeamID", "Score").
		Values(id, teamID1, 0).
		ExecContext(ctx, s.db)
	if err != nil {
		return "", err
	}

	_, err = s.b.InsertIntoTable("Scores").
		Fields("GameID", "TeamID", "Score").
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
