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
	"github.com/MichaelWasher/GoXO/grpc"
	"github.com/MichaelWasher/GoXO/input"
	"github.com/spf13/cobra"
	gameproject "github.com/MichaelWasher/GoXO/game"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a game of GoXO and host a server for others to join",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
		game := gameproject.Game{}
		game.InitGame()
		defer game.CloseGame()

		go input.HandleKeyEvents(game.Terminal, game.GetPlayerOneInputChannel())

		// TODO Set Args
		go grpc.SetupServer(7777, game.GetPlayerTwoInputChannel())
		game.GameLoop()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)


}
