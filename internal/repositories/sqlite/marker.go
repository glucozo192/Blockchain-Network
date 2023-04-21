package sqlite

import (
	"context"

	"github.com/blockchain-network/internal/models"
	"github.com/blockchain-network/internal/repositories"
)

type markerRepository struct {
	q *models.Queries
}

func (r *markerRepository) MarkBlock(ctx context.Context, blockID string) error {
	panic("not implemented") // TODO: Implement
}

func (r *markerRepository) GetByBlockID(ctx context.Context, blockID string) (*models.Marker, error) {
	panic("not implemented") // TODO: Implement
}

func NewMarkerRepository(q *models.Queries) repositories.MarkerRepository {
	return &markerRepository{q: q}
}
