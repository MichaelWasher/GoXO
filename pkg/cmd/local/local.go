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
package local

import (
	"log"

	"github.com/MichaelWasher/GoXO/pkg/game"
	"github.com/MichaelWasher/GoXO/pkg/io"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "local",
	Short: "Start a local game of GoXO and alternate between player turns",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLocalGame()
	},
}

func runLocalGame() error {
	log.Println("Starting the GoXO Local Game")

	// Configure Inputs
	terminal, _ := io.NewTerminal()
	defer terminal.Close() // Defer is LIFO ordering, Close is last.
	defer terminal.Restore()

	log.Println("Using Terminal as input")

	// Create the Game
	log.Println("Game Created")
	gameObj := game.NewGame(terminal, terminal)
	defer gameObj.CloseGame()

	// Configure Event Listeners
	log.Println("GameLoop Started")
	gameObj.GameLoop()
	log.Println("GameLoop Finished")
	return nil
}
