# file-share
This repository contains my take on a simple file share between clients, there is the server that takes requests to share a file and a client side that can send files to channels and recieve files once they have registered to one.

# Client
This Go programm can subscribe to channels to receive files and send files to an specific channel (right now can't send files to more than 1 channel).

## Usage
```
go run ./client.go send <channel>
go run ./client.go register <channel>
```

The channel can be any string (must not contain spaces)

# Server
This go programm manages the channels created when a register request is issued, the files sent to the channels must not surpass a size of 32MB.

## Usage
```
go run ./server start
```

# TODO
There are many improvements to be made to this personal implementation of a file share service:
- Users implementation.
- Change the way clients receive files since opening a port and receiving any request to store a file is dangerous(still testing the usage of net instead of net/http to handle this using Conn interfaces)
- Implement a way to receive requests to unsubscribe to a channel.
- Handling of error produced when a client subscribed to a channel stops listening for requests.
