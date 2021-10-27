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

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
		user_list: &pb.UserList{},
	}
}

//implementation of the grpc service
//to register this type with grpc we habe to embed the UnimplementedUserManagementServer inside of the type
type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
	user_list *pb.UserList
}

func (server *UserManagementServer) Run() error {
	//net.Listen function to begin listening on the port specified above
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//variable 's' is invoking the new server function from the grpc module
	s := grpc.NewServer()

	//after initialize this new server we're going to register the server
	//as a new grpc service
	pb.RegisterUserManagementServer(s, server)
	//printa no terminal do server o endereço
	log.Printf("server listening at %v", lis.Addr())
	//invoke the server
	return s.Serve(lis)
}

//begin defining our service method that we defined in the proto file
func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	//printa no terminal do server os nomes dos usuários
	log.Printf("Received: %v", in.GetName())
	var user_id int32 = int32(rand.Intn(1000))
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}
	s.user_list.Users = append(s.user_list.Users, created_user)
	return created_user, nil
}

func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	return s.user_list, nil
}

func main() {

	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to server: %v", err)
	}

}
