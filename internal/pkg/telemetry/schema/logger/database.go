package logger

type DatabaseLogField struct {
	NumberOfCalls    int `json:"number_of_calls"`
	NumberOfFailures int `json:"number_of_errors"`
}
