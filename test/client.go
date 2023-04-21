package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/blockchain-network/internal/models"
	"github.com/google/uuid"
)

func main() {
	blocks := [][]int{
		{1, 2, 3, 4, 5},
		{1, 3, 4, 5, 6},
		{1, 6, 4, 3, 2},
	}
	addr := ":9003"
	log.Println("len", len(blocks))
	for i := 0; i < 3; i++ {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Fatalf("can not connect :", addr)
		}
		req := &models.Request{
			Data:    blocks[i],
			Event:   models.PingEvent,
			BlockID: uuid.NewString(),
		}
		b, _ := json.Marshal(req)
		conn.Write(b)
		conn.Close()
	}
}
