package transferservice

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	transferv1 "github.com/gilwong00/file-streamer/internal/gen/proto/v1"
)

func (s *transferService) UploadFile(
	ctx context.Context,
	req *connect.BidiStream[transferv1.UploadFileRequest, transferv1.UploadFileResponse],
) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("this method is not implemented"))
}
