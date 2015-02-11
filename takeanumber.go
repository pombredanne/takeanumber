// Copyright 2015, Daniel Lindsley. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package takeanumber implements a simplistic queue server.

`takeanumber` is an in-memory networked queue server, allowing you to delay
processing by storing messages in one process & consuming those messages in
other processes.

This is especially useful in web applications, allowing you to take non-critical
processing out of the request-response cycle, though there are many other
applications.

`takeanumber` exposes its functionality over a plain-text protocol via a TCP
socket. These messages conform to a subset of the Redis Serialization Protocol
(http://redis.io/topics/protocol), so in theory an existing Redis client may
be able to talk to `takeanumber`.

Starting a server (on localhost, port `13331`):

    $ takeanumber -p 13331

You can then use tools like `telnet` to talk to `takeanumber`. Here's a sample
session:

	$ telnet localhost 13331
	LEN my_queue
	:0
	ADD my_queue 0 Hello, world!
	+bb713fbe-3c82-41c9-94f0-43c499bfac8c
	ADD my_queue 3 {"user_id": 5, "action": "send_welcome_email"}
	+4aaf88df-390b-4a0a-8352-1fe258d94d3d
	LEN my_queue
	:2
	RESERVE my_queue
	+bb713fbe-3c82-41c9-94f0-43c499bfac8c Hello, world!
	DONE my_queue bb713fbe-3c82-41c9-94f0-43c499bfac8c
	+OK
	LEN my_queue
	:1
	RESERVE my_queue
	+4aaf88df-390b-4a0a-8352-1fe258d94d3d {"user_id": 5, "action": "send_welcome_email"}
	RETRY my_queue 4aaf88df-390b-4a0a-8352-1fe258d94d3d
	+OK
	LEN my_queue
	:1
	CLOSE

This session did the following:

	* Checked the length of queue, bringing it into existence if not present
	* Added a "Hello, world!" message to the queue with no retries
	* Added a JSON message to the queue with 3 retries
	* Verified the length of the queue
	* Reserved the first item for processing then marked it as done
	* Reserved the second item, then marked it to be retried
	* Verified the message was in the queue
	* Closed the session

*/
package main

import (
	"flag"
	"fmt"
	"github.com/toastdriven/takeanumber/server"
)

const Version = "1.0.0"

func main() {
	var port int
	flag.IntVar(&port, "p", 13331, "The port to listen on")
	flag.Parse()

	fmt.Printf("takeanumber v%v\n", Version)
	s := server.New(port)

	fmt.Printf("Listening on port %v\n", port)
	s.Run()
}
