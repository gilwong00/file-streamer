package transferservice

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	transferv1 "github.com/gilwong00/file-streamer/internal/gen/proto/v1"
)

func (s *transferService) StreamFile(
	ctx context.Context,
	req *connect.Request[transferv1.StreamFileRequest],
	stream *connect.ServerStream[transferv1.StreamFileResponse],
) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("this method is not implemented"))
}
