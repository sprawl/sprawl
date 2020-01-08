package p2p

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/pb"
)

// Stream is a single streaming connection between two peers
type Stream struct {
	stream network.Stream
	input  chan pb.WireMessage
	output chan pb.WireMessage
}

func (p2p *P2p) handleStream(stream network.Stream) {
	if p2p.Logger != nil {
		p2p.Logger.Info("New stream opened")
	}
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	go p2p.receiveStream(rw)
	go p2p.writeToStream(rw)
}

func (p2p *P2p) receiveStream(rw *bufio.ReadWriter) {
	for {
		data, err := rw.ReadBytes('\n')
		if err != nil {
			p2p.Logger.Error(err)
		}
		if string(data) == "" {
			return
		}
	}
}

func (p2p *P2p) writeToStream(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}
		_, err = rw.Write(sendData)
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

func (stream *Stream) listenToStream(p2p *P2p, quitSignal chan bool) {
	go func(ctx context.Context) {
		for {
			select {
			case msg := <-stream.input:
				data := msg.GetData()

				if p2p.Receiver != nil {
					err := p2p.Receiver.Receive(data)
					if !errors.IsEmpty(err) {
						if p2p.Logger != nil {
							p2p.Logger.Error(errors.E(errors.Op("Receive data"), err))
						}
					}
				} else {
					if p2p.Logger != nil {
						p2p.Logger.Warn("Receiver not registered with p2p, not parsing any incoming data!")
					}
				}

				select {
				case quit := <-quitSignal: //Delete subscription
					if quit {
						return
					}
				default:
				}
			}
		}
	}(p2p.ctx)
}

// OpenStream opens a stream with another Sprawl peer
func (p2p *P2p) OpenStream(peerIDString string) error {
	p2p.host.SetStreamHandler(networkID, p2p.handleStream)
	peerID, err := peer.IDFromString(peerIDString)
	stream, err := p2p.host.NewStream(p2p.ctx, peerID, networkID)
	if err != nil {
		p2p.Logger.Errorf("Stream open failed: %s", err)
	} else {
		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
		go p2p.writeToStream(rw)
		go p2p.receiveStream(rw)
		p2p.Logger.Debugf("Stream opened with %s", peerID)
	}
	return err
}
