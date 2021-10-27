package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "tutorial-go-grpc-usermgmt/usermgmt"

	"google.golang.org/grpc"
)

//math/rand - is for generate random integers, used as the user Ids
//pb - is importing the protobuf module

const (
	port = ":50051"
)

//implementation of the grpc service
//to register this type with grpc we habe to embed the UnimplementedUserManagementServer inside of the type
type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

//begin defining our service method that we defined in the proto file
func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	//printa no terminal do server os nomes dos usuários
	log.Printf("Received: %v", in.GetName())
	var user_id int32 = int32(rand.Intn(1000))
	return &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}, nil
}

func main() {

	//net.Listen function to begin listening on the port specified above
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//variable 's' is invoking the new server function from the grpc module
	s := grpc.NewServer()

	//after initialize this new server we're going to register the server
	//as a new grpc service
	pb.RegisterUserManagementServer(s, &UserManagementServer{})
	//printa no terminal do server o endereço
	log.Printf("server listening at %v", lis.Addr())
	//invoke the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
