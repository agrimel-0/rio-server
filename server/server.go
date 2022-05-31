package server

import (
	"fmt"
	"log"
	"net"

	pb "github.com/agrimel-0/rio-grpc"
	"google.golang.org/grpc"
)

// Pins config
type Pins struct {
	GpioChip   string `mapstructure:"gpiochip"`
	LineOffset int    `mapstructure:"lineOffset"`
	Alias      string `mapstructure:"alias"`
	Value      int    `mapstructure:"value"`
	Output     bool   `mapstructure:"output"`
}

// Server config
type Server struct {
	Port  int    `mapstructure:"port"`
	Alias string `mapstructure:"alias"`
}

// Service config
type Config struct {
	Server  Server            `mapstructure:"server"`
	PinList []map[string]Pins `mapstructure:"pins"`
}

// Server struct
type server struct {
	pb.UnimplementedRioServer

	serverAlias string // Optional alias for the server
	serverPort  int    // Server port

	exportedPins []*IoPin // Slice containing the exported pins

	grpcInstance *grpc.Server
}

// Start the server
func Start(serverconfig Config) error {

	ioPins, errs := IoFromConfig(serverconfig.PinList)
	for _, err := range errs {
		log.Printf("io setup error: %v", err)
	}

	server := server{
		exportedPins: ioPins,
		serverAlias:  serverconfig.Server.Alias,
		serverPort:   serverconfig.Server.Port,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", server.serverPort))
	if err != nil {
		return err
	}
	s := grpc.NewServer()

	server.grpcInstance = s

	pb.RegisterRioServer(server.grpcInstance, &server)
	log.Printf("network %v\n", lis.Addr().Network())
	log.Printf("%s listening at %v", server.serverAlias, lis.Addr())

	if err := server.grpcInstance.Serve(lis); err != nil {
		return err
	}

	return nil
}
