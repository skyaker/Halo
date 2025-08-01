package cmd

import (
	"context"
	"halo/logger"
	"halo/ui"

	"github.com/urfave/cli/v3"
)

var NoteListCommand = &cli.Command{
	UseShortOptionHandling: true,
	Name:                   "list",
	Usage:                  "note list",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		err := ui.StartNoteSelector()
		if err != nil {
			logger.Logger.Fatal().Err(err).Msg("ui error")
			return err
		}
		return nil
	},
}
