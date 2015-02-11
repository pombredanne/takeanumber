package queue_test

import (
	"fmt"
	"github.com/toastdriven/takeanumber/queue"
)

func ExampleQueue() {
	q := queue.New()
	fmt.Println(q.Len())

	// Add an item with zero retries.
	id, err := q.Add("Hello, world!", 0)

	if err != nil {
		// Bad things happened. Bail out.
	}

	// The id of the item is returned.
	fmt.Println(id)
	fmt.Println(q.Len())

	// Fetch the topmost item from the queue.
	item, err := q.Reserve()
	fmt.Println(item.Body)

	// Mark it as Done.
	success := q.Done(item.Id)

	if success {
		// Huzzah, time to celebrate!
	}
}