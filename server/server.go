// Copyright 2015 Daniel Lindsley. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package server implements a TCP server that listens for queue commands.

For a complete description of the available commands, responses & errors, see
the included Protocol.md document that is included with `takeanumber`.

Example:

	import (
		"fmt"
		"github.com/toastdriven/takeanumber/server"
	)

	func main() {
		port := 13331

		// Create a Server.
		s := server.New(port)

		// Run the server.
		s.Run()
	}

*/
package server

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"github.com/toastdriven/takeanumber/queue"
)

// The Server itself.
type Server struct {
	Port int
	Queues map[string]*queue.Queue
}

// Returns a string version of the port (with preceding colon) for use with
// `net.Listen(...)`.
func (s *Server) NetPort() string {
	return fmt.Sprintf(":%v", s.Port)
}

// Fetches & returns a Queue by name.
//
// Accepts the name (string) of the Queue. If the queue does not already exist,
// a new queue will be created.
//
// Returns the Queue.
func (s *Server) GetQueue(name string) *queue.Queue {
	if q, ok := s.Queues[name]; ok {
		return q
	}

	s.Queues[name] = queue.New()
	return s.Queues[name]
}

// Formats a response for returning to the client.
//
// Accepts the response (interface{}), which may be a string, integer or error.
// Based on the type of the response, this will create a RESP encoded string.
//
// Returns the formatted response (string).
func (s *Server) FormatResponse(resp interface{}) string {
	var toFormat string

	switch resp.(type) {
	case string:
		toFormat = "+%s\r\n"
	case int, int8, int16, int32, int64:
		toFormat = ":%d\r\n"
	case error:
		toFormat = "-ERR %s\r\n"
	default:
		toFormat = "-ERR %s\r\n"
	}

	return fmt.Sprintf(toFormat, resp)
}

// Handles the LEN command.
//
// The command should include the name of the queue. The queue will be fetched
// & an integer count of the length of the queue will be returned.
//
// Returns a formatted integer string.
//
// Command Format:
//
//	LEN <queue_name>\r\n
//
// Response Format:
//
//	:<integer>\r\n
func (s *Server) HandleLen(command string) string {
	bits := strings.SplitN(command, " ", 2)

	if len(bits) != 2 {
		return s.FormatResponse(errors.New("Missing LEN parameters."))
	}

	q := s.GetQueue(bits[1])
	return s.FormatResponse(q.Len())
}

// Handles the ADD command.
//
// The command should include the name of the queue, the number of times it can
// be retried & the message body. The queue will be fetched
// & a new Item with the data will be placed at the end of the queue.
//
// Warning: Bodies may *not* be empty, nor can there be any bare newlines in
// the body.
//
// Returns a formatted string of the new item's Id.
//
// Command Format:
//
//	ADD <queue_name> <retries> <value>\r\n
//
// Response Format:
//
//	+<id>\r\n
func (s *Server) HandleAdd(command string) string {
	bits := strings.SplitN(command, " ", 4)

	if len(bits) != 4 {
		return s.FormatResponse(errors.New("Missing ADD parameters."))
	}

	q := s.GetQueue(bits[1])
	retries, err := strconv.Atoi(bits[2])

	if err != nil {
		return s.FormatResponse(errors.New("Invalid number of retries."))
	}

	id, err := q.Add(bits[3], retries)

	if err != nil {
		return s.FormatResponse(err)
	}

	return s.FormatResponse(id)
}

// Handles the RESERVE command.
//
// The command should include the name of the queue. The queue will be fetched
// & an item will be reserved off the front of the queue.
//
// Returns a formatted string of the item Id & body.
//
// Command Format:
//
//	RESERVE <queue_name>\r\n
//
// Response Format:
//
//	+<id> <body>\r\n
func (s *Server) HandleReserve(command string) string {
	bits := strings.SplitN(command, " ", 2)

	if len(bits) != 2 {
		return s.FormatResponse(errors.New("Missing RESERVE parameters."))
	}

	q := s.GetQueue(bits[1])
	i, err := q.Reserve()

	if err != nil {
		return s.FormatResponse(err)
	}

	resp := fmt.Sprintf("%s %s", i.Id, i.Body)
	return s.FormatResponse(resp)
}

// Handles the RETRY command.
//
// The command should include the name of the queue & the Id of the item to
// retry. The queue will be fetched & the item will marked to be retried.
//
// Returns a formatted "OK" string.
//
// Command Format:
//
//	RETRY <queue_name> <id>\r\n
//
// Response Format:
//
//	+OK\r\n
func (s *Server) HandleRetry(command string) string {
	bits := strings.SplitN(command, " ", 3)

	if len(bits) != 3 {
		return s.FormatResponse(errors.New("Missing RETRY parameters."))
	}

	q := s.GetQueue(bits[1])
	id := bits[2]
	res := q.Retry(id)

	if !res {
		return s.FormatResponse(errors.New("No retries remaining."))
	}

	return s.FormatResponse("OK")
}

// Handles the DONE command.
//
// The command should include the name of the queue & the Id of the item to
// be marked as done. The queue will be fetched & the item will be removed.
//
// Returns a formatted "OK" string.
//
// Command Format:
//
//	DONE <queue_name> <id>\r\n
//
// Response Format:
//
//	+OK\r\n
func (s *Server) HandleDone(command string) string {
	bits := strings.SplitN(command, " ", 3)

	if len(bits) != 3 {
		return s.FormatResponse(errors.New("Missing DONE parameters."))
	}

	q := s.GetQueue(bits[1])
	id := bits[2]
	res := q.Done(id)

	if !res {
		return s.FormatResponse(errors.New("No such Id."))
	}

	return s.FormatResponse("OK")
}

// Handles any command(s) sent by the client.
//
// The processing of each type of command is done by the other Handle*
// functions on the Server instance. This simply handles the
// reading/dispatching/writing flow.
func (s *Server) Handle(c net.Conn) {
	scanner := bufio.NewScanner(c)

	for scanner.Scan() {
		var resp string
		command := strings.TrimSpace(scanner.Text())

		switch {
		case strings.HasPrefix(command, "LEN "):
			resp = s.HandleLen(command)
		case strings.HasPrefix(command, "ADD "):
			resp = s.HandleAdd(command)
		case strings.HasPrefix(command, "RESERVE "):
			resp = s.HandleReserve(command)
		case strings.HasPrefix(command, "RETRY "):
			resp = s.HandleRetry(command)
		case strings.HasPrefix(command, "DONE "):
			resp = s.HandleDone(command)
		case strings.HasPrefix(command, "CLOSE"):
			c.Close()
			return
		default:
			resp = s.FormatResponse(errors.New("Unrecognized command."))
		}

		respBytes := []byte(resp)
		c.Write(respBytes)
	}
}

// Runs the server.
//
// This will use the preconfigured port, start listening on it & will spawn
// goroutines for each connection made. This will run forever & must be
// manually terminated.
func (s *Server) Run() {
	l, err := net.Listen("tcp", s.NetPort())

	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go s.Handle(conn)
	}
}

// New creates a new Server instance.
func New(port int) *Server {
	qs := map[string]*queue.Queue{}
	return &Server{port, qs}
}
