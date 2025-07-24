package cmd

import (
	"context"
	"fmt"
	"halo/client"
	"halo/config"
	"halo/localstore"
	"halo/logger"
	"halo/models"

	"github.com/google/uuid"
	"github.com/urfave/cli/v3"
)

var AddNoteCommand = &cli.Command{
	UseShortOptionHandling: true,
	Name:                   "add",
	Usage:                  "Add note",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		var noteInfo models.NoteStruct

		if cmd.Args().Len() < 1 {
			logger.Logger.Error().Msg("no content")
			return fmt.Errorf("no content")
		}

		noteInfo.Content = cmd.Args().Get(0)

		if noteId, err := uuid.NewV7(); err != nil {
			logger.Logger.Error().Err(err).Msg("new note id generation")
			return fmt.Errorf("new note id generation: %w", err)
		} else {
			noteInfo.Id = noteId
		}

		token, err := config.LoadToken()
		if err != nil {
			logger.Logger.Error().Err(err).Msg("get session token")
			return fmt.Errorf("get session token: %w", err)
		}

		err = client.SendNoteToService(token, noteInfo)
		if err != nil {
			logger.Logger.Error().Err(err).Msg("send note to service")
			return fmt.Errorf("send note to service: %w", err)
		}

		if err := localstore.AddNoteLocally(noteInfo); err != nil {
			logger.Logger.Error().Err(err).Msg("save note")
			return fmt.Errorf("save note error: %w", err)
		}

		// go localstore.SyncNotes(note)

		fmt.Println("Note added successfully.")
		return nil
	},
}
