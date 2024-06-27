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

	sr "github.com/dittrichlucas/poc-grpc-bookstore/auth/server"
	pb "github.com/dittrichlucas/poc-grpc-bookstore/proto"
	s "github.com/dittrichlucas/poc-grpc-bookstore/service"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		userStore := s.NewInMemoryUserStore()
		if err := seedUsers(userStore); err != nil {
			log.Fatal("cannot seed users: ", err)
		}

		jwtManager := s.NewJWTManager(secretKey, tokenDuration)
		authService := sr.NewAuthService(userStore, *jwtManager)

		// TODO: receber user e password dinamicamente por meio de flag
		req := pb.LoginRequest{
			Username: "user1",
			Password: "secret",
		}
		lr, err := authService.Login(context.Background(), &req)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(lr.AccessToken)
		}

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
