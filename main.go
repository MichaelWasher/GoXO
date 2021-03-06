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
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MichaelWasher/GoXO/pkg/cmd"
)

// TODO Add Multiplayer Support
// TODO Add Socket Support for Multiple Player Input
// TODO Flags for the Socket Connection
// TODO Implement CLI Options

// Configure Logging
func initLog() *os.File {
	f, err := os.OpenFile("log-file.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}

func main() {
	file := initLog()
	defer file.Close()

	goxo := cmd.Root()
	if err := goxo.Execute(); err != nil {
		log.Printf("A fatal error has occurred and GoXO must close. %v", err)
		fmt.Printf("A fatal error has occurred and GoXO must close. %v", err)
		os.Exit(1)
	}
}
