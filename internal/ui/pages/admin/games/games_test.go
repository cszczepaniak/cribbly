package games

import (
	"cmp"
	"fmt"
	"slices"
	"testing"

	"github.com/cszczepaniak/gotest/assert"

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
	assert.NoError(t, err)

	assert.Equal(t, []string{
		"1,2",
		"1,3",
		"1,4",
		"2,3",
		"2,4",
		"3,4",
	}, pairsToStrings(pairs))
}

func TestGeneratePairs_3(t *testing.T) {
	allTeams := []teams.Team{{ID: "1"}, {ID: "2"}, {ID: "3"}}
	pairs, err := generateMatchups(allTeams)
	assert.NoError(t, err)
	// Full round robin: 3 games total, 2 per team.
	assert.Equal(t, []string{"1,2", "1,3", "2,3"}, pairsToStrings(pairs))
}

func TestGeneratePairs_5(t *testing.T) {
	allTeams := []teams.Team{{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}, {ID: "5"}}
	pairs, err := generateMatchups(allTeams)
	assert.NoError(t, err)
	// Full round robin: 10 games, 4 per team.
	assert.Equal(t, []string{
		"1,2", "1,3", "1,4", "1,5",
		"2,3", "2,4", "2,5",
		"3,4", "3,5",
		"4,5",
	}, pairsToStrings(pairs))
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
	assert.NoError(t, err)

	// 3 games per team via 3 rounds of circle method (no duplicate pairings).
	assert.Equal(t, []string{
		"1,4", "1,5", "1,6",
		"2,3", "2,5", "2,6",
		"3,4", "3,5",
		"4,6",
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
