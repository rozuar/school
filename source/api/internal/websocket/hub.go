package websocket

import (
	"encoding/json"
	"log"

	fiberws "github.com/gofiber/websocket/v2"
)

// Hub mantiene el conjunto de clientes activos y difunde mensajes a los clientes
type Hub struct {
	// Clientes registrados
	clients map[*Client]bool

	// Canal para registrar clientes
	register chan *Client

	// Canal para desregistrar clientes
	unregister chan *Client

	// Canal para difundir mensajes a todos los clientes
	broadcast chan []byte
}

// NewHub crea un nuevo Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
	}
}

// Run inicia el hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected. Total clients: %d", len(h.clients))
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// BroadcastMessage envia un mensaje a todos los clientes conectados
func (h *Hub) BroadcastMessage(message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	h.broadcast <- data
	return nil
}

// Client representa un cliente WebSocket
type Client struct {
	hub *Hub

	// La conexion WebSocket
	conn *fiberws.Conn

	// Canal con los mensajes a enviar
	send chan []byte
}

// readPump bombea mensajes desde la conexion WebSocket al hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if fiberws.IsUnexpectedCloseError(err, fiberws.CloseGoingAway, fiberws.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

// writePump bombea mensajes desde el hub a la conexion WebSocket
func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(fiberws.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(fiberws.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Agregar mensajes pendientes al mensaje actual
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// HandleConnections maneja las conexiones WebSocket (Fiber).
func HandleConnections(hub *Hub, conn *fiberws.Conn) {
	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
