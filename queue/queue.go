// Copyright 2015 Daniel Lindsley. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package queue implements a simple FIFO queue.

Example:

	import (
		"fmt"
		"github.com/toastdriven/takeanumber/queue"
	)

	func Whatever() {
		q := queue.New()
		fmt.Println(q.Len())

		// Add an item with zero retries.
		id, err := q.Add("Hello, world!", 0)

		// The id of the item is returned.
		fmt.Println(id)
		fmt.Println(q.Len())

		// Fetch the topmost item from the queue.
		item, err := q.Reserve()
		fmt.Println(item.Body)

		// Mark it as Done.
		success := q.Done(item.Id)
	}

*/
package queue

import (
	"errors"
	"fmt"
	"github.com/toastdriven/takeanumber/item"
	"sync"
)

// An error for when there are no items in the queue.
var EmptyQueue = errors.New("No items available to reserve.")

// The Queue itself.
type Queue struct {
	Items []*item.Item
	lock *sync.Mutex
}

// Adds an item to the end of the queue.
//
// Accepts a body (string) & the number of times it can be retried (integer).
// This will create a new Item & push it onto the end of the queue.
//
// The Item's Id (uuid string) is returned.
func (q *Queue) Add(body string, retries int) (string, error) {
	i, err := item.New(body, retries)

	if err != nil {
		return "", err
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	q.Items = append(q.Items, i)
	return i.Id, nil
}

// Reserves an item from the front of the queue.
//
// This will fetch the first *non-reserved* Item from the queue, mark it as
// reserved & return it. If all the items are already reserved or there is
// nothing in the queue, an EmptyQueue error is returned.
func (q *Queue) Reserve() (*item.Item, error) {
	var i *item.Item
	found := false

	q.lock.Lock()
	defer q.lock.Unlock()

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

// Marks an item as completed.
//
// Accepts the Id (string) of the item to be marked done.
// If found, it will be removed from the queue.
//
// Returns whether the item was successfully removed or not (bool).
func (q *Queue) Done(id string) bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	for offset, current := range q.Items {
		if current.Id == id {
			q.Items = append(q.Items[:offset], q.Items[offset+1:]...)
			return true
		}
	}

	return false
}

// Marks an item to be retried.
//
// Accepts the Id (string) of the item to be retried. The item will become
// unreserved, its retry count will be decremented & it maintain its place
// early in the queue to be picked up again.
//
// If the retry count is zero, the item will be removed & this will return
// false, since the item will disappear from the queue.
//
// Returns whether the item was successfully marked to be retried (bool).
func (q *Queue) Retry(id string) bool {
	q.lock.Lock()
	defer q.lock.Unlock()

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

// Returns the length of *unreserved* items in the queue.
//
// This count can be used to determine if there are any items to be processed.
//
// Returns a count of items (integer).
func (q *Queue) Len() int {
	length := 0

	for _, current := range q.Items {
		if !current.IsReserved() {
			length++
		}
	}

	return length
}

// New creates a new Queue instance.
func New() *Queue {
	items := []*item.Item{}
	lock := &sync.Mutex{}
	return &Queue{items, lock}
}
