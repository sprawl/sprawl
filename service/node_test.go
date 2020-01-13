package service

import (
	"context"
	"testing"

	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/interfaces"
	"github.com/sprawl/sprawl/pb"
)

func TestNodeService(t *testing.T) {
	createNewServerInstance()
	defer p2pInstance.Close()
	defer storage.Close()
	defer conn.Close()
	leaveEveryChannel()

	var nodeService interfaces.NodeService = &NodeService{}
	nodeService.RegisterP2p(p2pInstance)
	pb.RegisterNodeHandlerServer(s, nodeService)

	go func() {
		if err := s.Serve(lis); !errors.IsEmpty(err) {
			log.Fatalf("Server exited with error: %v", err)
		}
		defer s.Stop()
	}()

	var nodeClient pb.NodeHandlerClient = pb.NewNodeHandlerClient(conn)
	peers, _ := nodeClient.GetAllPeers(context.Background(), &pb.Empty{})
	if len(peers.GetPeerIds()) != 0 {
		nodeClient.BlacklistPeer(context.Background(), &pb.PeerId{Id: peers.GetPeerIds()[0]})
	} else {
		nodeClient.BlacklistPeer(context.Background(), &pb.PeerId{Id: "Testi"})
	}
}
