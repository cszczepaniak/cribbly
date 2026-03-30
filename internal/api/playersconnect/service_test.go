package playersconnect

import (
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/cszczepaniak/gotest/assert"
	"github.com/google/uuid"

	cribblyv1 "github.com/cszczepaniak/cribbly/internal/gen/cribbly/v1"
	"github.com/cszczepaniak/cribbly/internal/persistence/database"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/server/middleware"
)

func newTestServer(t *testing.T) (*Server, players.Repository) {
	t.Helper()
	db := database.NewInMemory(t)
	repo := players.NewRepository(db)
	assert.NoError(t, repo.Init(t.Context()))
	return &Server{PlayerRepo: repo}, repo
}

func assertConnectCode(t *testing.T, err error, want connect.Code) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error")
	}
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		t.Fatalf("expected *connect.Error, got %T: %v", err, err)
	}
	assert.Equal(t, want, connectErr.Code())
}

func TestListPlayers_NotAdmin(t *testing.T) {
	svc, _ := newTestServer(t)
	_, err := svc.ListPlayers(
		t.Context(),
		connect.NewRequest(&cribblyv1.ListPlayersRequest{}),
	)
	assertConnectCode(t, err, connect.CodePermissionDenied)
}

func TestListPlayers_WithDevAdmin(t *testing.T) {
	svc, _ := newTestServer(t)
	ctx := middleware.WithDevAdminContext(t.Context())
	resp, err := svc.ListPlayers(
		ctx,
		connect.NewRequest(&cribblyv1.ListPlayersRequest{}),
	)
	assert.NoError(t, err)
	assert.SliceLen(t, resp.Msg.GetPlayers(), 0)
}

func TestCreatePlayer_NotAdmin(t *testing.T) {
	svc, _ := newTestServer(t)
	_, err := svc.CreatePlayer(
		t.Context(),
		connect.NewRequest(&cribblyv1.CreatePlayerRequest{
			FirstName: "A",
			LastName:  "B",
		}),
	)
	assertConnectCode(t, err, connect.CodePermissionDenied)
}

func TestCreatePlayer_InvalidArgument_EmptyName(t *testing.T) {
	svc, _ := newTestServer(t)
	ctx := middleware.WithDevAdminContext(t.Context())
	_, err := svc.CreatePlayer(
		ctx,
		connect.NewRequest(&cribblyv1.CreatePlayerRequest{
			FirstName: "",
			LastName:  "B",
		}),
	)
	assertConnectCode(t, err, connect.CodeInvalidArgument)
}

func TestCreatePlayer_Success_ReturnsList(t *testing.T) {
	svc, _ := newTestServer(t)
	ctx := middleware.WithDevAdminContext(t.Context())
	resp, err := svc.CreatePlayer(
		ctx,
		connect.NewRequest(&cribblyv1.CreatePlayerRequest{
			FirstName: "Ada",
			LastName:  "Lovelace",
		}),
	)
	assert.NoError(t, err)
	ps := resp.Msg.GetPlayers()
	assert.SliceLen(t, ps, 1)
	assert.Equal(t, "Ada", ps[0].GetFirstName())
	assert.Equal(t, "Lovelace", ps[0].GetLastName())
	assert.Equal(t, "", ps[0].GetTeamId())
}

func TestUpdatePlayer_NotFound(t *testing.T) {
	svc, _ := newTestServer(t)
	ctx := middleware.WithDevAdminContext(t.Context())
	_, err := svc.UpdatePlayer(
		ctx,
		connect.NewRequest(&cribblyv1.UpdatePlayerRequest{
			Id:        uuid.NewString(),
			FirstName: "A",
			LastName:  "B",
		}),
	)
	assertConnectCode(t, err, connect.CodeNotFound)
}

func TestUpdatePlayer_InvalidArgument_EmptyName(t *testing.T) {
	svc, repo := newTestServer(t)
	ctx := middleware.WithDevAdminContext(t.Context())
	id, err := repo.Create(t.Context(), "Old", "Name")
	assert.NoError(t, err)

	_, err = svc.UpdatePlayer(
		ctx,
		connect.NewRequest(&cribblyv1.UpdatePlayerRequest{
			Id:        id,
			FirstName: "",
			LastName:  "B",
		}),
	)
	assertConnectCode(t, err, connect.CodeInvalidArgument)
}

func TestUpdatePlayer_Success(t *testing.T) {
	svc, repo := newTestServer(t)
	ctx := middleware.WithDevAdminContext(t.Context())
	id, err := repo.Create(t.Context(), "Old", "Name")
	assert.NoError(t, err)

	resp, err := svc.UpdatePlayer(
		ctx,
		connect.NewRequest(&cribblyv1.UpdatePlayerRequest{
			Id:        id,
			FirstName: "New",
			LastName:  "Name",
		}),
	)
	assert.NoError(t, err)
	ps := resp.Msg.GetPlayers()
	assert.SliceLen(t, ps, 1)
	assert.Equal(t, "New", ps[0].GetFirstName())
	assert.Equal(t, "Name", ps[0].GetLastName())
}

func TestDeletePlayer_NotFound(t *testing.T) {
	svc, _ := newTestServer(t)
	ctx := middleware.WithDevAdminContext(t.Context())
	_, err := svc.DeletePlayer(
		ctx,
		connect.NewRequest(&cribblyv1.DeletePlayerRequest{
			Id: uuid.NewString(),
		}),
	)
	assertConnectCode(t, err, connect.CodeNotFound)
}

func TestDeletePlayer_EmptyId(t *testing.T) {
	svc, _ := newTestServer(t)
	ctx := middleware.WithDevAdminContext(t.Context())
	_, err := svc.DeletePlayer(
		ctx,
		connect.NewRequest(&cribblyv1.DeletePlayerRequest{Id: "  "}),
	)
	assertConnectCode(t, err, connect.CodeInvalidArgument)
}

func TestDeletePlayer_OnTeam_FailedPrecondition(t *testing.T) {
	svc, repo := newTestServer(t)
	ctx := middleware.WithDevAdminContext(t.Context())

	id, err := repo.Create(t.Context(), "On", "Team")
	assert.NoError(t, err)
	teamID := uuid.NewString()
	assert.NoError(t, repo.AssignToTeam(t.Context(), id, teamID))

	_, err = svc.DeletePlayer(
		ctx,
		connect.NewRequest(&cribblyv1.DeletePlayerRequest{Id: id}),
	)
	assertConnectCode(t, err, connect.CodeFailedPrecondition)
}

func TestDeletePlayer_Success(t *testing.T) {
	svc, repo := newTestServer(t)
	ctx := middleware.WithDevAdminContext(t.Context())

	id, err := repo.Create(t.Context(), "Free", "Agent")
	assert.NoError(t, err)

	resp, err := svc.DeletePlayer(
		ctx,
		connect.NewRequest(&cribblyv1.DeletePlayerRequest{Id: id}),
	)
	assert.NoError(t, err)
	assert.SliceLen(t, resp.Msg.GetPlayers(), 0)
}

func TestDeleteAllPlayers_UnassignsAndDeletes(t *testing.T) {
	svc, repo := newTestServer(t)
	ctx := middleware.WithDevAdminContext(t.Context())

	id1, err := repo.Create(t.Context(), "A", "One")
	assert.NoError(t, err)
	_, err = repo.Create(t.Context(), "B", "Two")
	assert.NoError(t, err)
	assert.NoError(t, repo.AssignToTeam(t.Context(), id1, uuid.NewString()))

	_, err = svc.DeleteAllPlayers(
		ctx,
		connect.NewRequest(&cribblyv1.DeleteAllPlayersRequest{}),
	)
	assert.NoError(t, err)

	all, err := repo.GetAll(t.Context())
	assert.NoError(t, err)
	assert.SliceLen(t, all, 0)
}

func TestGenerateRandomPlayers_NotAdmin(t *testing.T) {
	svc, _ := newTestServer(t)
	_, err := svc.GenerateRandomPlayers(
		t.Context(),
		connect.NewRequest(&cribblyv1.GenerateRandomPlayersRequest{Count: 3}),
	)
	assertConnectCode(t, err, connect.CodePermissionDenied)
}

func TestGenerateRandomPlayers_InvalidCount(t *testing.T) {
	svc, _ := newTestServer(t)
	ctx := middleware.WithDevAdminContext(t.Context())

	_, err := svc.GenerateRandomPlayers(
		ctx,
		connect.NewRequest(&cribblyv1.GenerateRandomPlayersRequest{Count: 0}),
	)
	assertConnectCode(t, err, connect.CodeInvalidArgument)

	_, err = svc.GenerateRandomPlayers(
		ctx,
		connect.NewRequest(&cribblyv1.GenerateRandomPlayersRequest{Count: 501}),
	)
	assertConnectCode(t, err, connect.CodeInvalidArgument)
}

func TestGenerateRandomPlayers_WithDevAdmin(t *testing.T) {
	svc, _ := newTestServer(t)
	ctx := middleware.WithDevAdminContext(t.Context())

	resp, err := svc.GenerateRandomPlayers(
		ctx,
		connect.NewRequest(&cribblyv1.GenerateRandomPlayersRequest{Count: 4}),
	)
	assert.NoError(t, err)
	assert.SliceLen(t, resp.Msg.GetPlayers(), 4)
}
