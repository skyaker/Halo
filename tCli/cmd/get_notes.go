package cmd

import (
	"context"
	"fmt"
	"halo/localstore"

	"github.com/urfave/cli/v3"
)

var GetNotesCommand = &cli.Command{
	UseShortOptionHandling: true,
	Name:                   "get",
	Usage:                  "Get notes",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		notes := localstore.GetNotesLocally(0, 0)
		for i, note := range notes {
			fmt.Printf("%v. %v\n", i+1, note.Content)
		}
		return nil
	},
}
