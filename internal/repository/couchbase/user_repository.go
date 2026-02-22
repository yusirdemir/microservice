package couchbase

import (
	"context"
	"errors"
	"time"

	cbopentelemetry "github.com/couchbase/gocb-opentelemetry"
	"github.com/couchbase/gocb/v2"
	"github.com/yusirdemir/microservice/internal/domain"
	"github.com/yusirdemir/microservice/internal/repository"
	"github.com/yusirdemir/microservice/pkg/config"
	"go.opentelemetry.io/otel"
)

type couchbaseUserRepository struct {
	cluster    *gocb.Cluster
	bucket     *gocb.Bucket
	collection *gocb.Collection
}

type UserDocument struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Type      string    `json:"type"`
}

func NewUserRepository(cfg *config.Config) (repository.UserRepository, error) {
	cluster, err := gocb.Connect(cfg.Database.Host, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: cfg.Database.Username,
			Password: cfg.Database.Password,
		},
		Tracer: cbopentelemetry.NewOpenTelemetryRequestTracer(otel.GetTracerProvider()),
	})
	if err != nil {
		return nil, err
	}

	bucket := cluster.Bucket(cfg.Database.Bucket)
	err = bucket.WaitUntilReady(30*time.Second, nil)
	if err != nil {
		return nil, err
	}

	collection := bucket.DefaultCollection()

	return &couchbaseUserRepository{
		cluster:    cluster,
		bucket:     bucket,
		collection: collection,
	}, nil
}

func (r *couchbaseUserRepository) Create(ctx context.Context, user *domain.User) error {
	doc := UserDocument{
		ID:        user.ID(),
		Name:      user.Name(),
		Email:     user.Email(),
		Password:  user.Password(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
		Type:      "user",
	}

	_, err := r.collection.Insert(user.ID(), doc, &gocb.InsertOptions{
		Context: ctx,
	})
	return err
}

func (r *couchbaseUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	result, err := r.collection.Get(id, &gocb.GetOptions{
		Context: ctx,
	})
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	var doc UserDocument
	err = result.Content(&doc)
	if err != nil {
		return nil, err
	}

	return domain.Reconstitute(
		doc.ID,
		doc.Name,
		doc.Email,
		doc.Password,
		doc.CreatedAt,
		doc.UpdatedAt,
	), nil
}

func (r *couchbaseUserRepository) Update(ctx context.Context, user *domain.User) error {
	doc := UserDocument{
		ID:        user.ID(),
		Name:      user.Name(),
		Email:     user.Email(),
		Password:  user.Password(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
		Type:      "user",
	}

	_, err := r.collection.Replace(user.ID(), doc, &gocb.ReplaceOptions{
		Context: ctx,
	})
	return err
}

func (r *couchbaseUserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.Remove(id, &gocb.RemoveOptions{
		Context: ctx,
	})
	return err
}
