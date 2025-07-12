package main

import (
	"context"
	cmd "halo/cmd"
	"halo/logger"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func loadEnv() {
	_ = godotenv.Load()
}

func isDebugEnabled() bool {
	return os.Getenv("HALO_DEBUG") == "1"
}

func main() {
	loadEnv()

	debug := isDebugEnabled()
	logger.Init(debug)

	cmd := &cli.Command{
		Name:  "halo",
		Usage: "Terminal client for habit tracking",
		Commands: []*cli.Command{
			cmd.LoginCommand,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Logger.Error().Err(err).Msg("Run command")
	}
}
