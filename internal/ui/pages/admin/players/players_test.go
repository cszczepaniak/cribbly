package players

import (
	"testing"

	"github.com/cszczepaniak/cribbly/internal/httptestutils"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/persistence/sqlite"
	"github.com/stretchr/testify/require"
)

func TestPlayerRegistrationForm(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	ps := players.NewService(db)
	require.NoError(t, ps.Init(t.Context()))

	// Seed some players into the DB to test that we render them properly in the list.
	_, err := ps.Create(t.Context(), "Mario")
	require.NoError(t, err)

	_, err = ps.Create(t.Context(), "Luigi")
	require.NoError(t, err)

	h := PlayersHandler{
		PlayerService: ps,
	}
	s := httptestutils.NewServerForComponent(t, h.RegistrationPage)
	resp := s.Get(t, "/")
	resp.CompareHTMLDocumentToSnapshot(t, "testdata/player_registration_page_snapshot.html")
}

func TestPostPlayer(t *testing.T) {
	db := sqlite.NewInMemoryForTest(t)
	ps := players.NewService(db)
	require.NoError(t, ps.Init(t.Context()))

	// Seed some players into the DB to test that we render them properly in the list.
	_, err := ps.Create(t.Context(), "Mario")
	require.NoError(t, err)

	_, err = ps.Create(t.Context(), "Luigi")
	require.NoError(t, err)

	h := PlayersHandler{
		PlayerService: ps,
	}
	s := httptestutils.NewServerForComponent(t, h.PostPlayer)

	// Add Waluigi to the database
	resp := s.Post(t, "/?name=Waluigi", nil)
	resp.CompareHTMLFragmentToSnapshot(t, "testdata/post_player_snapshot.html")
}
