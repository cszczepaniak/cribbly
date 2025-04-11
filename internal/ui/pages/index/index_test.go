package index

import (
	"testing"

	"github.com/cszczepaniak/cribbly/internal/httptestutils"
)

func TestIndex(t *testing.T) {
	s := httptestutils.NewServerForComponent(t, Index)
	resp := s.Get(t, "/")
	resp.CompareHTMLDocumentToSnapshot(t, "testdata/index_snapshot.html")
}
