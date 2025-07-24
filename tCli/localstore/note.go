package localstore

import (
	"halo/models"
)

func AddNoteLocally(note models.NoteStruct) error {
	query := `INSERT INTO notes (id, type_id, content, created_at, updated_at, ended_at, completed) 
						VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err := db.Exec(
		query,
		note.Id,
		note.Type_id,
		note.Content,
		note.Created_at,
		note.Updated_at,
		note.Ended_at,
		note.Completed,
	)
	if err != nil {
		return err
	}
	return nil
}
