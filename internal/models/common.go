package models

type Response struct {
	Err      error `json:"err,omitempty" bson:"err,omitempty" db:"err"`
	IsAccept bool  `json:"is_accept,omitempty" bson:"is_accept,omitempty" db:"is_accept"`
}

type Request struct {
	Data    []int  `json:"data,omitempty" bson:"data,omitempty" db:"data"`
	Event   Event  `json:"event,omitempty" bson:"event,omitempty" db:"event"`
	BlockID string `json:"block_id,omitempty" bson:"block_id,omitempty" db:"block_id"`
}

type Event int

const (
	PingEvent Event = iota + 1
	ValidateEvent
)
