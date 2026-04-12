package logger

type Event struct {
	Message         string                 `json:"message"`
	ContextBoundary string                 `json:"context_boundary"`
	Context         map[string]interface{} `json:"context,omitempty"`
}
