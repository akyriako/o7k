package plugin

import (
	"context"
	"fmt"

	pb "github.com/akyriako/o7k/plugin/proto"
	hplugin "github.com/hashicorp/go-plugin"
)

type grpcServer struct {
	pb.UnimplementedPluginServer

	impl   Plugin
	broker *hplugin.GRPCBroker
}

func (s *grpcServer) Initialize(ctx context.Context, request *pb.InitializeRequest) (*pb.Empty, error) {
	conn, err := s.broker.Dial(request.HostBrokerId)
	if err != nil {
		return nil, fmt.Errorf("dialing host: %w", err)
	}

	s.impl.SetHost(&grpcHostClient{
		client: pb.NewHostClient(conn),
	})

	return &pb.Empty{}, nil
}

func (s *grpcServer) GetMetadata(context.Context, *pb.Empty) (*pb.Metadata, error) {
	metadata := s.impl.Metadata()

	return &pb.Metadata{
		Name:    metadata.Name,
		Version: metadata.Version,
	}, nil
}

func (s *grpcServer) GetResources(context.Context, *pb.Empty) (*pb.Resources, error) {
	resources := s.impl.Resources()
	result := make([]*pb.Resource, 0, len(resources))

	for _, resource := range resources {
		columns := resource.Columns()
		pbColumns := make([]*pb.Column, 0, len(columns))

		for _, column := range columns {
			pbColumns = append(pbColumns, &pb.Column{
				Key:      column.Key,
				Title:    column.Title,
				MinWidth: int32(column.MinWidth),
				Flex:     int32(column.Flex),
			})
		}

		commands := resource.Commands()
		pbCommands := make([]*pb.Command, 0, len(commands))

		for _, command := range commands {
			pbCommands = append(pbCommands, &pb.Command{
				Key:         command.Key,
				Description: command.Description,
				Default:     command.Default,
			})
		}

		result = append(result, &pb.Resource{
			Service:  resource.Service(),
			Kind:     resource.Kind(),
			Title:    resource.Title(),
			Aliases:  resource.Aliases(),
			Columns:  pbColumns,
			Commands: pbCommands,
		})
	}

	return &pb.Resources{Resources: result}, nil
}

func (s *grpcServer) List(ctx context.Context, request *pb.ListRequest) (*pb.ListResponse, error) {
	resource, err := s.resource(request.Resource)
	if err != nil {
		return nil, err
	}

	rows, err := resource.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*pb.Row, 0, len(rows))

	for _, row := range rows {
		result = append(result, &pb.Row{
			Id:     row.ID,
			Fields: row.Fields,
		})
	}

	return &pb.ListResponse{Rows: result}, nil
}

func (s *grpcServer) Execute(ctx context.Context, request *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	resource, err := s.resource(request.Resource)
	if err != nil {
		return nil, err
	}

	result, err := resource.Execute(ctx, Command{
		Key:         request.Command.Key,
		Description: request.Command.Description,
		Default:     request.Command.Default,
	}, Row{
		ID:     request.Row.Id,
		Fields: request.Row.Fields,
	})
	if err != nil {
		return nil, err
	}

	response := &pb.Result{}

	if result.Details != nil {
		response.Details = &pb.Details{
			Id:      result.Details.ID,
			Content: result.Details.Content,
		}
	}

	if result.Navigate != nil {
		response.Navigate = &pb.Navigate{
			Resource: result.Navigate.Resource,
			Id:       result.Navigate.ID,
			Field:    result.Navigate.Field,
			Value:    result.Navigate.Value,
			Scope:    result.Navigate.Scope,
		}
	}

	return &pb.ExecuteResponse{Result: response}, nil
}

func (s *grpcServer) resource(kind string) (Resource, error) {
	for _, resource := range s.impl.Resources() {
		if resource.Kind() == kind {
			return resource, nil
		}
	}

	return nil, fmt.Errorf("resource %q not found", kind)
}
