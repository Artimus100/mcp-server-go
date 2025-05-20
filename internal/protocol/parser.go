package protocol

import (
	"fmt"
	"strings"
)

// Message types
const (
	TypePing    = "PING"
	TypePong    = "PONG"
	TypeContext = "CONTEXT"
	TypeAck     = "ACK"
	TypeError   = "ERROR"
	// TODO: Add more message types as needed
)

// Message represents a parsed MCP protocol message
type Message struct {
	Type   string
	Params map[string]string
}

// NewMessage creates a new message with the given type and parameters
func NewMessage(msgType string, params map[string]string) Message {
	if params == nil {
		params = make(map[string]string)
	}

	return Message{
		Type:   msgType,
		Params: params,
	}
}

// Parse converts a raw message string into a Message struct
// Format: TYPE:key=value;key2=value2
func Parse(raw string) (Message, error) {
	// Trim whitespace and any trailing newlines
	raw = strings.TrimSpace(raw)

	// Check for empty message
	if raw == "" {
		return Message{}, fmt.Errorf("empty message")
	}

	// Split message into type and parameters
	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 {
		return Message{}, fmt.Errorf("invalid message format: missing type separator")
	}

	msgType := strings.TrimSpace(parts[0])
	if msgType == "" {
		return Message{}, fmt.Errorf("missing message type")
	}

	// Parse parameters
	params := make(map[string]string)
	if parts[1] != "" {
		paramPairs := strings.Split(parts[1], ";")
		for _, pair := range paramPairs {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) != 2 {
				return Message{}, fmt.Errorf("invalid parameter format: %s", pair)
			}

			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])

			if key == "" {
				return Message{}, fmt.Errorf("empty parameter key")
			}

			params[key] = value
		}
	}

	return Message{
		Type:   msgType,
		Params: params,
	}, nil
}

// Format converts a Message struct back into a protocol string
func (m Message) Format() string {
	var params []string

	for key, value := range m.Params {
		params = append(params, fmt.Sprintf("%s=%s", key, value))
	}

	paramStr := strings.Join(params, ";")
	return fmt.Sprintf("%s:%s", m.Type, paramStr)
}

// String returns a string representation of the message for logging
func (m Message) String() string {
	return fmt.Sprintf("Message{Type: %s, Params: %v}", m.Type, m.Params)
}

// ValidateMessageType checks if a message type is valid
func ValidateMessageType(msgType string) bool {
	validTypes := map[string]bool{
		TypePing:    true,
		TypePong:    true,
		TypeContext: true,
		TypeAck:     true,
		TypeError:   true,
		// Add other valid types here
	}

	_, valid := validTypes[msgType]
	return valid
}

// TODO: Add more protocol helpers as needed, such as:
// - Message validation
// - Special message constructors for common message types
// - Protocol versioning
