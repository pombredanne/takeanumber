package takeanumber

import (
	"code.google.com/p/go-uuid/uuid"
	"strings"
	"time"
)

type EmptyBody struct {
	What string
}

func (err *EmptyBody) Error() string {
	return err.What
}

type Item struct {
	Id               string
	Body             string
	InitialRetries   int
	RemainingRetries int
	Reserved         bool
	Created          time.Time
}

func (i *Item) DecrRetries() bool {
	if !i.ShouldRetry() {
		return false
	}

	i.RemainingRetries--
	return true
}

func (i *Item) ShouldRetry() bool {
	return i.RemainingRetries > 0
}

func (i *Item) IsReserved() bool {
	return i.Reserved
}

func (i *Item) Reserve() {
	i.Reserved = true
}

func (i *Item) Release() {
	i.Reserved = false
}

func NewItem(body string, retries int) (Item, error) {
	if len(strings.TrimSpace(body)) <= 0 {
		return Item{}, &EmptyBody{"No body provided."}
	}

	id := uuid.New()
	created := time.Now()
	return Item{id, body, retries, retries, false, created}, nil
}
