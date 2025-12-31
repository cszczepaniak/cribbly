package games

import (
	"cmp"
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
)

func TestGeneratePairs_4(t *testing.T) {
	allTeams := []teams.Team{{
		ID: "1",
	}, {
		ID: "2",
	}, {
		ID: "3",
	}, {
		ID: "4",
	}}
	pairs, err := generateMatchups(allTeams)
	require.NoError(t, err)

	assert.Equal(t, []string{
		"1,2",
		"1,3",
		"1,4",
		"2,3",
		"2,4",
		"3,4",
	}, pairsToStrings(pairs))
}

func TestGeneratePairs_5(t *testing.T) {
	allTeams := []teams.Team{{
		ID: "1",
	}, {
		ID: "2",
	}, {
		ID: "3",
	}, {
		ID: "4",
	}, {
		ID: "5",
	}}
	_, err := generateMatchups(allTeams)
	require.Error(t, err)
}

func TestGeneratePairs_6(t *testing.T) {
	allTeams := []teams.Team{{
		ID: "1",
	}, {
		ID: "2",
	}, {
		ID: "3",
	}, {
		ID: "4",
	}, {
		ID: "5",
	}, {
		ID: "6",
	}}
	pairs, err := generateMatchups(allTeams)
	require.NoError(t, err)

	assert.Equal(t, []string{
		"1,2",
		"1,3",
		"1,6",
		"2,3",
		"2,4",
		"3,4",
		"4,5",
		// Note: with 6 teams, these two teams must either play a duplicate game, or one of the
		// other teams in the division would need to play 4 games.
		"5,6",
		"5,6",
	}, pairsToStrings(pairs))
}

func pairsToStrings(pairs [][2]teams.Team) []string {
	var s []string
	for _, p := range pairs {
		slices.SortFunc(p[:], func(a, b teams.Team) int {
			return cmp.Compare(a.ID, b.ID)
		})
		s = append(s, fmt.Sprintf("%s,%s", p[0].ID, p[1].ID))
	}
	slices.Sort(s)
	return s
}
