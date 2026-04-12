package logger

type ProtocolType string

const (
	GRPC ProtocolType = "grpc"
	HTTP ProtocolType = "http"
)

type Protocol struct {
	Type    ProtocolType `json:"type"`
	Version string       `json:"version"`
}

type RequestLogField struct {
	Path        string              `json:"path"`
	Method      string              `json:"method"`
	QueryParams map[string][]string `json:"query_params"`
	IP          string              `json:"ip"`
	UserAgent   string              `json:"user_agent"`
	Protocol    Protocol            `json:"protocol"`
}
