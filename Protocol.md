# takeanumber Protocol

`takeanumber` uses a plain-text protocol when communicating over a TCP socket.
This protocol tries to use Redis' [RESP](http://redis.io/topics/protocol) as
its basis. Only a portion of this is implemented (strings, integers, errors).

In the examples below, ``C: `` is the client talking, ``S: `` is the server.

## Length

**Request:**

    LEN <queue_name>\r\n

**Response:**

    :<length>\r\n

**Example:**

    // Existing queue
    C: LEN my_queue\r\n
    S: :15\r\n

    // Empty/non-existent queue
    C: LEN nopenopenope\r\n
    S: :0\r\n


## Add

**Request:**

    ADD <queue_name> <retries> <value>\r\n

**Response:**

    +<id>\r\n
    // ...or...
    -ERR <message>\r\n

**Example:**

    // Successful add
    C: ADD my_queue 3 {"thing": 1, "also": "abc"}\r\n
    S: +0269073f-f624-4cf9-8c53-ab3d194137b3\r\n

    // Failed add
    C: ADD nopenopenope 1 \r\n
    S: -ERR No body provided.\r\n

## Reserve

**Request:**

    RESERVE <queue_name>\r\n

**Response:**

    +<id> <body>\r\n
    // ...or...
    :-1\r\n

**Example:**

    // Successful reserve
    C: RESERVE my_queue\r\n
    S: +0269073f-f624-4cf9-8c53-ab3d194137b3 {"thing": 1, "also": "abc"}\r\n

    // Empty queue
    C: RESERVE my_queue\r\n
    S: :-1\r\n

## Retry

**Request:**

    RETRY <queue_name> <id>\r\n

**Response:**

    +OK\r\n
    // ...or...
    -ERR <message>\r\n

**Example:**

    // Can retry
    C: RETRY my_queue 0269073f-f624-4cf9-8c53-ab3d194137b3\r\n
    S: +OK\r\n

    // Out of retries
    C: RETRY my_queue 0269073f-f624-4cf9-8c53-ab3d194137b3\r\n
    S: -ERR No retries remaining.\r\n

## Done

**Request:**

    DONE <queue_name> <id>\r\n

**Response:**

    +OK\r\n
    // ...or...
    -ERR <message>\r\n

**Example:**

    // Successful done
    C: DONE my_queue 0269073f-f624-4cf9-8c53-ab3d194137b3\r\n
    S: +OK\r\n

    // Non-existent ID
    C: DONE nopenopenope 0269073f-ffff-4444-8888-ab3d194137b3
    S: -ERR No such Id.\r\n


## Close

**Request:**

    CLOSE\r\n

**Response:**

    No response, closed TCP connection

**Example:**

    C: CLOSE\r\n
