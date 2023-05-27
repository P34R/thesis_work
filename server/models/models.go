package models

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	PubKey    string `json:"pubKey"`
	Nonce     string `json:"nonce"`
	NonceSign int    `json:"nonceSign"`
}

type Message struct {
	Id     int    `json:"id"`
	ChatId int    `json:"chat"`
	From   int    `json:"from"`
	Mess   string `json:"message"`
	Stamp  int64  `json:"stamp"`
}

type Packet struct {
	Type    int    `json:"type"`
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}
type Socket struct {
	Conn *websocket.Conn
	In   chan Packet
	Out  chan Packet
	Quit chan int
}

func NewSocket(conn *websocket.Conn) *Socket {
	return &Socket{
		Conn: conn,
		In:   make(chan Packet),
		Out:  make(chan Packet),
		Quit: make(chan int),
	}
}

type Connections struct {
	mu    sync.Mutex
	conns map[string]*Socket
}

func NewConnections() *Connections {
	return &Connections{
		mu:    sync.Mutex{},
		conns: make(map[string]*Socket),
	}
}
func (c *Connections) AddSocket(username string, conn *Socket) {
	c.mu.Lock()
	c.conns[username] = conn
	log.Println(c.conns[username].In)
	c.mu.Unlock()
}

func (c *Connections) GetSocket(username string) *Socket {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conns[username]
}
func (c *Connections) IsPresent(username string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.conns[username]; ok {
		return true
	}
	return false
} /*
func (c *Connections) SendToFront(mess Packet) bool{
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.conns[username].Conn.WriteJSON(mess)
	if err != nil {
		log.Println("error in ")
		return false
	}
	return true
}*/
func (c *Connections) SendMessage(mess Packet) bool {
	if c.IsPresent(mess.To) {
		c.mu.Lock()
		c.conns[mess.To].In <- mess
		c.mu.Unlock()
		log.Println(mess.From, " SENT wow ", mess.To)
		return true
	}
	return false
}
func (c *Connections) CloseConn(username string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conns[username].Quit <- 1
	close(c.conns[username].In)
	close(c.conns[username].Quit)
	c.conns[username].Conn.Close()
	delete(c.conns, username)
}
func (c *Connections) CloseAllConns() {
	for k, _ := range c.conns {
		c.CloseConn(k)
	}
}
