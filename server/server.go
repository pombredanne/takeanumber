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

type Server struct {
	Port int
	Queues map[string]*queue.Queue
}

func (s *Server) NetPort() string {
	return fmt.Sprintf(":%v", s.Port)
}

func (s *Server) GetQueue(name string) *queue.Queue {
	if q, ok := s.Queues[name]; ok {
		return q
	}

	s.Queues[name] = queue.New()
	return s.Queues[name]
}

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

func (s *Server) HandleLen(command string) string {
	bits := strings.SplitN(command, " ", 2)

	if len(bits) != 2 {
		return s.FormatResponse(errors.New("Missing LEN parameters."))
	}

	q := s.GetQueue(bits[1])
	return s.FormatResponse(q.Len())
}

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

func New(port int) *Server {
	qs := map[string]*queue.Queue{}
	return &Server{port, qs}
}
