package p2p

import (
	"bufio"

	"github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/interfaces"
)

// Stream is a single streaming connection between two peers
type Stream struct {
	stream     network.Stream
	remotePeer peer.ID
	input      *bufio.Writer
	output     *bufio.Reader
}

func (p2p *P2p) handleStream(buf network.Stream) {
	p2p.Logger.Debugf("New stream opened with %s", buf.Conn().RemotePeer())
	reader := bufio.NewReader(bufio.NewReader(buf))
	remotePeer := buf.Conn().RemotePeer()
	stream := &Stream{stream: buf, output: reader, remotePeer: remotePeer}
	go func() {
		stream.receiveStream(reader, p2p.Receiver)
		stream.stream.Close()
	}()
}

func (stream *Stream) receiveStream(reader *bufio.Reader, receiver interfaces.Receiver) error {
	data := []byte{}
	for {
		line, _ := reader.ReadByte()
		data = append(data, line)
		if reader.Buffered() == 0 {
			err := receiver.Receive(data, stream.remotePeer)
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Passing data from stream to receiver"), err)
			}
			return nil
		}
	}
}

// WriteToStream writes data as bytes to specified stream
func (stream *Stream) WriteToStream(data []byte) error {
	_, err := stream.input.Write(data)
	if err != nil {
		return errors.E(errors.Op("Write to stream"), err)
	}
	err = stream.input.Flush()
	return errors.E(errors.Op("Flush the stream"), err)
}

// OpenStream opens a stream with another Sprawl peer
func (p2p *P2p) OpenStream(peerID peer.ID) (interfaces.Stream, error) {
	stream, err := p2p.host.NewStream(p2p.ctx, peerID, networkID)
	var newStream *Stream
	if err != nil {
		p2p.Logger.Errorf("Stream open failed with peer %s on network %s: %s", peerID, networkID, err)
	} else {
		writer := bufio.NewWriter(bufio.NewWriter(stream))
		newStream = &Stream{stream: stream, input: writer, remotePeer: peerID}
		p2p.streams[peerID.String()] = newStream
	}
	return newStream, err
}

// CloseStream removes and closes a stream
func (p2p *P2p) CloseStream(peerID peer.ID) error {
	err := p2p.streams[peerID.String()].stream.Close()
	delete(p2p.streams, peerID.String())
	return err
}
