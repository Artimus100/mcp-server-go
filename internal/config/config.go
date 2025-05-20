package config

// Server configuration constants
const (
	// DefaultPort is the default port for the MCP server
	DefaultPort = 8080

	// MaxMessageSize is the maximum allowed size for incoming messages in bytes
	MaxMessageSize = 4096

	// ReadTimeout is the default read timeout in seconds for client connections
	ReadTimeout = 60

	// WriteTimeout is the default write timeout in seconds for client connections
	WriteTimeout = 10

	// IdleTimeout is the default idle timeout in seconds for client connections
	IdleTimeout = 300

	// MaxConnections is the maximum number of simultaneous connections
	MaxConnections = 1000
)

// Protocol configuration constants
const (
	// ProtocolVersion is the current version of the MCP protocol
	ProtocolVersion = "1.0"

	// MessageDelimiter is the character used to separate messages
	MessageDelimiter = '\n'

	// TypeParamSeparator is the character that separates message type from parameters
	TypeParamSeparator = ':'

	// ParamPairSeparator is the character that separates parameter pairs
	ParamPairSeparator = ';'

	// ParamKeyValueSeparator is the character that separates parameter keys from values
	ParamKeyValueSeparator = '='
)

// Logging constants
const (
	// LogTimeFormat is the time format used in log messages
	LogTimeFormat = "2006-01-02 15:04:05.000"

	// LogLevel defines the default log level
	LogLevel = "info"
)

// TODO: Add other application-wide constants as needed
// TODO: Consider making these configurable via environment variables or config file
