package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
	"golang.org/x/term"

	client "halo/client"
	config "halo/config"
)

var RegisterCommand = &cli.Command{
	Name:  "register",
	Usage: "Register with login and password",
	Action: func(ctx context.Context, c *cli.Command) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter login: ")
		login, _ := reader.ReadString('\n')

		var passwordBytes []byte
		var err error

		for {
			fmt.Print("Enter password: ")
			passwordBytes, err = term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}

			fmt.Print("Confirm password: ")
			confirmPasswordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}

			if passwordBytes == nil {
				fmt.Println("Password is empty")
				continue
			}

			if bytes.Equal(passwordBytes, confirmPasswordBytes) {
				break
			}

			fmt.Println("Passwords don't match")
		}

		fmt.Print("Enter username (optional): ")
		username, _ := reader.ReadString('\n')

		fmt.Print("Enter email (optional): ")
		email, _ := reader.ReadString('\n')

		login = strings.TrimSpace(login)
		password := strings.TrimSpace(string(passwordBytes))
		username = strings.TrimSpace(username)
		email = strings.TrimSpace(email)

		token, err := client.Register(login, password, username, email)
		if err != nil {
			return fmt.Errorf("register failed: %w", err)
		}

		err = config.SaveToken(token)
		if err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}

		fmt.Println("Registered successfully.")
		return nil
	},
}
