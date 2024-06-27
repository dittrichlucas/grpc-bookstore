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
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	ca "github.com/dittrichlucas/poc-grpc-bookstore/auth/client"
	pb "github.com/dittrichlucas/poc-grpc-bookstore/proto"
)

type BookClient struct {
	service pb.BookClient
}

// const (
// 	address  = "localhost:9000"
// 	username = "user"
// 	password = "secret"
// )

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Query the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		// serverAddress := flag.String("address", address, "the server address")
		// flag.Parse()
		// log.Printf("dial server %s", *serverAddress)

		// cc1, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
		// if err != nil {
		// 	log.Fatal("cannot dial server: ", err)
		// }

		// refreshDuration := 15 * time.Second
		// token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzI1NDA4NjksInVzZXJuYW1lIjoidXNlcjEiLCJyb2xlIjoidXNlciJ9.MpRuAXpaAUDy5MYgMvjzVUFia7AwLPPZZEgnYgv5gWE"
		// authClient := ci.NewAuthClient(cc1, token)                                                  // Vai autenticar o serviço com o usuário e a senha (preciso entender bem como isso funciona)
		// interceptor, err := ci.NewAuthClientInterceptor(authClient, authMethods(), refreshDuration) // Esse cara é o responsável por pegar o token e garantir que o usuário está autorizado a consumir as rotas
		// if err != nil {
		// 	log.Fatal("cannot create auth interceptor: ", err)
		// }

		// cc2, err := grpc.Dial(
		// 	*serverAddress,
		// 	grpc.WithInsecure(),
		// 	grpc.WithUnaryInterceptor(interceptor.Unary()),
		// 	grpc.WithStreamInterceptor(interceptor.Stream()),
		// )
		// if err != nil {
		// 	log.Fatal("cannot dial server: ", err)
		// }

		cc2 := ca.ClientAuthenticator()
		client := pb.NewBookClient(cc2)

		var id string
		if len(os.Args) > 2 {
			id = os.Args[2]
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := client.GetBook(ctx, &pb.Request{Id: id})
		if err != nil {
			log.Fatalf("could not gree: %v", err)
		}
		log.Printf("%s", r.GetMessage())
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}

func NewBookClient(cc *grpc.ClientConn) *BookClient {
	service := pb.NewBookClient(cc)
	return &BookClient{service}
}

// func authMethods() map[string]bool {
// 	const bookServicePath = "/books.Book/"
// 	return map[string]bool{
// 		bookServicePath + "GetBook": true,
// 		// bookServicePath + "": {"admin"},
// 	}
// }
