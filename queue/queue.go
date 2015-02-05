package queue

import (
	"errors"
	"fmt"
	"github.com/toastdriven/takeanumber/item"
)

var EmptyQueue = errors.New("No items available to reserve.")

type UnknownElement struct {
	What  string
	Value interface{}
}

func (err *UnknownElement) Error() string {
	return fmt.Sprintf("%v: %v", err.What, err.Value)
}

type Queue struct {
	Items []*item.Item
}

func (q *Queue) Add(body string, retries int) (string, error) {
	i, err := item.New(body, retries)

	if err != nil {
		return "", err
	}

	q.Items = append(q.Items, i)
	return i.Id, nil
}

func (q *Queue) Reserve() (*item.Item, error) {
	var i *item.Item
	found := false

	for _, current := range q.Items {
		if !current.IsReserved() {
			i = current
			found = true
			break
		}
	}

	if !found {
		return &item.Item{}, EmptyQueue
	}

	i.Reserve()
	return i, nil
}

func (q *Queue) Done(id string) bool {
	for offset, current := range q.Items {
		if current.Id == id {
			q.Items = append(q.Items[:offset], q.Items[offset+1:]...)
			return true
		}
	}

	return false
}

func (q *Queue) Retry(id string) bool {
	for offset, current := range q.Items {
		if current.Id == id {
			success := current.DecrRetries()

			if !success {
				q.Items = append(q.Items[:offset], q.Items[offset+1:]...)
				return false
			}

			current.Release()
			return true
		}
	}

	return false
}

func (q *Queue) Len() int {
	length := 0

	for _, current := range q.Items {
		if !current.IsReserved() {
			length++
		}
	}

	return length
}

func New() *Queue {
	items := []*item.Item{}
	return &Queue{items}
}
