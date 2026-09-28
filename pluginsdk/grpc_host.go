package pluginsdk

import (
	"context"

	pb "github.com/akyriako/o7k/pluginsdk/proto"
)

type grpcHostClient struct {
	client pb.HostClient
}

func (c *grpcHostClient) Context(ctx context.Context) (Context, error) {
	response, err := c.client.GetContext(ctx, &pb.Empty{})
	if err != nil {
		return Context{}, err
	}

	return Context{
		Generation: response.Generation,
		Cloud:      response.Cloud,
		CloudsPath: response.CloudsPath,
		Region:     response.Region,
	}, nil
}

type grpcHostServer struct {
	pb.UnimplementedHostServer

	host Host
}

func (s *grpcHostServer) GetContext(ctx context.Context, _ *pb.Empty) (*pb.Context, error) {
	current, err := s.host.Context(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.Context{
		Generation: current.Generation,
		Cloud:      current.Cloud,
		CloudsPath: current.CloudsPath,
		Region:     current.Region,
	}, nil
}
