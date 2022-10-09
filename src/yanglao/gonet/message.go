package gonet

const (
	contextMessageCount uint32 = 1024
)

type contextMessage struct {
	src   *context
	fname string
	args  []interface{}
	reply chan []interface{}
}
