package repositories

import (
	"context"

	"github.com/blockchain-network/internal/models"
)

type NodeRepository interface {
	Create(ctx context.Context, node *models.Node) error
	// GetByID(ctx context.Context, id string) (*models.Node, error)
	GetRandom(ctx context.Context, r int) ([]*models.Node, error)
	GetAll(ctx context.Context) ([]*models.Node, error)
}
