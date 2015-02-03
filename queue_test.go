package takeanumber

import "testing"

func TestQueue(t *testing.T) {
	// Test initialization.
	q := NewQueue()

	if q.Len() != 0 {
		t.Error("Queue length wasn't zeroed out, saw:", q.Len())
	}

	id_1, err := q.Add("test 1", 2)
	// fmt.Print("Id #1:", id_1)

	if err != nil {
		t.Error("Saw error:", err)
	}

	if id_1 == "" {
		t.Error("No valid Id for #1 returned!")
	}

	id_2, err := q.Add("test 2", 3)

	if err != nil {
		t.Error("Saw error:", err)
	}

	if id_2 == "" {
		t.Error("No valid Id for #2 returned!")
	}

	id_3, err := q.Add("test 3", 0)

	if err != nil {
		t.Error("Saw error:", err)
	}

	if id_3 == "" {
		t.Error("No valid Id for #3 returned!")
	}

	if q.Len() != 3 {
		t.Error("Queue length is wrong, expected 3, got:", q.Len())
	}

	reserve_1, err := q.Reserve()

	if err != nil {
		t.Error("Failed to reserve #1:", err)
	}

	if reserve_1.Id != id_1 {
		t.Error("Got the wrong item back first, saw:", reserve_1.Body)
	}
}
