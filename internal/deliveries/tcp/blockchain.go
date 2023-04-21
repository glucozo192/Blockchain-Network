package tcp

import (
	"context"
	"encoding/json"
	"log"

	"github.com/blockchain-network/internal/domains"
	"github.com/blockchain-network/internal/models"
)

type BlockchainDelivery interface {
	RetrievePingEvent(ctx context.Context, req *models.Request) (*models.Response, error)
	ValidateData(ctx context.Context, req *models.Request) (*models.Response, error)
}

type blockchainDelivery struct {
	blockchainDomain domains.BlockchainDomain
	logger           *log.Logger
}

func NewBlockchainDelivery(blockchainDomain domains.BlockchainDomain, logger *log.Logger) BlockchainDelivery {
	return &blockchainDelivery{
		blockchainDomain: blockchainDomain,
		logger:           logger,
	}
}

func (d *blockchainDelivery) RetrievePingEvent(ctx context.Context, req *models.Request) (*models.Response, error) {
	b, _ := json.Marshal(req)
	d.logger.Println(string(b))
	// go d.blockchainDomain.PingNeighbourNodes(ctx, req)
	if err := d.blockchainDomain.SnowBall(ctx, req); err != nil {
		return nil, err
	}
	return nil, nil
}
func (d *blockchainDelivery) ValidateData(ctx context.Context, req *models.Request) (*models.Response, error) {
	b, _ := json.Marshal(req)
	d.logger.Println(string(b))
	err := d.blockchainDomain.Validate(ctx, req.Data)
	if err != nil {
		return nil, err
	}
	return &models.Response{
		IsAccept: true,
	}, nil
}
