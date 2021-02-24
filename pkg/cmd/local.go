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
	"log"

	. "github.com/MichaelWasher/GoXO/pkg/game"
	"github.com/MichaelWasher/GoXO/pkg/io"
	"github.com/spf13/cobra"
)

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Start a local game of GoXO and alternate between player turns",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting the GoXO Local Game")

		// Configure Inputs

		terminal := io.NewTerminal()
		defer terminal.Close() // Defer is LIFO ordering, Close is last.
		// defer terminal.Restore()
		log.Println("Using Terminal as input")

		// Create the Game
		log.Println("Game Created")
		game := NewGame(terminal, terminal)
		defer game.CloseGame()

		// Configure Event Listeners
		log.Println("GameLoop Started")
		game.GameLoop()
		log.Println("GameLoop Finished")

	},
}

func init() {
	rootCmd.AddCommand(localCmd)
}
