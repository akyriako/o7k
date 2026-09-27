package plugin

import (
	"context"

	pb "github.com/akyriako/o7k/plugin/proto"
)

type grpcClient struct {
	client pb.PluginClient
}

func (c *grpcClient) Metadata() Metadata {
	response, err := c.client.GetMetadata(context.Background(), &pb.Empty{})
	if err != nil {
		return Metadata{}
	}

	return Metadata{
		Name:    response.Name,
		Version: response.Version,
	}
}

func (c *grpcClient) Resources() []Resource {
	response, err := c.client.GetResources(context.Background(), &pb.Empty{})
	if err != nil {
		return nil
	}

	resources := make([]Resource, 0, len(response.Resources))

	for _, resource := range response.Resources {
		resources = append(resources, &grpcResource{
			client:   c.client,
			resource: resource,
		})
	}

	return resources
}

type grpcResource struct {
	client   pb.PluginClient
	resource *pb.Resource
}

func (r *grpcResource) Service() string {
	return r.resource.Service
}

func (r *grpcResource) Kind() string {
	return r.resource.Kind
}

func (r *grpcResource) Title() string {
	return r.resource.Title
}

func (r *grpcResource) Aliases() []string {
	return r.resource.Aliases
}

func (r *grpcResource) Columns() []Column {
	columns := make([]Column, 0, len(r.resource.Columns))

	for _, column := range r.resource.Columns {
		columns = append(columns, Column{
			Key:      column.Key,
			Title:    column.Title,
			MinWidth: int(column.MinWidth),
			Flex:     int(column.Flex),
		})
	}

	return columns
}

func (r *grpcResource) Commands() []Command {
	commands := make([]Command, 0, len(r.resource.Commands))

	for _, command := range r.resource.Commands {
		commands = append(commands, Command{
			Key:         command.Key,
			Description: command.Description,
			Default:     command.Default,
		})
	}

	return commands
}

func (r *grpcResource) List(ctx context.Context) ([]Row, error) {
	response, err := r.client.List(ctx, &pb.ListRequest{
		Resource: r.Kind(),
	})
	if err != nil {
		return nil, err
	}

	rows := make([]Row, 0, len(response.Rows))

	for _, row := range response.Rows {
		rows = append(rows, Row{
			ID:     row.Id,
			Fields: row.Fields,
		})
	}

	return rows, nil
}

func (r *grpcResource) Execute(ctx context.Context, command Command, row Row) (Result, error) {
	response, err := r.client.Execute(ctx, &pb.ExecuteRequest{
		Resource: r.Kind(),
		Command: &pb.Command{
			Key:         command.Key,
			Description: command.Description,
			Default:     command.Default,
		},
		Row: &pb.Row{
			Id:     row.ID,
			Fields: row.Fields,
		},
	})
	if err != nil {
		return Result{}, err
	}

	result := Result{}

	if response.Result.Details != nil {
		result.Details = &Details{
			ID:      response.Result.Details.Id,
			Content: response.Result.Details.Content,
		}
	}

	if response.Result.Navigate != nil {
		result.Navigate = &Navigate{
			Resource: response.Result.Navigate.Resource,
			ID:       response.Result.Navigate.Id,
			Field:    response.Result.Navigate.Field,
			Value:    response.Result.Navigate.Value,
			Scope:    response.Result.Navigate.Scope,
		}
	}

	return result, nil
}
