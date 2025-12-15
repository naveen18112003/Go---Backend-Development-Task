package service

import (
	"context"
	"errors"
	"time"

	"user-age-api/db/sqlc"
	"user-age-api/internal/models"
	"user-age-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

var (
	// ErrNotFound indicates an entity was not found.
	ErrNotFound = errors.New("user not found")
)

// UserService contains business logic for users.
type UserService struct {
	repo     repository.UserRepository
	validate *validator.Validate
}

// NewUserService creates a new service.
func NewUserService(repo repository.UserRepository, validate *validator.Validate) *UserService {
	return &UserService{repo: repo, validate: validate}
}

// CreateUser validates input, persists, and returns a response.
func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return models.UserResponse{}, err
	}

	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.CreateUser(ctx, db.CreateUserParams{
		Name: req.Name,
		Dob:  dob,
	})
	if err != nil {
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Format("2006-01-02"),
	}, nil
}

// GetUser fetches a user by id and adds computed age.
func (s *UserService) GetUser(ctx context.Context, id int32) (models.UserResponse, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserResponse{}, ErrNotFound
		}
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Format("2006-01-02"),
		Age:  models.CalculateAge(user.Dob, time.Now()),
	}, nil
}

// UpdateUser updates attributes and returns the updated user.
func (s *UserService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return models.UserResponse{}, err
	}

	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.UpdateUser(ctx, db.UpdateUserParams{
		ID:   id,
		Name: req.Name,
		Dob:  dob,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserResponse{}, ErrNotFound
		}
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Format("2006-01-02"),
	}, nil
}

// DeleteUser removes a user by id.
func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	if err := s.repo.DeleteUser(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// ListUsers returns paginated users with computed ages.
func (s *UserService) ListUsers(ctx context.Context, limit, offset int32) ([]models.UserResponse, error) {
	users, err := s.repo.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	resp := make([]models.UserResponse, 0, len(users))
	now := time.Now()
	for _, u := range users {
		resp = append(resp, models.UserResponse{
			ID:   u.ID,
			Name: u.Name,
			Dob:  u.Dob.Format("2006-01-02"),
			Age:  models.CalculateAge(u.Dob, now),
		})
	}
	return resp, nil
}



