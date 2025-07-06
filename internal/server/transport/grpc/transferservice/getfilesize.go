package transferservice

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	transferv1 "github.com/gilwong00/file-streamer/internal/gen/proto/v1"
)

func (s *transferService) GetFileSize(
	ctx context.Context,
	req *connect.Request[transferv1.GetFileSizeRequest],
) (*connect.Response[transferv1.GetFileSizeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("this method is not implemented"))
}
