package handler

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Artimus100/mcp-server-go/internal/protocol"
	"github.com/Artimus100/mcp-server-go/internal/state"
	"github.com/Artimus100/mcp-server-go/internal/utils"
)

// Connection represents a client connection to the MCP server
type Connection struct {
	id         string
	conn       net.Conn
	store      *state.ContextStore
	logger     *utils.Logger
	closeChan  chan struct{}
	closedOnce sync.Once
}

// Server handles incoming TCP connections
type Server struct {
	port        int
	listener    net.Listener
	store       *state.ContextStore
	logger      *utils.Logger
	connections map[string]*Connection
	mu          sync.RWMutex
	closeChan   chan struct{}
}

// NewServer creates a new MCP server
func NewServer(port int, store *state.ContextStore, logger *utils.Logger) *Server {
	return &Server{
		port:        port,
		store:       store,
		logger:      logger,
		connections: make(map[string]*Connection),
		closeChan:   make(chan struct{}),
	}
}

// Start begins listening for connections
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", addr, err)
	}
	s.listener = listener

	go s.acceptConnections()
	return nil
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown() error {
	close(s.closeChan)

	// Close listener
	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			return err
		}
	}

	// Clean up connections
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, conn := range s.connections {
		s.logger.Info("Closing connection %s", id)
		conn.Close()
		delete(s.connections, id)
	}

	return nil
}

// acceptConnections handles incoming connections
func (s *Server) acceptConnections() {
	for {
		select {
		case <-s.closeChan:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.closeChan:
					return
				default:
					s.logger.Error("Error accepting connection: %v", err)
					continue
				}
			}

			// Create new connection
			connID := fmt.Sprintf("%s-%d", conn.RemoteAddr().String(), time.Now().UnixNano())
			c := &Connection{
				id:        connID,
				conn:      conn,
				store:     s.store,
				logger:    s.logger.WithPrefix(fmt.Sprintf("conn[%s]", connID)),
				closeChan: make(chan struct{}),
			}

			// Add to connections map
			s.mu.Lock()
			s.connections[connID] = c
			s.mu.Unlock()

			// Handle connection in goroutine
			go c.Handle()
		}
	}
}

// BroadcastMessage sends a message to all connected clients
// TODO: Implement message broadcasting logic
func (s *Server) BroadcastMessage(msg protocol.Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for range s.connections {
		// TODO: Format message according to protocol and send
	}
}

// Handle processes incoming messages from a client
func (c *Connection) Handle() {
	defer c.Close()

	c.logger.Info("New connection established")

	reader := bufio.NewReader(c.conn)

	for {
		select {
		case <-c.closeChan:
			return
		default:
			// Set read deadline
			err := c.conn.SetReadDeadline(time.Now().Add(time.Minute))
			if err != nil {
				c.logger.Error("Failed to set read deadline: %v", err)
				return
			}

			// Read line from connection
			line, err := reader.ReadString('\n')
			if err != nil {
				c.logger.Error("Error reading from connection: %v", err)
				return
			}

			// Parse message
			msg, err := protocol.Parse(line)
			if err != nil {
				c.logger.Error("Failed to parse message: %v", err)
				continue
			}

			// Process message
			c.handleMessage(msg)
		}
	}
}

// handleMessage processes a parsed message
func (c *Connection) handleMessage(msg protocol.Message) {
	c.logger.Info("Received message: %s", msg.String())

	switch msg.Type {
	case protocol.TypePing:
		// Handle ping message
		c.handlePing(msg)

	case protocol.TypeContext:
		// Handle context update
		c.handleContextUpdate(msg)

	default:
		c.logger.Warning("Unknown message type: %s", msg.Type)
	}
}

// handlePing responds to ping messages
func (c *Connection) handlePing(msg protocol.Message) {
	// TODO: Implement proper ping response
	// For now, just log it
	c.logger.Info("Ping received with params: %v", msg.Params)

	// Prepare response
	response := protocol.NewMessage(protocol.TypePong, map[string]string{
		"time": fmt.Sprintf("%d", time.Now().Unix()),
	})

	// Send response
	c.Send(response)
}

// handleContextUpdate processes context updates
func (c *Connection) handleContextUpdate(msg protocol.Message) {
	// Log context update
	c.logger.Info("Context update received with params: %v", msg.Params)

	// Update context in the store
	// TODO: Add proper context update logic
	for key, value := range msg.Params {
		c.store.Set(c.id, key, value)
	}

	// Respond with acknowledgment
	response := protocol.NewMessage(protocol.TypeAck, map[string]string{
		"status": "ok",
	})

	c.Send(response)
}

// Send transmits a message to the client
func (c *Connection) Send(msg protocol.Message) {
	formatted := msg.Format()
	_, err := c.conn.Write([]byte(formatted + "\n"))
	if err != nil {
		c.logger.Error("Failed to send message: %v", err)
	}
}

// Close terminates the connection
func (c *Connection) Close() {
	c.closedOnce.Do(func() {
		close(c.closeChan)
		c.conn.Close()
		c.logger.Info("Connection closed")
	})
}
