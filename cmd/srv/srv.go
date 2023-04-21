package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/blockchain-network/configs"
	"github.com/blockchain-network/internal/deliveries/tcp"
	"github.com/blockchain-network/internal/domains"
	"github.com/blockchain-network/internal/models"
	"github.com/blockchain-network/internal/repositories"
	"github.com/blockchain-network/internal/repositories/sqlite"
	"github.com/blockchain-network/pkg/tcp_server"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type server struct {
	//* server
	tcpServer *tcp_server.TcpServer

	//* deliveries
	blockchainDelivery tcp.BlockchainDelivery

	//* domains
	blockchainDomain domains.BlockchainDomain

	//* repositories
	blockRepo  repositories.BlockRepository
	nodeRepo   repositories.NodeRepository
	markerRepo repositories.MarkerRepository

	db *sql.DB

	//* config
	config *configs.Config

	processors []processor
	factories  []factory

	logFile *os.File
}

type processor interface {
	Init(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type factory interface {
	Connect(ctx context.Context) error
	Stop(ctx context.Context) error
}

func (s *server) loadDatabaseClients(ctx context.Context) error {
	dbFilePath := fmt.Sprintf("data/%s.db", s.config.NodeID)
	os.Remove(dbFilePath) // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	file, err := os.Create(dbFilePath) // Create SQLite file
	if err != nil {
		logger.Fatal(err.Error())
	}
	file.Close()
	db, _ := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return err
	}

	s.db = db
	return nil

}

func (s *server) migrate() error {
	driver, err := sqlite3.WithInstance(s.db, &sqlite3.Config{})
	if err != nil {
		logger.Fatal(err)
	}
	migrationPath := "file://./database/migrations"
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"sqlite3", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		logger.Println(err)
		return err
	}
	return nil
}

func (s *server) loadRepositories() error {

	q := models.New(s.db)
	s.blockRepo = sqlite.NewBlockRepository(q)
	s.markerRepo = sqlite.NewMarkerRepository(q)
	s.nodeRepo = sqlite.NewNodeRepository(q)
	return nil
}

func (s *server) loadDomains() error {
	s.blockchainDomain = domains.NewBlockchainDomain(
		s.nodeRepo,
		s.blockRepo,
		s.markerRepo,
		s.config,
		logger,
	)
	return nil
}

func (s *server) loadDeliveries() error {
	s.blockchainDelivery = tcp.NewBlockchainDelivery(s.blockchainDomain, logger)
	return nil
}

func (s *server) loadConfig(ctx context.Context) error {
	s.config = &configs.Config{}
	s.config.SampleSize, _ = strconv.Atoi(os.Getenv("SAMPLE_SIZE"))
	s.config.QuorumSize, _ = strconv.Atoi(os.Getenv("QUORUM_SIZE"))
	s.config.DecisionThreshHold, _ = strconv.Atoi(os.Getenv("DECISION_THRESHOLD"))
	s.config.Tcp.Port = os.Getenv("PORT")
	s.config.NodeID = os.Getenv("NODE_ID")
	return nil
}

func (s *server) loadServers(ctx context.Context) error {
	s.tcpServer = &tcp_server.TcpServer{
		Addr: configs.ConnectionAddr{
			Port: s.config.Tcp.Port,
		},
		Handlers: map[models.Event]func(ctx context.Context, req *models.Request) (*models.Response, error){
			models.PingEvent:     s.blockchainDelivery.RetrievePingEvent,
			models.ValidateEvent: s.blockchainDelivery.ValidateData,
		},
		Logger: logger,
	}
	s.processors = append(s.processors, s.tcpServer)
	return nil
}

func (s *server) loadLogger(ctx context.Context) error {
	filePath := fmt.Sprintf("./logs/%s.log", s.config.NodeID)
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logger = log.New(logFile, "app ", log.LstdFlags)
	return nil
}
