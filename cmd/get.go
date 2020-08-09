package cmd

import (
	"fmt"
	"os"

	"github.com/gotha/gopswdcli/pkg/secrets"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "return secret",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("You should provide argument")
		}
		allowedCommands := []string{"username", "password", "secret"}
		for _, allowedCommand := range allowedCommands {
			if allowedCommand == args[0] {
				return nil
			}
		}
		return fmt.Errorf("Invalid argument %s", args[0])
	},
}

var usernameCmd = &cobra.Command{
	Use:   "username",
	Short: "print username from secret",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("You should provide secret name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		keyring, err := secrets.NewKeyring(keychainName)
		if err != nil {
			fmt.Printf("error opening keychain: %s", err)
			os.Exit(1)
		}

		username, _, err := keyring.Get(args[0])
		if err != nil {
			fmt.Printf("error getting secret: %s", err)
			os.Exit(1)
		}
		fmt.Print(username)
	},
}

var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "print password from secret",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("You should provide secret name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		keyring, err := secrets.NewKeyring(keychainName)
		if err != nil {
			fmt.Printf("error opening keychain: %s", err)
			os.Exit(1)
		}

		_, password, err := keyring.Get(args[0])
		if err != nil {
			fmt.Printf("error getting secret: %s", err)
			os.Exit(1)
		}
		fmt.Print(password)
	},
}

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "print the whole secret",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("You should provide secret name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		secretKey := args[0]

		keyring, err := secrets.NewKeyring(keychainName)
		if err != nil {
			fmt.Printf("error opening keychain: %s", err)
			os.Exit(1)
		}

		username, password, err := keyring.Get(secretKey)
		if err != nil {
			fmt.Printf("error getting secret: %s", err)
			os.Exit(1)
		}

		fmt.Printf("%s\t%s\t%s", username, password, secretKey)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(usernameCmd)
	getCmd.AddCommand(passwordCmd)
	getCmd.AddCommand(secretCmd)
}
