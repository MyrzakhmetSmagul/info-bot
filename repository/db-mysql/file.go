package db_mysql

import (
	"fmt"
	"tg-bot/repository"
)

func (s storage) AddFileToMessage(file *repository.File) error {
	query := `INSERT INTO file(msg_group_id, file_name, file_type) VALUES (?, ?, ?)`

	exec, err := s.db.Exec(query, file.MsgGroupID, file.FileName, file.FileType)
	if err != nil {
		return fmt.Errorf("add file to message failed: %w", err)
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return fmt.Errorf("add file to message failed: %w", err)
	}

	file.ID = int(id)

	return nil
}

func (s storage) GetFileByID(fileID int) (repository.File, error) {
	query := `SELECT msg_group_id, file_name, file_type FROM file WHERE id=?`

	file := repository.File{ID: fileID}
	err := s.db.QueryRow(query, fileID).Scan(&file.MsgGroupID, &file.FileName, &file.FileType)
	if err != nil {
		return repository.File{}, fmt.Errorf("get file by id failed: %w", err)
	}

	return file, err
}

func (s storage) GetFilesOfMsgGroup(msgGroupID int) ([]repository.File, error) {
	query := `SELECT id, file_name, file_type FROM file WHERE msg_group_id=?`

	rows, err := s.db.Query(query, msgGroupID)
	if err != nil {
		return nil, fmt.Errorf("get files of msgGroup: %w", err)
	}

	files := make([]repository.File, 0)
	defer rows.Close()

	for rows.Next() {
		file := repository.File{MsgGroupID: msgGroupID}

		err = rows.Scan(&file.ID, &file.FileName, &file.FileType)
		if err != nil {
			return nil, fmt.Errorf("get files of msgGroup: %w", err)
		}

		files = append(files, file)
	}

	return files, nil
}

func (s storage) DeleteFile(fileID int) error {
	query := `DELETE FROM file WHERE id=?`

	//you should check how many rows affected
	//if affected rows is zero, return custom error which describe it
	_, err := s.db.Exec(query, fileID)
	if err != nil {
		return fmt.Errorf("delete file failed: %w", err)
	}

	return nil
}
