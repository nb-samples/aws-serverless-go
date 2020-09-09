package response

type (
	// Error response body
	Error struct {
		Code    int         `json:"code,omitempty"`    // code (use either code or reason)
		Reason  string      `json:"reason,omitempty"`  // reason (use either code or reason)
		Message string      `json:"message"`           // human-readable description
		Details interface{} `json:"details,omitempty"` // public details
		Errors  []Error     `json:"errors,omitempty"`  // detailed errors
		Meta    interface{} `json:"-"`                 // private meta data (excluded from marshalling)
	}
)

func (e *Error) Error() string {
	return e.Message
}
