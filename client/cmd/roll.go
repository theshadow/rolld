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
	"context"
	"fmt"

	"time"

	"github.com/spf13/cobra"

	"google.golang.org/grpc"

	"github.com/theshadow/dice"
	"github.com/theshadow/dice/formula"

	"github.com/theshadow/roll/encoding"

	"github.com/theshadow/rolld/server"
)

var address string

// startCmd represents the start command
var rollCmd = &cobra.Command{
	Use:   "roll",
	Short: "Rolls some dice",
	Long:  `Takes a roll formula and rolls the corresponding dice returning the result.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("missing argument <formula>")
		}

		// create connection.
		conn, err := grpc.Dial(":50051", grpc.WithInsecure())
		if err != nil {
			return err
		}

		req := server.RollRequest{
			Formula: args[0],
		}

		context.WithTimeout(context.Background(), time.Second * 5)

		clnt := server.NewRollerClient(conn)
		res, err := clnt.Roll(context.Background(), &req)
		if err != nil {
			return err
		}

		rs, err := FromgRPCRoll(*res.Rolls)
		if err != nil {
			return err
		}

		enc := encoding.JSONEncoder{}
		s, err := enc.Encode(rs)
		if err != nil {
			return err
		}

		fmt.Println(s)

		return nil
	},
}

func FromgRPCRoll(roll server.Result) (dice.Results, error) {
	var rs []int
	for _, r := range roll.Roll {
		rs = append(rs, int(r))
	}

	return dice.Results{
		Formula: formula.Formula(roll.Formula),
		Rolls: rs,
		Extensions: roll.Extensions,
	}, nil
}

func init() {
	rootCmd.AddCommand(rollCmd)

	rollCmd.Flags().StringVarP(&address, "address", "a",
		":50051", "A host and port in the form of 'host:port' to listen on.")
}
