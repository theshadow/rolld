// Copyright © 2018 Xander Guzman <xander.guzman@xanderguzman.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/theshadow/rolld/server"
)

var address string
var seed string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the daemon",
	Long:  `Starts the daemon as a long running process`,
	RunE: func(cmd *cobra.Command, args []string) error {
		lis, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			return err
		}

		s := time.Now().UTC().UnixNano()
		if len(seed) > 0 {
			s, err = strconv.ParseInt(seed, 10, 64)
			if err != nil {
				log.Fatalf("failed to convert seed %s: %s", seed, err)
				return err
			}
		}

		rand.Seed(s)

		done := make(chan struct{})

		srv := grpc.NewServer()
		server.RegisterRollerServer(srv, server.New(srv, done))
		reflection.Register(srv)

		go func() {
			err := srv.Serve(lis)
			if err != nil {
				log.Fatalf("failed to start rpc server: %s", err)
			}
		}()

		<-done

		srv.GracefulStop()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&address, "address", "a",
		":50051", "A host and port in the form of 'host:port' to listen on.")
	startCmd.Flags().StringVarP(&seed, "seed", "s",
		"", "Set the seed for the RNG, must be a valid in64 value.")
}
