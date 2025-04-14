package httptestutils

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/a-h/htmlformat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Response struct {
	bytes []byte
}

// String returns the contents of the response's body as a string.
func (r Response) String() string {
	return string(r.bytes)
}

// Reader returns an [io.Reader] with the contents of the response.
func (r Response) Reader() io.Reader {
	return bytes.NewReader(r.bytes)
}

// CompareHTMLDocumentToSnapshot compares the respone contents to the given snapshot file. If the snapshot
// file doesn't exist, it will be created with the contents of the request and the test will fail.
// If the contents of the file don't match the contents of the request, the test will fail and the
// contents of the request will replace the contents of the file. If this happens and the changes
// are expected, the new snapshot file should be committed.
//
// This function differs from [CompareHTMLFragmentToSnapshot] in that it should be used when a full
// HTML document is expected. That is, if your HTML has an <html> tag, this is the more appropriate
// function to use.
func (r Response) CompareHTMLDocumentToSnapshot(t *testing.T, snapshotFile string) {
	t.Helper()

	r.compareHTMLToSnapshot(t, snapshotFile, htmlformat.Document)
}

// CompareHTMLFragmentToSnapshot compares the response contents to the given snapshot file. If the
// snapshot file doesn't exist, it will be created with the contents of the request and the test
// will fail. If the contents of the file don't match the contents of the request, the test will
// fail and the contents of the request will replace the contents of the file. If this happens and
// the changes are expected, the new snapshot file should be committed.
//
// This function differs from [CompareHTMLDocumentToSnapshot] in that it should be used when only a
// fragment of HTML is expected, not a full document. That is, if your HTML doesn't have an <html>
// tag, this is the more appropriate function to use.
func (r Response) CompareHTMLFragmentToSnapshot(t *testing.T, snapshotFile string) {
	t.Helper()

	r.compareHTMLToSnapshot(t, snapshotFile, htmlformat.Fragment)
}

func (r Response) compareHTMLToSnapshot(
	t *testing.T,
	snapshotFile string,
	formatFn func(w io.Writer, r io.Reader) error,
) {
	t.Helper()

	formatHTMLIntoFile := func(r io.Reader, flag int) {
		t.Helper()

		f, err := os.OpenFile(snapshotFile, flag, os.ModePerm)
		require.NoError(t, err)
		defer f.Close()

		require.NoError(t, formatFn(f, r))
		require.NoError(t, f.Close())
	}

	snapshot, err := os.ReadFile(snapshotFile)
	if os.IsNotExist(err) {
		require.NoError(t, os.MkdirAll(filepath.Dir(snapshotFile), os.ModePerm))

		// The snapshot file didn't exist; we're going to create it and then fail the test.
		formatHTMLIntoFile(r.Reader(), os.O_WRONLY|os.O_CREATE)
		require.FailNow(t, "no shapshot existed; we created one for you. rerun the test!")
	}

	require.NoError(t, err)

	contents := &strings.Builder{}
	formatFn(contents, r.Reader())

	if !assert.Equal(t, string(snapshot), contents.String()) {
		// Format the new HTML response into the file.
		formatHTMLIntoFile(
			// We already read the whole body, so we have to use the string builder as our reader
			// since it now hold those contents.
			strings.NewReader(contents.String()),
			// Truncate the existing file.
			os.O_WRONLY|os.O_TRUNC,
		)
		require.FailNow(t, "snapshot didn't match! if the changes are expected, commit the new snapshot.")
	}
}
