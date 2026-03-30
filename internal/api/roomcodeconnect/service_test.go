package roomcodeconnect

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/cszczepaniak/gotest/assert"

	cribblyv1 "github.com/cszczepaniak/cribbly/internal/gen/cribbly/v1"
	cribblyv1connect "github.com/cszczepaniak/cribbly/internal/gen/cribbly/v1/cribblyv1connect"

	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/roomcodes"
	"github.com/cszczepaniak/cribbly/internal/server/middleware"
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
		t.Context(),
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
		t.Context(),
		connect.NewRequest(&cribblyv1.SetRoomCodeRequest{Code: "BADCODE"}),
	)
	assert.Error(t, err)
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		t.Fatalf("expected *connect.Error, got %T: %v", err, err)
	}
	assert.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
}

func TestCheckRoomAccess_NoCookie(t *testing.T) {
	db := database.NewInMemory(t)
	repo := roomcodes.NewRepository(db)
	assert.NoError(t, repo.Init(t.Context()))

	svc := &Server{Repo: repo}
	_, h := cribblyv1connect.NewRoomCodeServiceHandler(svc)
	ts := httptest.NewServer(http.StripPrefix("/api", h))
	defer ts.Close()

	client := cribblyv1connect.NewRoomCodeServiceClient(http.DefaultClient, ts.URL+"/api")
	resp, err := client.CheckRoomAccess(
		t.Context(),
		connect.NewRequest(&cribblyv1.CheckRoomAccessRequest{}),
	)
	assert.NoError(t, err)
	assert.Equal(t, false, resp.Msg.GetHasAccess())
}

func TestCheckRoomAccess_ValidRoomCookie(t *testing.T) {
	db := database.NewInMemory(t)
	repo := roomcodes.NewRepository(db)
	assert.NoError(t, repo.Init(t.Context()))
	assert.NoError(t, repo.Create(t.Context(), "GOODCODE", time.Now().Add(time.Hour)))

	svc := &Server{Repo: repo}
	_, h := cribblyv1connect.NewRoomCodeServiceHandler(svc)
	ts := httptest.NewServer(http.StripPrefix("/api", h))
	defer ts.Close()

	client := cribblyv1connect.NewRoomCodeServiceClient(http.DefaultClient, ts.URL+"/api")
	req := connect.NewRequest(&cribblyv1.CheckRoomAccessRequest{})
	req.Header().Set("Cookie", "room_code=GOODCODE")

	resp, err := client.CheckRoomAccess(t.Context(), req)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Msg.GetHasAccess())
}

func TestGenerateRoomCode_NotAdmin(t *testing.T) {
	db := database.NewInMemory(t)
	repo := roomcodes.NewRepository(db)
	assert.NoError(t, repo.Init(t.Context()))

	svc := &Server{Repo: repo}
	_, err := svc.GenerateRoomCode(
		t.Context(),
		connect.NewRequest(&cribblyv1.GenerateRoomCodeRequest{}),
	)
	assert.Error(t, err)
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		t.Fatalf("expected *connect.Error, got %T: %v", err, err)
	}
	assert.Equal(t, connect.CodePermissionDenied, connectErr.Code())
}

func TestGenerateRoomCode_WithDevAdminContext(t *testing.T) {
	db := database.NewInMemory(t)
	repo := roomcodes.NewRepository(db)
	assert.NoError(t, repo.Init(t.Context()))

	svc := &Server{Repo: repo}
	ctx := middleware.WithDevAdminContext(t.Context())
	resp, err := svc.GenerateRoomCode(
		ctx,
		connect.NewRequest(&cribblyv1.GenerateRoomCodeRequest{}),
	)
	assert.NoError(t, err)
	if len(resp.Msg.GetCode()) != 6 {
		t.Fatalf("expected 6-char code, got %q", resp.Msg.GetCode())
	}
	if resp.Msg.GetExpiresAt() == "" {
		t.Fatal("expected expires_at")
	}
}
