"""
takeanumber
-----------

A client library for talking to the ``takeanumber`` queue server.

Usage::

    >>> from takeanumber import Client
    >>> c = Client()
    >>> c.len('my_queue')
    0
    >>> c.add('my_queue', 'Hello, world!')
    '4488ea2c-197d-4443-9467-74f75317917c'
    >>> c.len('my_queue')
    1
    >>> c.reserve('my_queue')
    ('4488ea2c-197d-4443-9467-74f75317917c', 'Hello, world!')
    >>> c.done('my_queue', '4488ea2c-197d-4443-9467-74f75317917c')
    'OK'
    >>> c.close()

"""
import socket


class TakeANumberError(Exception): pass
class EmptyBodyError(TakeANumberError): pass
class EmptyQueueError(TakeANumberError): pass
class MissingParmetersError(TakeANumberError): pass
class InvalidRetriesError(TakeANumberError): pass
class NoRetriesRemainingError(TakeANumberError): pass


class Client(object):
    def __init__(self, host='127.0.0.1', port=13331, timeout=30):
        self.host = host
        self.port = port
        self.timeout = timeout
        self.sock = None

    def connect(self):
        self.sock = socket.create_connection(
            (self.host, self.port),
            self.timeout
        )

    def _send(self, command):
        sent = self.sock.sendall(command)

        if sent == 0:
            raise RuntimeError("Broken connection")

    def _receive(self):
        chunks = []

        while True:
            chunk = self.sock.recv(4096)

            if chunk == b'':
                raise RuntimeError("Broken connection")

            chunks.append(chunk)

            if chunk.endswith('\r\n'):
                break

        return str(b''.join(chunks))

    def decode(self, resp):
        if resp[0] not in (':', '+', '-', '$', '*'):
            return resp

        resp = resp.rstrip('\r\n')
        ty, clean_resp = resp[0], resp[1:]

        if ty == '+':
            return clean_resp
        elif ty == ':':
            return int(clean_resp)
        elif ty == '-':
            # We need to do more here, because exceptions.
            if 'Missing' in clean_resp and 'parameters' in clean_resp:
                raise MissingParmetersError(clean_resp)
            elif 'Invalid number of retries' in clean_resp:
                raise InvalidRetriesError(clean_resp)
            elif 'No retries' in clean_resp:
                raise NoRetriesRemainingError(clean_resp)
            elif 'No body' in clean_resp:
                raise EmptyBodyError(clean_resp)
            elif 'No items available' in clean_resp:
                raise EmptyQueueError(clean_resp)
            else:
                raise TakeANumberError(clean_resp)

        # We can't (& don't need to) handle the other RESP cases. Just return
        # the raw response for now.
        return resp

    def len(self, queue_name):
        if self.sock is None:
            self.connect()

        command = "LEN {}\r\n".format(queue_name)
        self._send(command)
        return self.decode(self._receive())

    def add(self, queue_name, body, retries=0):
        if self.sock is None:
            self.connect()

        command = "ADD {} {} {}\r\n".format(
            queue_name,
            retries,
            body
        )
        self._send(command)
        return self.decode(self._receive())

    def reserve(self, queue_name):
        if self.sock is None:
            self.connect()

        command = "RESERVE {}\r\n".format(queue_name)
        self._send(command)
        raw_body = self.decode(self._receive())
        ident, body = raw_body.split(' ', 1)
        return ident, body

    def retry(self, queue_name, ident):
        if self.sock is None:
            self.connect()

        command = "RETRY {} {}\r\n".format(queue_name, ident)
        self._send(command)
        return self.decode(self._receive())

    def done(self, queue_name, ident):
        if self.sock is None:
            self.connect()

        command = "DONE {} {}\r\n".format(queue_name, ident)
        self._send(command)
        return self.decode(self._receive())

    def close(self):
        if self.sock is None:
            self.connect()

        command = "CLOSE\r\n"
        self._send(command)
        self.sock.close()
