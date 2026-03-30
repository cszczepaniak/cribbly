package playersconnect

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"github.com/jaswdr/faker/v2"

	cribblyv1 "github.com/cszczepaniak/cribbly/internal/gen/cribbly/v1"
	"github.com/cszczepaniak/cribbly/internal/moreiter"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
	"github.com/cszczepaniak/cribbly/internal/server/middleware"
)

type Server struct {
	PlayerRepo players.Repository
}

func requireAdmin(ctx context.Context) error {
	if !middleware.IsAdmin(ctx) {
		return connect.NewError(connect.CodePermissionDenied, errors.New("must be an admin"))
	}
	return nil
}

func toProto(ps []players.Player) []*cribblyv1.Player {
	out := make([]*cribblyv1.Player, 0, len(ps))
	for _, p := range ps {
		out = append(out, &cribblyv1.Player{
			Id:        p.ID,
			FirstName: p.FirstName,
			LastName:  p.LastName,
			TeamId:    p.TeamID,
		})
	}
	return out
}

func (s *Server) ListPlayers(
	ctx context.Context,
	_ *connect.Request[cribblyv1.ListPlayersRequest],
) (*connect.Response[cribblyv1.ListPlayersResponse], error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	ps, err := s.PlayerRepo.GetAll(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&cribblyv1.ListPlayersResponse{Players: toProto(ps)}), nil
}

func (s *Server) CreatePlayer(
	ctx context.Context,
	req *connect.Request[cribblyv1.CreatePlayerRequest],
) (*connect.Response[cribblyv1.CreatePlayerResponse], error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	first := strings.TrimSpace(req.Msg.GetFirstName())
	last := strings.TrimSpace(req.Msg.GetLastName())
	if first == "" || last == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("first and last name are required"))
	}

	if _, err := s.PlayerRepo.Create(ctx, first, last); err != nil {
		if strings.Contains(err.Error(), "must have a first and last name") {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	ps, err := s.PlayerRepo.GetAll(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&cribblyv1.CreatePlayerResponse{Players: toProto(ps)}), nil
}

func (s *Server) DeletePlayer(
	ctx context.Context,
	req *connect.Request[cribblyv1.DeletePlayerRequest],
) (*connect.Response[cribblyv1.DeletePlayerResponse], error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	id := strings.TrimSpace(req.Msg.GetId())
	if id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("id is required"))
	}

	p, err := s.PlayerRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("player not found"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if p.TeamID != "" {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("remove the player from their team before deleting"))
	}

	if err := s.PlayerRepo.Delete(ctx, id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	ps, err := s.PlayerRepo.GetAll(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&cribblyv1.DeletePlayerResponse{Players: toProto(ps)}), nil
}

func (s *Server) DeleteAllPlayers(
	ctx context.Context,
	_ *connect.Request[cribblyv1.DeleteAllPlayersRequest],
) (*connect.Response[cribblyv1.DeleteAllPlayersResponse], error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	ps, err := s.PlayerRepo.GetAll(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	for _, p := range ps {
		if p.TeamID != "" {
			err := s.PlayerRepo.UnassignFromTeam(ctx, p.TeamID, moreiter.Of(p.ID))
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
		}

		if err := s.PlayerRepo.Delete(ctx, p.ID); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return connect.NewResponse(&cribblyv1.DeleteAllPlayersResponse{}), nil
}

func (s *Server) GenerateRandomPlayers(
	ctx context.Context,
	req *connect.Request[cribblyv1.GenerateRandomPlayersRequest],
) (*connect.Response[cribblyv1.GenerateRandomPlayersResponse], error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	n := int(req.Msg.GetCount())
	if n < 1 || n > 500 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("count must be between 1 and 500"))
	}

	fake := faker.New()
	for range n {
		firstName := fake.Person().FirstName()
		lastName := fake.Person().LastName()
		if _, err := s.PlayerRepo.Create(ctx, firstName, lastName); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	ps, err := s.PlayerRepo.GetAll(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&cribblyv1.GenerateRandomPlayersResponse{Players: toProto(ps)}), nil
}
