package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"gitlab.silkrode.com.tw/team_golang/kbc/km/storage/graph/generated"
	"gitlab.silkrode.com.tw/team_golang/kbc/km/storage/graph/model"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	fmt.Println("CreateTodo")
	return nil, nil
}

func (r *mutationResolver) SingleUpload(ctx context.Context, file graphql.Upload) (bool, error) {
	fmt.Println("Filename", file.Filename)
	fmt.Println("ContentType~~~~", file.ContentType)
	fmt.Println("read~~~~", file.File.Read)
	fmt.Println("File~~~~", file.File)
	return true, nil
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
