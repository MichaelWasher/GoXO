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
	"github.com/MichaelWasher/GoXO/pkg/cmd/local"
	"github.com/spf13/cobra"
)

var cfgFile string

func Root() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var cmd = &cobra.Command{
		Use:          "goxo",
		Short:        "A simple noughts and crosses game written in Go",
		Long:         ``,
		SilenceUsage: true,
	}
	// TODO Perform Additional Template functions

	// Add the Sub Commands
	cmd.AddCommand(
		local.Command,
	)
	// Execute adds all child commands to the root command and sets flags appropriately.

	return cmd
}
