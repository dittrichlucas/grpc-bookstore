package client

import (
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
)

const address = "localhost:9000"

func ClientAuthenticator() grpc.ClientConnInterface {
	serverAddress := flag.String("address", address, "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	cc1, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	refreshDuration := 15 * time.Second
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzI2MDUwMzIsInVzZXJuYW1lIjoidXNlcjEiLCJyb2xlIjoidXNlciJ9.puHOdITQh9MoSfmrb9qwtE_LT-BwslcJ94aA5VXVXFY"
	authClient := NewAuthClient(cc1, token)                                                  // Vai autenticar o serviço com o usuário e a senha (preciso entender bem como isso funciona)
	interceptor, err := NewAuthClientInterceptor(authClient, authMethods(), refreshDuration) // Esse cara é o responsável por pegar o token e garantir que o usuário está autorizado a consumir as rotas
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	cc2, err := grpc.Dial(
		*serverAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	return cc2
}

func authMethods() map[string]bool {
	const bookServicePath = "/books.Book/"
	return map[string]bool{
		bookServicePath + "GetBook":  true,
		bookServicePath + "ListBook": true,
		// bookServicePath + "": {"admin"},
	}
}
