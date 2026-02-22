package couchbase

import (
	"context"
	"errors"
	"fmt"
	"time"

	cbopentelemetry "github.com/couchbase/gocb-opentelemetry"
	"github.com/couchbase/gocb/v2"
	"github.com/yusirdemir/microservice/internal/domain"
	"github.com/yusirdemir/microservice/internal/repository"
	"github.com/yusirdemir/microservice/pkg/config"
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type couchbaseProductRepository struct {
	cluster    *gocb.Cluster
	bucket     *gocb.Bucket
	collection *gocb.Collection
}

type ProductDocument struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Type      string    `json:"type"`
}

func NewProductRepository(cfg *config.Config) (repository.ProductRepository, error) {
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

	return &couchbaseProductRepository{
		cluster:    cluster,
		bucket:     bucket,
		collection: collection,
	}, nil
}

func (r *couchbaseProductRepository) Create(ctx context.Context, product *domain.Product) error {
	doc := ProductDocument{
		ID:        product.ID,
		UserID:    product.UserID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
		Type:      "product",
	}

	_, err := r.collection.Insert(product.ID, doc, &gocb.InsertOptions{
		Context:    ctx,
		ParentSpan: cbopentelemetry.NewOpenTelemetryRequestSpan(ctx, oteltrace.SpanFromContext(ctx)),
	})
	return err
}

func (r *couchbaseProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	result, err := r.collection.Get(id, &gocb.GetOptions{
		Context:    ctx,
		ParentSpan: cbopentelemetry.NewOpenTelemetryRequestSpan(ctx, oteltrace.SpanFromContext(ctx)),
	})
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	var doc ProductDocument
	err = result.Content(&doc)
	if err != nil {
		return nil, err
	}

	return domain.ReconstituteProduct(
		doc.ID,
		doc.UserID,
		doc.Name,
		doc.Price,
		doc.Stock,
		doc.CreatedAt,
		doc.UpdatedAt,
	), nil
}

func (r *couchbaseProductRepository) FindAllByUserID(ctx context.Context, userID string) ([]*domain.Product, error) {
	query := fmt.Sprintf("SELECT x.* FROM `%s` x WHERE x.type = 'product' AND x.user_id = $1", r.bucket.Name())
	rows, err := r.cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []any{userID},
		Context:              ctx,
		ParentSpan:           cbopentelemetry.NewOpenTelemetryRequestSpan(ctx, oteltrace.SpanFromContext(ctx)),
	})
	if err != nil {
		return nil, err
	}

	var products []*domain.Product
	for rows.Next() {
		var doc ProductDocument
		if err := rows.Row(&doc); err != nil {
			return nil, err
		}
		products = append(products, domain.ReconstituteProduct(
			doc.ID,
			doc.UserID,
			doc.Name,
			doc.Price,
			doc.Stock,
			doc.CreatedAt,
			doc.UpdatedAt,
		))
	}
	return products, nil
}

func (r *couchbaseProductRepository) Update(ctx context.Context, product *domain.Product) error {
	product.UpdatedAt = time.Now()

	doc := ProductDocument{
		ID:        product.ID,
		UserID:    product.UserID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
		Type:      "product",
	}

	_, err := r.collection.Replace(product.ID, doc, &gocb.ReplaceOptions{
		Context:    ctx,
		ParentSpan: cbopentelemetry.NewOpenTelemetryRequestSpan(ctx, oteltrace.SpanFromContext(ctx)),
	})
	return err
}

func (r *couchbaseProductRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.Remove(id, &gocb.RemoveOptions{
		Context:    ctx,
		ParentSpan: cbopentelemetry.NewOpenTelemetryRequestSpan(ctx, oteltrace.SpanFromContext(ctx)),
	})
	return err
}
