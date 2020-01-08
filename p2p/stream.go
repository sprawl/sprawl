package p2p

import (
	"bufio"

	"github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/sprawl/sprawl/errors"
)

// Stream is a single streaming connection between two peers
type Stream struct {
	stream network.Stream
	input  *bufio.Writer
}

func (p2p *P2p) handleStream(stream network.Stream) {
	if p2p.Logger != nil {
		p2p.Logger.Info("New stream opened")
	}
	reader := bufio.NewReader(bufio.NewReader(stream))
	go p2p.receiveStream(reader)
}

func (p2p *P2p) receiveStream(rw *bufio.Reader) error {
	for {
		data, err := rw.ReadBytes('\n')
		if err != nil {
			return errors.E(errors.Op("Reading bytes from stream"), err)
		} else {
			p2p.Logger.Info(data)
			if p2p.Receiver != nil {
				err := p2p.Receiver.Receive(data)
				if !errors.IsEmpty(err) {
					if p2p.Logger != nil {
						return errors.E(errors.Op("Passing data from stream to receiver"), err)
					}
				}
			} else {
				if p2p.Logger != nil {
					p2p.Logger.Warn("Receiver not registered with p2p, not parsing any incoming data!")
				}
			}
		}
		if string(data) == "" {
			return nil
		}
	}
}

func (stream *Stream) writeToStream(data []byte) error {
	_, err := stream.input.Write(data)
	err = stream.input.Flush()
	return err
}

// OpenStream opens a stream with another Sprawl peer
func (p2p *P2p) OpenStream(peerIDString string) error {
	peerID, err := peer.IDFromString(peerIDString)
	stream, err := p2p.host.NewStream(p2p.ctx, peerID, networkID)
	if err != nil {
		p2p.Logger.Errorf("Stream open failed: %s", err)
	} else {
		writer := bufio.NewWriter(bufio.NewWriter(stream))
		p2p.streams[peerIDString] = Stream{stream: stream, input: writer}
		p2p.Logger.Debugf("Stream opened with %s", peerID)
	}
	return err
}

// CloseStream removes and closes a stream
func (p2p *P2p) CloseStream(peerIDString string) error {
	err := p2p.streams[peerIDString].stream.Close()
	delete(p2p.streams, peerIDString)
	return err
}
