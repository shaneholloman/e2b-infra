package filesystem

import (
	"context"
	"fmt"
	"os"

	"connectrpc.com/connect"

	"github.com/e2b-dev/infra/packages/envd/internal/permissions"
	rpc "github.com/e2b-dev/infra/packages/envd/internal/services/spec/filesystem"
)

func (Service) Stat(ctx context.Context, req *connect.Request[rpc.StatRequest]) (*connect.Response[rpc.StatResponse], error) {
	u, err := permissions.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}

	path, err := permissions.ExpandAndResolve(req.Msg.GetPath(), u)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("file not found: %w", err))
		}

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("error statting file: %w", err))
	}

	return connect.NewResponse(
		&rpc.StatResponse{
			Entry: &rpc.EntryInfo{
				Name: fileInfo.Name(),
				Type: getEntryType(fileInfo),
				Path: path,
			},
		},
	), nil
}
