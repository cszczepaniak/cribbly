package games

import (
	"fmt"
	"testing"

	"github.com/cszczepaniak/cribbly/internal/assert"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"

	"github.com/cszczepaniak/cribbly/internal/notifier"
)

func TestGames(t *testing.T) {
	n := &notifier.Notifier{}
	db := database.NewInMemory(t)
	s := NewRepository(db, n)
	assert.NoError(t, s.Init(t.Context()))

	t1 := "a"
	t2 := "b"
	t3 := "c"

	g1, err := s.Create(t.Context(), t1, t2)
	assert.NoError(t, err)

	g2, err := s.Create(t.Context(), t1, t3)
	assert.NoError(t, err)

	g3, err := s.Create(t.Context(), t2, t3)
	assert.NoError(t, err)

	gs, err := s.GetForTeam(t.Context(), t1)
	assert.NoError(t, err)

	assert.MapLen(t, gs, 2)
	assert.MapHasKey(t, gs, g1)
	assert.Equal(t, [2]Score{
		{GameID: g1, TeamID: "a", Score: 0},
		{GameID: g1, TeamID: "b", Score: 0},
	}, gs[g1])
	assert.MapHasKey(t, gs, g2)
	assert.Equal(t, [2]Score{
		{GameID: g2, TeamID: "a", Score: 0},
		{GameID: g2, TeamID: "c", Score: 0},
	}, gs[g2])

	gs, err = s.GetForTeam(t.Context(), t2)
	assert.NoError(t, err)

	assert.MapLen(t, gs, 2)
	assert.MapHasKey(t, gs, g1)
	assert.Equal(t, [2]Score{
		{GameID: g1, TeamID: "a", Score: 0},
		{GameID: g1, TeamID: "b", Score: 0},
	}, gs[g1])
	assert.MapHasKey(t, gs, g3)
	assert.Equal(t, [2]Score{
		{GameID: g3, TeamID: "b", Score: 0},
		{GameID: g3, TeamID: "c", Score: 0},
	}, gs[g3])

	gs, err = s.GetForTeam(t.Context(), t3)
	assert.NoError(t, err)

	assert.MapLen(t, gs, 2)
	assert.MapHasKey(t, gs, g2)
	assert.Equal(t, [2]Score{
		{GameID: g2, TeamID: "a", Score: 0},
		{GameID: g2, TeamID: "c", Score: 0},
	}, gs[g2])
	assert.MapHasKey(t, gs, g3)
	assert.Equal(t, [2]Score{
		{GameID: g3, TeamID: "b", Score: 0},
		{GameID: g3, TeamID: "c", Score: 0},
	}, gs[g3])
}

func TestGames_Notifications(t *testing.T) {
	n := &notifier.Notifier{}
	db := database.NewInMemory(t)
	s := NewRepository(db, n)
	assert.NoError(t, s.Init(t.Context()))

	t1 := "a"
	t2 := "b"
	t3 := "c"

	g1, err := s.Create(t.Context(), t1, t2)
	assert.NoError(t, err)

	g2, err := s.Create(t.Context(), t1, t3)
	assert.NoError(t, err)

	sub, cancel := n.Subscribe()
	t.Cleanup(cancel)

	assert.NoError(t, s.UpdateScore(t.Context(), g1, t1, 100))
	// should get notified
	<-sub

	assert.NoError(t, s.UpdateScores(t.Context(), g2, t1, 121, t3, 99))
	// should get notified
	<-sub
}

func TestGames_UpdateScores_TeamsMustExistForGame(t *testing.T) {
	n := &notifier.Notifier{}
	db := database.NewInMemory(t)
	s := NewRepository(db, n)
	assert.NoError(t, s.Init(t.Context()))

	t1 := "a"
	t2 := "b"

	g, err := s.Create(t.Context(), t1, t2)
	assert.NoError(t, err)

	err = s.UpdateScores(t.Context(), g, t1, 121, "not a team", 100)
	assert.Error(t, err)

	var exists bool
	err = db.QueryRowContext(
		t.Context(),
		"SELECT EXISTS(SELECT 1 FROM Scores WHERE TeamID = ?)",
		"not a team",
	).Scan(&exists)
	assert.False(t, exists)
}

func TestTournamentGames(t *testing.T) {
	n := &notifier.Notifier{}
	db := database.NewInMemory(t)
	s := NewRepository(db, n)
	assert.NoError(t, s.Init(t.Context()))

	err := s.InitializeTournament(t.Context(), 17)
	assert.Error(t, err)

	assert.NoError(t, s.InitializeTournament(t.Context(), 32))

	tourney, err := s.LoadTournament(t.Context())
	assert.NoError(t, err)

	assert.SliceLen(t, tourney.Rounds, 5)
	assert.SliceLen(t, tourney.Rounds[0].Games, 16)
	assert.SliceLen(t, tourney.Rounds[1].Games, 8)
	assert.SliceLen(t, tourney.Rounds[2].Games, 4)
	assert.SliceLen(t, tourney.Rounds[3].Games, 2)
	assert.SliceLen(t, tourney.Rounds[4].Games, 1)

	for i := range 32 {
		if i%2 == 0 {
			assert.NoError(t, s.PutTeam1IntoTournamentGame(t.Context(), 0, i/2, fmt.Sprintf("team%d", i)))
		} else {
			assert.NoError(t, s.PutTeam2IntoTournamentGame(t.Context(), 0, i/2, fmt.Sprintf("team%d", i)))
		}
	}

	tourney, err = s.LoadTournament(t.Context())
	assert.NoError(t, err)
	assert.SliceLen(t, tourney.Rounds, 5)
	for i, g := range tourney.Rounds[0].Games {
		t1 := fmt.Sprintf("team%d", 2*i)
		t2 := fmt.Sprintf("team%d", 2*i+1)
		assert.Equal(t, t1, g.TeamIDs[0])
		assert.Equal(t, t2, g.TeamIDs[1])

		// Advance team1
		assert.NoError(t, s.SetTournamentGameWinner(t.Context(), g.Round, i, t1))

		if i%2 == 0 {
			assert.NoError(t, s.PutTeam1IntoTournamentGame(t.Context(), g.Round+1, i/2, t1))
		} else {
			assert.NoError(t, s.PutTeam2IntoTournamentGame(t.Context(), g.Round+1, i/2, t1))
		}
	}

	tourney, err = s.LoadTournament(t.Context())
	assert.NoError(t, err)
	assert.SliceLen(t, tourney.Rounds, 5)
	for i, g := range tourney.Rounds[0].Games {
		t1 := fmt.Sprintf("team%d", 2*i)
		assert.Equal(t, t1, g.Winner)
	}
	for i, g := range tourney.Rounds[1].Games {
		t1 := fmt.Sprintf("team%d", 4*i)
		t2 := fmt.Sprintf("team%d", 4*i+2)
		assert.Equal(t, t1, g.TeamIDs[0])
		assert.Equal(t, t2, g.TeamIDs[1])
	}
}
