package tui

// MsgUpdateStatus updates the status in UI
type MsgUpdateStatus struct {
	Latency string
	Addr    string
}

// MsgUpdateUserID updates the user ID in UI
type MsgUpdateUserID string

// MsgUpdateImage updates the image in UI
type MsgUpdateImage string

// Push pushes a message to the TUI
func Push(msg interface{}) {
	// TODO: Implement message push logic
}

// Start starts the TUI
func Start() {
	// TODO: Implement TUI start logic
}
