package roomcodeconnect

import (
	"context"
	"errors"
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
