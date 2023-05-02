package db_mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"tg-bot/repository"
)

func (s *storage) AddFile(file *repository.File) error {
	query := `INSERT INTO file(message_id, file_name, file_type)
				VALUES (?, ?, ?)`

	_, err := s.db.Exec(query, file.MessageID, file.FileName, file.FileType)
	if err != nil {
		return fmt.Errorf("can't add file: %w", err)
	}

	return nil
}

func (s *storage) DeleteFile(fileID int) error {
	query := `DELETE FROM file WHERE id=?`
	_, err := s.db.Exec(query, fileID)
	if err != nil {
		return fmt.Errorf("can't delete file: %w", err)
	}

	return nil
}

func (s *storage) GetFileByID(fileID int) (*repository.File, error) {
	query := `SELECT message_id, file_name, file_type FROM file WHERE id=?`

	file := &repository.File{ID: fileID}
	err := s.db.QueryRow(query, fileID).Scan(&file.MessageID, &file.FileName, &file.FileType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get file by id: %w", err)
	}

	return file, nil
}

func (s *storage) GetFilesByMessageID(messageID int) ([]*repository.File, error) {
	query := `SELECT id, message_id, file_name, file_type FROM file WHERE message_id=?`

	files := make([]*repository.File, 0)
	rows, err := s.db.Query(query, messageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get files by message id: %w", err)
	}

	defer rows.Close()
	
	for rows.Next() {
		file := &repository.File{}
		err = rows.Scan(&file.ID, &file.MessageID, &file.FileName, &file.FileType)
		if err != nil {
			return nil, fmt.Errorf("can't get files by message id: %w", err)
		}
		files = append(files, file)
	}

	return files, nil
}

func (s *storage) GetAllFiles() ([]*repository.File, error) {
	query := `SELECT id, message_id, file_name, file_type FROM file`

	files := make([]*repository.File, 0)
	rows, err := s.db.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get files by message id: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		file := &repository.File{}
		err = rows.Scan(&file.ID, &file.MessageID, &file.FileName, &file.FileType)
		if err != nil {
			return nil, fmt.Errorf("can't get files by message id: %w", err)
		}
		files = append(files, file)
	}

	return files, nil
}
