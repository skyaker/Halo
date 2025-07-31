package localstore

import (
	"halo/logger"
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

func GetNotesLocally(currentPage int, pageSize int) []models.NoteStruct {
	if currentPage < 0 {
		currentPage = 0
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := currentPage * pageSize

	query := `SELECT id, type_id, content, created_at, updated_at, ended_at, completed
						FROM notes
						ORDER BY created_at DESC
						LIMIT $1 OFFSET $2;`
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("note receiving")
		return nil
	}

	defer rows.Close()

	notes := make([]models.NoteStruct, 0)

	for rows.Next() {
		var noteInfo models.NoteStruct
		var createdAt, updatedAt, endedAt int

		err := rows.Scan(
			&noteInfo.Id,
			&noteInfo.Type_id,
			&noteInfo.Content,
			&createdAt,
			&updatedAt,
			&endedAt,
			&noteInfo.Completed,
		)
		if err != nil {
			logger.Logger.Error().Err(err).Msg("note info scan")
			return nil
		}
		notes = append(notes, noteInfo)
	}

	return notes
}

func DeleteNoteLocally(id string) error {
	query := `DELETE FROM notes WHERE id = $1`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func GetNumberOfNotes() (int, error) {
	query := `SELECT COUNT(*) FROM notes`

	row := db.QueryRow(query)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
