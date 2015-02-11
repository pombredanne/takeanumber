package server

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	// Test initialization.
	s := New(13331)

	port := s.NetPort()

	if port != ":13331" {
		t.Error("NetPort is wrong, saw: %v", port)
	}

	// Queue shouldn't exist, but should spring to life.
	if _, ok := s.Queues["test_queue"]; ok {
		t.Error("Somehow the test queue already exists. That's not right.")
	}

	q := s.GetQueue("test_queue")

	if q.Len() != 0 {
		t.Error("Queue already has items?")
	}

	if _, ok := s.Queues["test_queue"]; !ok {
		t.Error("The test queue didn't get retained.")
	}

	// Test formatting of responses.
	res_string := s.FormatResponse("A test string")

	if res_string != "+A test string\r\n" {
		t.Error("String didn't format right, saw: ", res_string)
	}

	res_integer := s.FormatResponse(251)

	if res_integer != ":251\r\n" {
		t.Error("Integer didn't format right, saw: ", res_integer)
	}

	res_error := s.FormatResponse(errors.New("Stuff broke"))

	if res_error != "-ERR Stuff broke\r\n" {
		t.Error("Error didn't format right, saw: ", res_error)
	}

	// LEN command
	if s.HandleLen("LEN test_queue") != ":0\r\n" {
		t.Error("Test queue already has items in it, got: ", s.HandleLen("LEN test_queue"))
	}

	// ADD command
	id := s.HandleAdd("ADD test_queue 3 Hello")

	if !strings.HasPrefix(id, "+") {
		t.Error("Add didn't work, got: ", id)
	}

	new_len := s.HandleLen("LEN test_queue")

	if new_len != ":1\r\n" {
		t.Error("Length is wrong, got: ", new_len)
	}

	// RESERVE command
	resp := s.HandleReserve("RESERVE test_queue")
	bits := strings.SplitN(resp, " ", 2)
	id = strings.TrimPrefix(bits[0], "+")
	body := bits[1]

	if body != "Hello\r\n" {
		t.Error("Incorrect body returned, got: ", body)
	}

	// RETRY command
	resp = s.HandleRetry(fmt.Sprintf("RETRY test_queue %v", id))

	if resp != "+OK\r\n" {
		t.Error("Retry failed, got: ", resp)
	}

	// DONE command
	resp = s.HandleDone(fmt.Sprintf("DONE test_queue %v", id))

	if resp != "+OK\r\n" {
		t.Error("Done failed, got: ", resp)
	}
}
