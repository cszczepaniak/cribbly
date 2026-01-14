package games

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cszczepaniak/cribbly/internal/notifier"
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
	n  *notifier.Notifier
}

func NewService(db *sql.DB, n *notifier.Notifier) Service {
	return Service{
		db: db,
		b:  sqlbuilder.New(formatter.Sqlite{}),
		n:  n,
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
	if err != nil {
		return err
	}

	s.n.Notify()
	return nil
}

func (s Service) UpdateScores(
	ctx context.Context,
	gameID string,
	team1ID string,
	team1Score int,
	team2ID string,
	team2Score int,
) (finalErr error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil || finalErr != nil {
			_ = tx.Rollback()
		}
	}()

	mustUpdateOneRow := func(res sql.Result, err error) error {
		n, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if n != 1 {
			return errors.New("unknown team ID")
		}
		return nil
	}

	err = mustUpdateOneRow(tx.ExecContext(
		ctx,
		"UPDATE Scores SET Score = ? WHERE GameID = ? AND TeamID = ?",
		team1Score,
		gameID,
		team1ID,
	))
	if err != nil {
		return err
	}

	err = mustUpdateOneRow(tx.ExecContext(
		ctx,
		"UPDATE Scores SET Score = ? WHERE GameID = ? AND TeamID = ?",
		team2Score,
		gameID,
		team2ID,
	))
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	s.n.Notify()
	return nil
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

func (s Service) Get(ctx context.Context, id string) ([2]Score, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT GameID, TeamID, Score FROM Scores WHERE GameID = ? ORDER BY TeamID`,
		id,
	)
	if err != nil {
		return [2]Score{}, err
	}
	defer rows.Close()

	var res [2]Score
	idx := 0
	for rows.Next() {
		var gameID, teamID string
		var score int
		err := rows.Scan(&gameID, &teamID, &score)
		if err != nil {
			return [2]Score{}, err
		}

		if idx == 2 {
			return [2]Score{}, errors.New("too many scores for game")
		}

		res[idx] = Score{
			GameID: gameID,
			TeamID: teamID,
			Score:  score,
		}
		idx++
	}

	return res, nil
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

type Standing struct {
	TeamID     string
	TeamName   string
	Wins       int
	Losses     int
	TotalScore int
}

func (s Service) GetStandings(ctx context.Context) ([]Standing, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			TeamID,
			t.Name,
			SUM(s.Score) as totalScore,
			SUM(s.Score >= 121) as wins,
			SUM(s.Score > 0 AND s.Score < 121) as losses
		FROM Scores s INNER JOIN Teams t ON s.TeamID = t.ID
		GROUP BY TeamID ORDER BY wins DESC, losses ASC, totalScore DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Standing
	for rows.Next() {
		var s Standing
		err := rows.Scan(&s.TeamID, &s.TeamName, &s.TotalScore, &s.Wins, &s.Losses)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}

	return res, nil
}
