package transferservice

import (
	"context"

	"connectrpc.com/connect"
	transferv1 "github.com/gilwong00/file-streamer/internal/gen/proto/v1"
)

func (s *transferService) GetFileSize(
	ctx context.Context,
	req *connect.Request[transferv1.GetFileSizeRequest],
) (*connect.Response[transferv1.GetFileSizeResponse], error) {
	info, err := s.storageClient.GetObjectInfo(ctx, s.bucketName, req.Msg.FileName)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	return connect.NewResponse(&transferv1.GetFileSizeResponse{
		Size: info.Size,
	}), nil
}
