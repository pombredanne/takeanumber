package item_test

import (
	"fmt"
	"github.com/toastdriven/takeanumber/item"
)

func ExampleItem() {
	// Create an item with a body & 5 retries.
	i, err := item.New("Hello, world!", 5)

	if err != nil {
		// Bad things happened. Bail out.
	}

	// A unique id (uuid) is created for each Item.
	fmt.Println(i.Id)
	fmt.Println(i.Body)
	fmt.Println(i.InitialRetries)
	// There's also a time created, to track how long something has
	// been in the queue.
	fmt.Println(i.Created)

	// Check the reserve status on the item.
	// This will be false initially.
	fmt.Println(i.IsReserved())

	// Reserving the Item marks it as in-progress, so that other
	// workers can't pick it up.
	i.Reserve()

	// This will now be true.
	fmt.Println(i.IsReserved())

	// You can then release it.
	i.Release()
	// You should try to decrement the retries on releasing, so
	// that the proper number of retries is maintained.
	success := i.DecrRetries()

	if success {
		// Hooray, it worked.
	}

	// Back to false, open for anyone to pick up.
	fmt.Println(i.IsReserved())
}