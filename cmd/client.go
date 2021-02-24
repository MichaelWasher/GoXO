/*
Copyright © 2021 Michael Washer <michael.washer@icloud.com>

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
	"fmt"
	gameproject "github.com/MichaelWasher/GoXO/game"
	"github.com/MichaelWasher/GoXO/grpc"
	"github.com/MichaelWasher/GoXO/input"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Connect to a game of GoXO",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("client called")
		game := gameproject.Game{}


		game.InitGame()
		defer game.CloseGame()

		go input.HandleKeyEvents(game.Terminal, game.GetPlayerTwoInputChannel())

		// TODO Set Args
		grpc.SetupClient(7777, game.GetPlayerTwoInputChannel())
		
		//game.GameLoop()
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	// TODO add client Details
}
