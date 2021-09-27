package main

import (
	controller "awesomeProject/server/api/controller"
	proto "awesomeProject/server/api/proto/generated"
	"awesomeProject/server/dataBase"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

// main start a gRPC server and waits for connection
func main() {

	go func() {
		dataBase.ConnectDataBase()
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := controller.Server{}
	grpcServer := grpc.NewServer()
	proto.RegisterAuthenticationServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
