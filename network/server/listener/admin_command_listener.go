package listener

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/andantan/vote-blockchain-server/config"
	pb "github.com/andantan/vote-blockchain-server/network/gRPC/admin_commands_message"
	"github.com/andantan/vote-blockchain-server/util"
	"google.golang.org/grpc"
)

type AdminCommandListenerOption struct {
	network string
	port    uint16
}

func NewAdminCommandListenerOption(network string, port int) *AdminCommandListenerOption {
	return &AdminCommandListenerOption{
		network: network,
		port:    uint16(port),
	}
}

type AdminCommandListener struct {
	grpcServer *grpc.Server
	listener   net.Listener
	option     *AdminCommandListenerOption
	exitCh     chan uint8
}

func NewAdminCommandListener(exitCh chan uint8) *AdminCommandListener {
	connectionGrpcCommandListenerNetwork := config.GetEnvVar("CONNECTION_GRPC_ADMIN_COMMANDS_LISTENER_NETWORK")
	connectionGrpcCommandListenerPort := config.GetIntEnvVar("CONNECTION_GRPC_ADMIN_COMMANDS_LISTENER_PORT")

	opts := NewAdminCommandListenerOption(connectionGrpcCommandListenerNetwork, connectionGrpcCommandListenerPort)
	grpcServer := grpc.NewServer()

	l4Service := &L4CommandsServiceImpl{}

	pb.RegisterL4CommandsServer(grpcServer, l4Service)

	return &AdminCommandListener{
		grpcServer: grpcServer,
		option:     opts,
		exitCh:     exitCh,
	}

}

type L4CommandsServiceImpl struct {
	pb.UnimplementedL4CommandsServer
}

func (s *L4CommandsServiceImpl) CheckHealth(ctx context.Context, req *pb.L4HealthCheckRequest) (*pb.L4HealthCheckResponse, error) {
	currentListeningPorts := []uint32{
		uint32(config.GetIntEnvVar("CONNECTION_GRPC_PROPOSAL_LISTENER_PORT")),
		uint32(config.GetIntEnvVar("CONNECTION_GRPC_SUBMIT_LISTENER_PORT")),
		uint32(config.GetIntEnvVar("CONNECTION_GRPC_ADMIN_COMMANDS_LISTENER_PORT")),
		uint32(config.GetIntEnvVar("CONNECTION_REST_EXPLORER_LISTENER_PORT")),
	}

	serverIP := getOutboundIP()

	return &pb.L4HealthCheckResponse{
		Connected: true,
		Status:    "OK",
		Pong:      req.Ping,
		Ip:        serverIP,
		Ports:     currentListeningPorts,
	}, nil
}

func (acl *AdminCommandListener) Start(wg *sync.WaitGroup) {
	address := fmt.Sprintf(":%d", acl.option.port)

	lis, err := net.Listen(acl.option.network, address)

	if err != nil {
		log.Printf(util.FatalString("failed to listen on port %d (AdminCommand): %v"), acl.option.port, err)
		acl.exitCh <- 0
		return
	}

	acl.listener = lis

	go acl.stopAdminCommandListener()

	log.Printf(util.SystemString("SYSTEM: Admin command gRPC listener opened { port: %d }"), acl.option.port)
	wg.Done()

	if err := acl.grpcServer.Serve(lis); err != nil {
		log.Printf(util.FatalString("failed to serve gRPC listener (AdminCommand) over port %d: %v"), acl.option.port, err)
		acl.exitCh <- 0
	}
}

func (acl *AdminCommandListener) stopAdminCommandListener() {
	val := <-acl.exitCh
	log.Printf("AdminCommandListener: Received exit signal %d, shutting down...", val)
	acl.Stop()
}

func (acl *AdminCommandListener) Stop() {
	if acl.grpcServer != nil {
		log.Println("AdminCommandListener: Shutting down gRPC server...")
		acl.grpcServer.GracefulStop()
	}
	if acl.listener != nil {
		acl.listener.Close()
	}
	log.Println("AdminCommandListener: gRPC server stopped.")
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		log.Printf("ERROR: Failed to dial to get outbound IP: %v", err)
		return "127.0.0.1"
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
