package httputility

const (
	ErrNoEmpty = "fields cannot be empty"
)

type HTTPError struct {
	Message string `json:"message"`
}
