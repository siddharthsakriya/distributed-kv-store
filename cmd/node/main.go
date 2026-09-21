package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/siddharthsakriya/distributed-kv-store/cluster"
	raftpb "github.com/siddharthsakriya/distributed-kv-store/proto/gen/raft/v1"
	"github.com/siddharthsakriya/distributed-kv-store/raft"
	"github.com/siddharthsakriya/distributed-kv-store/server"
	grpctransport "github.com/siddharthsakriya/distributed-kv-store/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	id := flag.String("id", "", "this node's id (must match the config)")
	configPath := flag.String("config", "raft-cluster.yml", "path to cluster config")
	flag.Parse()

	if *id == "" {
		log.Fatal("-id is required")
	}

	cfg, err := cluster.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	myAddr, myClientAddr, peers, err := cfg.View(*id)
	if err != nil {
		log.Fatalf("config view: %v", err)
	}

	peerIDs := make([]string, 0, len(peers))
	for pid := range peers {
		peerIDs = append(peerIDs, pid)
	}
	transport := grpctransport.NewTransport(peers)
	applyCh := make(chan raft.ApplyMsg, 256)
	kv := &server.FakeKV{Store: make(map[string][]byte)}
	node := raft.NewNode(raft.Config{
		ID:        *id,
		Peers:     peerIDs,
		Transport: transport,
		ApplyCh:   applyCh,
	})

	lis, err := net.Listen("tcp", myAddr)
	if err != nil {
		log.Fatalf("listen on %s: %v", myAddr, err)
	}
	grpcSrv := grpc.NewServer()
	raftpb.RegisterRaftServiceServer(grpcSrv, grpctransport.NewServer(node))

	go func() {
		log.Printf("[%s] raft gRPC listening on %s", *id, myAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("serve: %v", err)
		}
	}()

	kvServer := server.NewKVServer(node, kv, applyCh)
	go kvServer.RunApplyLoop()
	go http.ListenAndServe(myClientAddr, server.Handler(kvServer))

	node.Start()
	log.Printf("[%s] node started, peers=%v", *id, peerIDs)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Printf("[%s] shutting down", *id)
	node.Stop()
	grpcSrv.GracefulStop()
	transport.Close()
}
