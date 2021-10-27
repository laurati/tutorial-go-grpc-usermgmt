package main

import (
	"context"
	"log"
	"time"

	pb "tutorial-go-grpc-usermgmt/usermgmt"

	"google.golang.org/grpc"
)

//address of the grpc server
const (
	address = "localhost:50051"
)

func main() {

	//dial a connection to grpc server
	//withBlock() means that this function will not return
	//until the connection is made
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//create a new client (to pass the connection to that function)
	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var new_users = make(map[string]int32)
	new_users["Alice"] = 43
	new_users["Bob"] = 30

	//call the create new user function by looping over the new users map
	for name, age := range new_users {
		//r - response from the grpc server
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatalf("could not create user: %v", err)
		}
		//printa no terminal do client os detalhes dos usuários criados
		log.Printf(`User Details:
	NAME: %s
	AGE: %d
	ID: %d`, r.GetName(), r.GetAge(), r.GetId())
	}
}
