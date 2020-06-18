package main

import (
	"BlueprintChatApp/trace"
	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
	"log"
	"net/http"
)

type room struct {
	forward chan *message      //this is used to holds incoming messages and forward it to other user
	join    chan *client     //channel use if user wants to join room object
	leave   chan *client     //channel use if user wants to leave room object
	clients map[*client]bool //map who is in the room
	tracer  trace.Tracer     //track activity in the room
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
		tracer: trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <- r.join: //jika ada sesuatu di channel join
			r.clients[client] = true  // menambahkan client ke map
			r.tracer.Trace("New User Joined")
		case client := <- r.leave: //jika ada sesuatu di channel leave
			delete(r.clients, client) // menghapus client dari map
			close(client.send)
			r.tracer.Trace("User Left")
		case msg := <-r.forward:   //jika ada sesuatu di channel forward
			r.tracer.Trace("Message received: ", msg.Message)
			r.tracer.Trace("Time: ", msg.When)
			for client := range r.clients {
				client.send <- msg // send message to all client to all user in clients map
				r.tracer.Trace(" -- sent to client ")
			}
		}
	}
}

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

// upgrade http connection
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP (w http.ResponseWriter, req *http.Request) {
	// get the socket
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	//ngedapetin data dari cookie untuk mengisi userdata
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}
	//membuat objeck client
	client := &client{
		socket  : socket,                               //pass socket object
		send    : make(chan *message, messageBufferSize), //membuat channel sebanyak messageBufferSize : 256
		room	: r,                                    //pass object r
		userData: objx.MustFromBase64(authCookie.Value),
	}
	//pass client to join channel
	r.join <- client
	defer func() { r.leave <- client } () //memastikan client dihapus dari leave channel
	go client.write() //menjalankan write dan read function dengan thread yang berbeda
	client.read() //read function run forever
}