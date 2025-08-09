package model

type Devices map[string]Device

// Device represents a single device under a user
type Device struct {
	Name       string `json:"name"`       // device name
	DeviceType string `json:"deviceType"` // device type or model

	Signals map[string]interface{} `json:"signals"` // sensor data, key-value pairs

	Commands []Command  `json:"commands"` // commands sent to and responses from device
	Logs     []LogEntry `json:"logs"`     // historical logs or events
}

type CommandData struct {
	Path string    `json:"path"`
	Data []Command `json:"data"`
}

// Command represents a command with bidirectional info
type Command struct {
	Action      string `json:"action"`                // command action name
	IssuedAt    int64  `json:"issuedAt"`              // when command was issued
	Status      string `json:"status"`                // e.g., "pending", "executing", "done", "error"
	Output      string `json:"output,omitempty"`      // optional result/output from device
	CompletedAt int64  `json:"completedAt,omitempty"` // optional completion timestamp
	ErrorMsg    string `json:"errorMsg,omitempty"`    // optional error description
}

// LogEntry represents one log entry/event from device or app
type LogEntry struct {
	LogID     string `json:"logId"`     // unique log entry ID
	Timestamp int64  `json:"timestamp"` // event time
	Message   string `json:"message"`   // log message text
}
