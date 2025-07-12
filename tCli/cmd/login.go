package cmd

import (
	"bufio"
	"context"
	"fmt"
	client "halo/client"
	config "halo/config"
	"halo/logger"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

var LoginCommand = &cli.Command{
	Name:  "login",
	Usage: "Login with login and password",
	Action: func(ctx context.Context, c *cli.Command) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter login: ")
		login, _ := reader.ReadString('\n')

		fmt.Print("Enter password: ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		logger.Logger.Err(err).Msg("Read password")

		login = strings.TrimSpace(login)
		password := strings.TrimSpace(string(passwordBytes))

		token, err := client.Login(login, password)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		err = config.SaveToken(token)
		if err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}

		fmt.Println("Logged in successfully.")
		return nil
	},
}
