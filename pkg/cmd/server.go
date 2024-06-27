/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"

	sr "github.com/dittrichlucas/poc-grpc-bookstore/auth/server"
	pb "github.com/dittrichlucas/poc-grpc-bookstore/proto"
	s "github.com/dittrichlucas/poc-grpc-bookstore/service"
)

const (
	port          = ":9000"
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

type Server struct {
	pb.UnimplementedBookServer
}

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

/*
func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Println("--> unary interceptor: ", info.FullMethod)
	return handler(ctx, req)
}

func streamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Println("--> stream interceptor: ", info.FullMethod)
	return handler(srv, stream)
}
*/

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the Schema gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		userStore := s.NewInMemoryUserStore()
		if err := seedUsers(userStore); err != nil {
			log.Fatal("cannot seed users: ", err)
		}

		jwtManager := s.NewJWTManager(secretKey, tokenDuration)
		authService := sr.NewAuthService(userStore, *jwtManager)

		interceptor := sr.NewAuthInterceptor(jwtManager, accessibleRoles())
		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(interceptor.Unary()),
			grpc.StreamInterceptor(interceptor.Stream()),
		)

		// Register services
		pb.RegisterBookServer(grpcServer, &Server{})
		pb.RegisterAuthServiceServer(grpcServer, authService)

		log.Printf("GRPC server listening on %v", lis.Addr())

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func (s *Server) GetBook(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	res := &pb.Response{}

	// Check request
	if req == nil {
		fmt.Println("request mus not be nil")
		return res, xerrors.Errorf("request must not be nil")
	}

	if req.Id == "" {
		fmt.Println("id must not be empty in the request")
		return res, xerrors.Errorf("id must not be empty in the request")
	}

	log.Printf("Received: %v", req.GetId())

	// Get book
	file := filepath.Join("books.json")

	// Mock data
	book := []Book{}
	book = append(book, Book{ID: req.Id, Title: "Whatever", Author: "Mock ok"})
	bookJSON, _ := json.Marshal(book)
	ioutil.WriteFile(file, bookJSON, os.ModePerm)

	data, _ := ioutil.ReadFile(file)

	var books []Book
	if err := json.Unmarshal(data, &books); err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	var filterBook Book
	for _, myBook := range books {
		if myBook.ID == req.Id {
			filterBook = Book{
				myBook.ID,
				myBook.Author,
				myBook.Title,
			}
		}
	}

	res.Message = fmt.Sprintf("Call made successfully!\tAuthor: %s", filterBook.Author)

	return res, nil
}

func (s *Server) ListBook(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	res := &pb.Response{}

	// // Check request
	// if req == nil {
	// 	fmt.Println("request mus not be nil")
	// 	return res, xerrors.Errorf("request must not be nil")
	// }

	// if req.Id == "" {
	// 	fmt.Println("id must not be empty in the request")
	// 	return res, xerrors.Errorf("id must not be empty in the request")
	// }

	// log.Printf("Received: %v", req.GetId())

	// // Get book
	// file := filepath.Join("books.json")

	// // Mock data
	// book := []Book{}
	// book = append(book, Book{ID: req.Id, Title: "Whatever", Author: "Mock ok"})
	// bookJSON, _ := json.Marshal(book)
	// ioutil.WriteFile(file, bookJSON, os.ModePerm)

	// data, _ := ioutil.ReadFile(file)

	// var books []Book
	// if err := json.Unmarshal(data, &books); err != nil {
	// 	log.Fatalf("failed to unmarshal JSON: %v", err)
	// }

	// var filterBook Book
	// for _, myBook := range books {
	// 	if myBook.ID == req.Id {
	// 		filterBook = Book{
	// 			myBook.ID,
	// 			myBook.Author,
	// 			myBook.Title,
	// 		}
	// 	}
	// }

	res.Message = fmt.Sprintf("Call made successfully!\tID received: %s", req.Id)

	return res, nil
}

func createUser(userStore s.UserStore, username, password, role string) error {
	user, err := s.NewUser(username, password, role)
	if err != nil {
		return err
	}

	return userStore.Save(user)
}

func seedUsers(userStore s.UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}

	return createUser(userStore, "user1", "secret", "user")
}

func accessibleRoles() map[string][]string {
	const bookServicePath = "/books.Book/"
	return map[string][]string{
		bookServicePath + "GetBook":  {"user"},
		bookServicePath + "ListBook": {"admin"},
		// bookServicePath + "": {"admin"},
	}
}
