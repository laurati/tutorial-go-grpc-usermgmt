package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	pb "tutorial-go-grpc-usermgmt/usermgmt"

	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

//math/rand - is for generate random integers, used as the user Ids
//pb - is importing the protobuf module

const (
	port = ":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

//implementation of the grpc service
//to register this type with grpc we habe to embed the UnimplementedUserManagementServer inside of the type
type UserManagementServer struct {
	conn *pgx.Conn
	pb.UnimplementedUserManagementServer
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
func (server *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	//printa no terminal do server os nomes dos usuários
	log.Printf("Received: %v", in.GetName())

	createSql := `
	create table if not exixts users(
		id SERIAL PRIMARY KEY,
		name text,
		age int
	);
	`
	_, err := server.conn.Exec(context.Background(), createSql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Table creation failed: %v\n", err)
		os.Exit(1)
	}

	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge()}
	tx, err := server.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed: %v", err)
	}

	_, err = tx.Exec(context.Background(), "insert into users(name,age) values ($1, $2)", created_user.Name, created_user.Age)
	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}

	tx.Commit(context.Background())

	return created_user, nil
}

func (server *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {

	var users_list *pb.UserList = &pb.UserList{}

	rows, err := server.conn.Query(context.Background(), "select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := pb.User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		users_list.Users = append(users_list.Users, &user)
	}

	return users_list, nil
}

func main() {
	database_url := "postgres://postgres:mysecretpassword@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("unable to establish connection: %v", err)
	}

	defer conn.Close(context.Background())

	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	user_mgmt_server.conn = conn
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to server: %v", err)
	}

}
