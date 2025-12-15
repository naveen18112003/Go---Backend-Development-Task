package repository

import (
	"context"

	"user-age-api/db/sqlc"
)

// UserRepository abstracts the data access layer.
type UserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUser(ctx context.Context, id int32) (db.User, error)
	UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
	DeleteUser(ctx context.Context, id int32) error
	ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error)
}

// SQLCUserRepository implements UserRepository using generated sqlc queries.
type SQLCUserRepository struct {
	queries *db.Queries
}

// NewUserRepository returns a repository backed by sqlc queries.
func NewUserRepository(queries *db.Queries) UserRepository {
	return &SQLCUserRepository{queries: queries}
}

func (r *SQLCUserRepository) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return r.queries.CreateUser(ctx, arg)
}

func (r *SQLCUserRepository) GetUser(ctx context.Context, id int32) (db.User, error) {
	return r.queries.GetUser(ctx, id)
}

func (r *SQLCUserRepository) UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error) {
	return r.queries.UpdateUser(ctx, arg)
}

func (r *SQLCUserRepository) DeleteUser(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *SQLCUserRepository) ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	return r.queries.ListUsers(ctx, arg)
}



