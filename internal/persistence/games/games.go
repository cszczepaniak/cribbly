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

	"github.com/cszczepaniak/cribbly/internal/notifier"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
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

type Repository struct {
	db            database.Database
	b             *sqlbuilder.Builder
	scoreNotifier *notifier.Notifier
}

func NewRepository(db database.Database, scoreNotifier *notifier.Notifier) Repository {
	return Repository{
		db:            db,
		b:             sqlbuilder.New(formatter.Sqlite{}),
		scoreNotifier: scoreNotifier,
	}
}

func (s Repository) Init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS Scores (
			GameID VARCHAR(36),
			TeamID VARCHAR(36),
			Score SMALLINT,
			
			PRIMARY KEY (GameID, TeamID)
		)`)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS TournamentGames (
			Round   SMALLINT,
			Idx     SMALLINT,
			TeamID1 VARCHAR(36),
			TeamID2 VARCHAR(36),
			Winner  VARCHAR(36),
			
			PRIMARY KEY (Round, Idx)
		)`)

	return err
}

func (s Repository) Create(ctx context.Context, teamID1, teamID2 string) (string, error) {
	id := uuid.NewString()

	_, err := s.b.InsertIntoTable("Scores").
		Fields("GameID", "TeamID", "Score").
		Values(id, teamID1, 0).
		Values(id, teamID2, 0).
		ExecContext(ctx, s.db)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s Repository) InitializeTournament(ctx context.Context, numTeams int) error {
	if (numTeams-1)&(numTeams) != 0 {
		return errors.New("number of teams must be a power of two")
	}

	b := s.b.InsertIntoTable("TournamentGames").
		Fields("Round", "Idx")

	numGamesInRound := numTeams / 2
	round := 0
	for numGamesInRound > 0 {
		for idx := range numGamesInRound {
			b = b.Values(round, idx)
		}
		numGamesInRound /= 2
		round++
	}

	_, err := b.ExecContext(ctx, s.db)
	return err
}

func (s Repository) DeleteTournament(ctx context.Context) error {
	return s.db.ExecVoid(ctx, `DELETE FROM TournamentGames`)
}

type TournamentGame struct {
	Round   int
	TeamIDs [2]string
	Winner  string
}

type Round struct {
	N     int
	Games []TournamentGame
}

type Tournament struct {
	Rounds []Round
}

func (s Repository) LoadTournament(ctx context.Context) (Tournament, error) {
	rows, err := s.db.QueryContext(
		ctx,
		// Ordering by DESC here allows us to allocate the exact size of the various arrays below.
		`SELECT Round, Idx, TeamID1, TeamID2, Winner FROM TournamentGames
		ORDER BY Round DESC, Idx DESC`,
	)
	if err != nil {
		return Tournament{}, err
	}

	var tourney Tournament
	for rows.Next() {
		var round, idx int
		var teamID1, teamID2, winner sql.Null[string]
		err := rows.Scan(&round, &idx, &teamID1, &teamID2, &winner)
		if err != nil {
			return Tournament{}, err
		}

		if len(tourney.Rounds) == 0 {
			tourney.Rounds = make([]Round, round+1)
		}
		thisRound := tourney.Rounds[round]
		if len(thisRound.Games) == 0 {
			thisRound.Games = make([]TournamentGame, idx+1)
		}
		thisGame := thisRound.Games[idx]
		thisGame.TeamIDs = [2]string{teamID1.V, teamID2.V}
		thisGame.Winner = winner.V
		thisRound.Games[idx] = thisGame
		tourney.Rounds[round] = thisRound
	}

	return tourney, nil
}

func (s Repository) PutTeam1IntoTournamentGame(ctx context.Context, round, idx int, teamID string) error {
	return s.db.ExecOne(
		ctx,
		`UPDATE TournamentGames SET TeamID1 = ? 
		WHERE Round = ? AND Idx = ? AND TeamID1 IS NULL`,
		teamID, round, idx,
	)
}

func (s Repository) PutTeam2IntoTournamentGame(ctx context.Context, round, idx int, teamID string) error {
	return s.db.ExecVoid(
		ctx,
		`UPDATE TournamentGames SET TeamID2 = ? 
		WHERE Round = ? AND Idx = ? AND TeamID2 IS NULL`,
		teamID, round, idx,
	)
}

func (s Repository) SetTournamentGameWinner(ctx context.Context, round, idx int, winner string) error {
	return s.db.ExecOne(
		ctx,
		`UPDATE TournamentGames SET Winner = ? WHERE Round = ? AND Idx = ?`,
		winner,
		round,
		idx,
	)
}

func (s Repository) UpdateScore(ctx context.Context, gameID, teamID string, score int) error {
	_, err := s.b.UpdateTable("Scores").SetFieldTo("Score", score).WhereAll(
		filter.Equals("GameID", gameID),
		filter.Equals("TeamID", teamID),
	).ExecContext(ctx, s.db)
	if err != nil {
		return err
	}

	s.scoreNotifier.Notify()
	return nil
}

func (s Repository) UpdateScores(
	ctx context.Context,
	gameID string,
	team1ID string,
	team1Score int,
	team2ID string,
	team2Score int,
) error {
	return s.db.WithTx(ctx, func(ctx context.Context) error {
		err := s.db.ExecOne(
			ctx,
			"UPDATE Scores SET Score = ? WHERE GameID = ? AND TeamID = ?",
			team1Score,
			gameID,
			team1ID,
		)
		if err != nil {
			return err
		}

		err = s.db.ExecOne(
			ctx,
			"UPDATE Scores SET Score = ? WHERE GameID = ? AND TeamID = ?",
			team2Score,
			gameID,
			team2ID,
		)

		if err != nil {
			return err
		}

		s.scoreNotifier.Notify()
		return nil
	})
}

func (s Repository) GetScore(ctx context.Context, gameID, teamID string) (int, error) {
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

func (s Repository) GetAll(ctx context.Context) ([]Score, error) {
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

func (s Repository) DeleteAll(ctx context.Context) error {
	_, err := s.b.DeleteFromTable("Scores").ExecContext(ctx, s.db)
	return err
}

func (s Repository) Get(ctx context.Context, id string) ([2]Score, error) {
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

func (s Repository) GetForTeam(ctx context.Context, teamID string) (map[string][2]Score, error) {
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

func (s Repository) GetStandings(ctx context.Context) ([]Standing, error) {
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
