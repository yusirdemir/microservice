package service

import (
	"context"

	"github.com/yusirdemir/microservice/internal/domain"
	"github.com/yusirdemir/microservice/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var userTracer = otel.Tracer("microservice/service/user")

type UserService interface {
	CreateUser(ctx context.Context, name, email, password string) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, id, name, email string) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	ctx, span := userTracer.Start(ctx, "UserService.CreateUser")
	defer span.End()

	user, err := domain.NewUser(name, email, password)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.String("app.user.id", user.ID()))

	if err := s.repo.Create(ctx, user); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	ctx, span := userTracer.Start(ctx, "UserService.GetUser")
	defer span.End()

	span.SetAttributes(attribute.String("app.user.id", id))

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id, name, email string) (*domain.User, error) {
	ctx, span := userTracer.Start(ctx, "UserService.UpdateUser")
	defer span.End()

	span.SetAttributes(attribute.String("app.user.id", id))

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if name != "" {
		if err := user.UpdateName(name); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}

	if err := s.repo.Update(ctx, user); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	ctx, span := userTracer.Start(ctx, "UserService.DeleteUser")
	defer span.End()

	span.SetAttributes(attribute.String("app.user.id", id))

	if err := s.repo.Delete(ctx, id); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
