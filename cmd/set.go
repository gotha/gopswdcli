/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gotha/gopswdcli/pkg/secrets"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set a secret",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("You should provide secret name as parameter")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		keyring, err := secrets.NewLinuxKeyring(keychainName)
		if err != nil {
			fmt.Printf("error opening keychain: %s\n", err)
			os.Exit(1)
		}

		usernameReader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter username:")
		username, err := usernameReader.ReadString('\n')
		if err != nil {
			fmt.Printf("error reading username: %s\n", err)
			os.Exit(1)
		}
		username = strings.TrimSuffix(username, "\n")

		passwordReader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter password:")
		password, err := passwordReader.ReadString('\n')
		if err != nil {
			fmt.Printf("error reading password: %s\n", err)
			os.Exit(1)
		}
		password = strings.TrimSuffix(password, "\n")

		fmt.Printf("%s:%s\n", username, password)

		err = keyring.Set(args[0], username, password)
		if err != nil {
			fmt.Println("error saving secret")
			os.Exit(1)
		}
		fmt.Printf("savied secret %s in %s keychain\n", args[0], keychainName)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
