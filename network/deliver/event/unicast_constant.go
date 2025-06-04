package event

const (
	JSON_CONTENT_TYPE = "application/json"
)

type DefaultEndPoint struct {
	protocol string
	address  string
	port     uint16
	path     string
}
