package takeanumber

import (
	"fmt"
	"container/list"
)


type EmptyQueue struct {
	What string
}

func (err *EmptyQueue) Error() string {
	return err.What
}


type UnknownElement struct {
	What string
	Value interface{}
}

func (err *UnknownElement) Error() string {
	return fmt.Sprintf("%v: %v", err.What, err.Value)
}


type Queue struct {
	Items list.List
}

func (q *Queue) Add(body string, retries int) (string, error) {
	i, err := NewItem(body, retries)

	if err != nil {
		return "", err
	}

	q.Items.PushBack(i)
	return i.Id, nil
}

func (q *Queue) Reserve() (*Item, error) {
	var i *Item
	found := false


	for current := q.Items.Front(); current != nil; current = current.Next() {
		if i, ok := current.Value.(Item); ok {
			i = (Item)(i)

			if !i.IsReserved() {
				found = true
				break
			}
		}
	}

	if !found {
		return nil, &EmptyQueue{"No items available to reserve."}
	}

	i.Reserve()
	return i, nil
}

func (q *Queue) Done(id string) bool {
	for current := q.Items.Front(); current != nil; current = current.Next() {
		if i, ok := current.Value.(Item); ok {
			// i = (Item)(i)

			if i.Id == id {
				q.Items.Remove(current)
				return true
			}
		}
	}

	return false
}

func (q *Queue) Retry(i Item) bool {
	success := i.DecrRetries()

	if !success {
		return false
	}

	i.Release()
	q.Items.PushFront(i)
	return true
}

func (q *Queue) Len() int {
	return q.Items.Len()
}


func NewQueue() *Queue {
	// items := list.New{}
	return &Queue{} //items}
}