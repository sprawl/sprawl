package interfaces

// Receiver receives and parses all Wiremessages from p2p
type Receiver interface {
	Receive(data []byte) error
}
