package roomcodeconnect

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"

	cribblyv1 "github.com/cszczepaniak/cribbly/internal/gen/cribbly/v1"
	"github.com/cszczepaniak/cribbly/internal/persistence/roomcodes"
	"github.com/cszczepaniak/cribbly/internal/server/middleware"
)

// Server implements cribbly.v1.RoomCodeService (Connect).
type Server struct {
	Repo roomcodes.Repository
}

func (s *Server) SetRoomCode(
	ctx context.Context,
	req *connect.Request[cribblyv1.SetRoomCodeRequest],
) (*connect.Response[cribblyv1.SetRoomCodeResponse], error) {
	code := strings.TrimSpace(req.Msg.GetCode())
	if code == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("code is required"))
	}

	ok, err := s.Repo.Validate(ctx, code)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if !ok {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid or expired room code"))
	}

	resp := connect.NewResponse(&cribblyv1.SetRoomCodeResponse{})
	cookie := &http.Cookie{
		Name:     "room_code",
		Value:    code,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   middleware.IsProd(ctx),
		SameSite: http.SameSiteLaxMode,
	}
	resp.Header().Add("Set-Cookie", cookie.String())

	return resp, nil
}

func (s *Server) DoSomething(
	ctx context.Context,
	req *connect.Request[cribblyv1.SomethingRequest],
	stream *connect.ServerStream[cribblyv1.SomethingResponse],
) error {
	n := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}

		n++
		slog.Info("sending stream event", "val", n)

		err := stream.Send(&cribblyv1.SomethingResponse{
			Data: fmt.Sprint(n),
		})
		if err != nil {
			return err
		}
	}
}
