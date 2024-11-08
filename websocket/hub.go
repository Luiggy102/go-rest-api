package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

// upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// constructor
func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

// handle websocket
func (hub *Hub) HandleWebSocket(w http.ResponseWriter,
	r *http.Request) {
	// upgrade the http to websocket
	socker, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error upgrading connection",
			http.StatusInternalServerError)
		return
	}
	// upgrade succesful
	// new client
	c := NewClient(hub, socker)
	// register the new client
	hub.register <- c

	// the cliente begin to write to the websocket
	go c.Write()
}
