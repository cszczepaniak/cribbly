package roomcodeconnect

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	cribblyv1 "github.com/cszczepaniak/cribbly/internal/gen/cribbly/v1"
	cribblyv1connect "github.com/cszczepaniak/cribbly/internal/gen/cribbly/v1/cribblyv1connect"
	"github.com/cszczepaniak/gotest/assert"

	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/roomcodes"
)

func TestSetRoomCode_SetsCookieHeader(t *testing.T) {
	db := database.NewInMemory(t)
	repo := roomcodes.NewRepository(db)
	assert.NoError(t, repo.Init(t.Context()))
	assert.NoError(t, repo.Create(t.Context(), "GOODCODE", time.Now().Add(time.Hour)))

	svc := &Server{Repo: repo}
	_, h := cribblyv1connect.NewRoomCodeServiceHandler(svc)
	ts := httptest.NewServer(http.StripPrefix("/api", h))
	defer ts.Close()

	client := cribblyv1connect.NewRoomCodeServiceClient(http.DefaultClient, ts.URL+"/api")
	resp, err := client.SetRoomCode(
		context.Background(),
		connect.NewRequest(&cribblyv1.SetRoomCodeRequest{Code: "GOODCODE"}),
	)
	assert.NoError(t, err)
	if resp == nil {
		t.Fatal("expected response")
	}
	if len(resp.Header().Values("Set-Cookie")) == 0 {
		t.Fatal("expected Set-Cookie header")
	}
}

func TestSetRoomCode_InvalidCode(t *testing.T) {
	db := database.NewInMemory(t)
	repo := roomcodes.NewRepository(db)
	assert.NoError(t, repo.Init(t.Context()))

	svc := &Server{Repo: repo}
	_, h := cribblyv1connect.NewRoomCodeServiceHandler(svc)
	ts := httptest.NewServer(http.StripPrefix("/api", h))
	defer ts.Close()

	client := cribblyv1connect.NewRoomCodeServiceClient(http.DefaultClient, ts.URL+"/api")
	_, err := client.SetRoomCode(
		context.Background(),
		connect.NewRequest(&cribblyv1.SetRoomCodeRequest{Code: "BADCODE"}),
	)
	assert.Error(t, err)
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		t.Fatalf("expected *connect.Error, got %T: %v", err, err)
	}
	assert.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
}
