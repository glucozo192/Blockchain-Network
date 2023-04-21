package domains

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/blockchain-network/configs"
	"github.com/blockchain-network/internal/models"
	"github.com/blockchain-network/internal/repositories"
)

type BlockchainDomain interface {
	Validate(ctx context.Context, data []int) error
	SnowBall(ctx context.Context, req *models.Request) error
	PingNeighbourNodes(ctx context.Context, req *models.Request)
}

type blockchainDomain struct {
	nodeRepo   repositories.NodeRepository
	blockRepo  repositories.BlockRepository
	markerRepo repositories.MarkerRepository
	configs    *configs.Config
	logger     *log.Logger
}

func NewBlockchainDomain(
	nodeRepo repositories.NodeRepository,
	blockRepo repositories.BlockRepository,
	markerRepo repositories.MarkerRepository,
	configs *configs.Config,
	logger *log.Logger,
) BlockchainDomain {
	return &blockchainDomain{nodeRepo: nodeRepo,
		blockRepo:  blockRepo,
		markerRepo: markerRepo,
		configs:    configs,
		logger:     logger,
	}
}

func (d *blockchainDomain) checkValid(ctx context.Context, source []int32, target []int) error {
	m := make(map[int]int)
	for pos, num := range source {
		m[int(num)] = pos
	}
	pointer := -1
	for _, num := range target {
		if _, ok := m[num]; ok {
			if m[num] < pointer {
				return fmt.Errorf("wrong position")
			}
			pointer = m[num]
		}
	}
	return nil
}

func (d *blockchainDomain) Validate(ctx context.Context, data []int) error {
	blocks, err := d.blockRepo.GetAll(ctx)
	if err != nil {
		return err
	}
	for _, block := range blocks {
		if err := d.checkValid(ctx, block.Data, data); err != nil {
			return fmt.Errorf("data is not valid")
		}
	}
	return nil
}

func (d *blockchainDomain) PingNeighbourNodes(ctx context.Context, req *models.Request) {
	_, err := d.markerRepo.GetByBlockID(ctx, req.BlockID)
	if err != nil {
		d.logger.Printf("unable to get block: %v", err)
		return
	}
	if err := d.markerRepo.MarkBlock(ctx, req.BlockID); err != nil {
		d.logger.Printf("unable to mark block: %v", err)
		return
	}
	nodes, err := d.nodeRepo.GetAll(ctx)
	if err != nil {
		d.logger.Printf("unable to get nodes: %v", err)
		return
	}
	var wg sync.WaitGroup
	for _, node := range nodes {
		if node.Address == ":"+d.configs.Tcp.Port {
			continue
		}
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				d.logger.Printf("unable to connect %s: %v", addr, err)
				return
			}
			defer conn.Close()

			b, _ := json.Marshal(&req)
			if _, err := conn.Write(b); err != nil {
				d.logger.Printf("unable to write data : %v", err)
				return
			}
		}(node.Address)
	}
	wg.Wait()
}

func (d *blockchainDomain) SnowBall(ctx context.Context, req *models.Request) error {
	preference := false
	consecutiveSuccesses := 0
	wg := sync.WaitGroup{}

	startTime := time.Now()

	for {
		d.logger.Printf("consecutiveSuccesses=%d\n", consecutiveSuccesses)
		if time.Since(startTime) > 10*time.Second {
			d.logger.Println("more than 10s")
			break
		}
		chosenNodes, err := d.nodeRepo.GetRandom(ctx, d.configs.SampleSize)
		if err != nil {
			continue
		}
		totalAccept := int32(0)
		totalDecline := int32(0)
		for _, node := range chosenNodes {
			if node.Address == ":"+d.configs.Tcp.Port {
				continue
			}
			wg.Add(1)
			go func(ctx context.Context, addr string) {
				defer wg.Done()

				conn, err := net.Dial("tcp", addr)
				if err != nil {
					d.logger.Printf("unable to connect %s: %v", addr, err)
					return
				}
				defer conn.Close()
				var (
					request models.Request
					resp    models.Response
				)
				request = *req
				request.Event = models.ValidateEvent
				b, err := json.Marshal(&request)
				if err != nil {
					d.logger.Printf("unable to marshal : %v", err)
					return

				}
				if _, err := conn.Write(b); err != nil {
					d.logger.Printf("unable to write data : %v", err)
					return
				}

				data := make([]byte, 1024)
				n, err := conn.Read(data)
				if err != nil {
					d.logger.Printf("unable to reading :%v\n", err)
					return
				}

				if err := json.Unmarshal(data[:n], &resp); err != nil {
					d.logger.Printf("unable to unmarshal :%v\n", err)
					return
				}
				if resp.Err != nil {
					d.logger.Printf("unable to unmarshal :%v\n", err)
					return
				}

				if resp.IsAccept {
					atomic.AddInt32(&totalAccept, 1)
				} else {
					atomic.AddInt32(&totalDecline, 1)
				}
			}(ctx, node.Address)
		}
		wg.Wait()
		time.Sleep(1 * time.Second)
		d.logger.Printf("totalAccept=%d totalDecline=%d\n", totalAccept, totalDecline)
		if int(totalAccept) < d.configs.QuorumSize && int(totalDecline) < d.configs.QuorumSize {
			consecutiveSuccesses = 0
			continue
		}
		currenntPref := totalAccept > totalDecline
		if currenntPref != preference {
			preference = currenntPref
			consecutiveSuccesses = 1
		} else {
			consecutiveSuccesses++
		}

		if consecutiveSuccesses >= d.configs.DecisionThreshHold {
			d.decide(ctx, req)
			break
		}
	}
	return nil
}

func (d *blockchainDomain) decide(ctx context.Context, req *models.Request) {
	var nums []int32
	for _, num := range req.Data {
		nums = append(nums, int32(num))
	}
	d.logger.Println(nums)
	d.blockRepo.Create(ctx, &models.Block{
		ID:   req.BlockID,
		Data: nums,
	})
}
