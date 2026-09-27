package plugin

import (
	"context"

	hplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	pb "github.com/akyriako/o7k/plugin/proto"
)

const PluginName = "o7k"

type GRPCPlugin struct {
	hplugin.NetRPCUnsupportedPlugin

	Impl Plugin
	Host Host
}

func (p *GRPCPlugin) GRPCServer(broker *hplugin.GRPCBroker, server *grpc.Server) error {
	pb.RegisterPluginServer(server, &grpcServer{
		impl:   p.Impl,
		broker: broker,
	})

	return nil
}

func (p *GRPCPlugin) GRPCClient(ctx context.Context, broker *hplugin.GRPCBroker, client *grpc.ClientConn) (interface{}, error) {
	grpcClient := &grpcClient{
		client: pb.NewPluginClient(client),
	}

	if p.Host == nil {
		return grpcClient, nil
	}

	id := broker.NextId()

	go broker.AcceptAndServe(id, func(options []grpc.ServerOption) *grpc.Server {
		server := grpc.NewServer(options...)
		pb.RegisterHostServer(server, &grpcHostServer{
			host: p.Host,
		})
		return server
	})

	_, err := grpcClient.client.Initialize(ctx, &pb.InitializeRequest{
		HostBrokerId: id,
	})
	if err != nil {
		return nil, err
	}

	return grpcClient, nil
}
