package sqlite

import (
	"context"

	"github.com/blockchain-network/internal/models"
	"github.com/blockchain-network/internal/repositories"
)

type blockRepository struct {
	q *models.Queries
}

func NewBlockRepository(q *models.Queries) repositories.BlockRepository {
	return &blockRepository{q: q}
}

func (r *blockRepository) Create(ctx context.Context, block *models.Block) error {
	return r.q.CreateBlock(ctx, models.CreateBlockParams{
		ID:   block.ID,
		Data: block.Data,
	})
}

func (r *blockRepository) GetByID(ctx context.Context, id string) (*models.Block, error) {
	panic("not implemented") // TODO: Implement
}

func (r *blockRepository) GetLatestBlock(ctx context.Context) (*models.Block, error) {
	panic("not implemented") // TODO: Implement
}

func (r *blockRepository) GetAll(ctx context.Context) ([]*models.Block, error) {
	return r.q.GetAllBlock(ctx)
}
