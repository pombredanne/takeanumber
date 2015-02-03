package takeanumber

import (
	"fmt"
	"container/list"
	"github.com/toastdriven/takeanumber/item"
)


type EmptyQueue struct {
	What string
}

func (err *EmptyQueue) Error() string {
	return err.What
}


type Queue struct {
	Items list.List
}

func (q *Queue) Add(body string, retries uint) string {
	item := item.New{body, retries}
	q.Items.PushBack(item)
	return item.Id
}

func (q *Queue) Reserve() (item.Item, error) {
	var i item.Item
	found := false

	for current := q.Items.Next() {
		if !current.IsReserved() {
			i := current
			found = true
			break
		}
	}

	if !found {
		return nil, &EmptyQueue{"No items available to reserve."}
	}

	i.Reserve()
	return &i
}

func (q *Queue) Done(id string) bool {
	for i := q.Items.Next() {
		if i.Id == id {
			q.Items.Remove(i)
			return true
		}
	}

	return false
}

func (q *Queue) Retry(i item.Item) bool {
	success := i.DecrRetries()

	if !success {
		return false
	}

	i.Release()
	q.Items.PushFront(i)
	return true
}


func New() Queue {
	// items := list.New{}
	return &Queue{} //items}
}
