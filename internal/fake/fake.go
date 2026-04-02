package fake

import (
	_ "embed"
	"math/rand/v2"
	"strings"
)

var (
	//go:embed names.txt
	namesEmbed string

	firstNames []string
	lastNames  []string

	//go:embed states.txt
	statesEmbed string
	states      []string
)

func init() {
	for l := range strings.SplitSeq(namesEmbed, "\n") {
		first, last, _ := strings.Cut(l, " ")
		firstNames = append(firstNames, strings.TrimSpace(first))
		lastNames = append(firstNames, strings.TrimSpace(last))
	}

	for l := range strings.SplitSeq(statesEmbed, "\n") {
		states = append(states, strings.TrimSpace(l))
	}
}

func FirstName() string {
	return firstNames[rand.IntN(len(firstNames))]
}

func LastName() string {
	return lastNames[rand.IntN(len(lastNames))]
}

func USState() string {
	return states[rand.IntN(len(states))]
}
