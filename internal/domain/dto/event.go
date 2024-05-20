package dto

import "time"

type EventListElem struct {
	Kind      string    `json:"kind"`
	Extra     string    `json:"extra"`
	Timestamp time.Time `json:"timestamp"`
}

type EventListOutput []EventListElem
