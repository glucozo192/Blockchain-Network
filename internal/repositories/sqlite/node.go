package sqlite

import (
	"context"

	"github.com/blockchain-network/internal/models"
	"github.com/blockchain-network/internal/repositories"
)

type nodeRepository struct {
	q *models.Queries
}

func NewNodeRepository(q *models.Queries) repositories.NodeRepository {
	return &nodeRepository{q: q}
}

func (r *nodeRepository) GetRandom(ctx context.Context, ra int) ([]*models.Node, error) {
	return r.q.GetRandomNode(ctx, int32(ra))
}

func (r *nodeRepository) Create(ctx context.Context, node *models.Node) error {
	return r.q.CreateNode(ctx, models.CreateNodeParams{
		ID:      node.ID,
		Address: node.Address,
	})
}

func (r *nodeRepository) GetAll(ctx context.Context) ([]*models.Node, error) {
	return r.q.GetAll(ctx)
}
