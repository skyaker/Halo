package main

import (
	"context"
	"halo/cmd"
	"halo/localstore"
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

	localstore.GetLocalDbConnection()

	cmd := &cli.Command{
		Name:                   "halo",
		Usage:                  "Terminal client for habit tracking",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			cmd.LoginCommand,
			cmd.AddNoteCommand,
			cmd.NoteListCommand,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Logger.Error().Err(err).Msg("Run command")
	}
}
