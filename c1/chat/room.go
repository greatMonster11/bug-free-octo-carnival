package main

type room struct {
	// forward is a chanel that holds incomming message
	// that should be forwarded to the other clients
	forward chan []byte
	// join is a chanel for clients wishing to join the room
	join chan *client
	// leave is a chandle for clients wishing to leave the room
	leave chan *client
	// clients hold add curren in the room
	clients map[*client]bool
}
