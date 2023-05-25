package db_mysql

import (
	"fmt"
	"tg-bot/repository"
)

func (s storage) CreateReplyMarkup(rm *repository.ReplyMarkup) error {
	query := `INSERT INTO reply_markup(msg_group_id, reply_msg_group_id) VALUES (?, ?)`

	exec, err := s.db.Exec(query, rm.MessageGroupID, rm.ReplyMessageGroupID)
	if err != nil {
		return fmt.Errorf("create replyMarkup failed: %w", err)
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return fmt.Errorf("create replyMarkup failed: %w", err)
	}

	rm.ID = int(id)

	return nil
}

func (s storage) GetReplyMarkupByID(id int) (repository.ReplyMarkup, error) {
	query := `SELECT msg_group_id, reply_msg_group_id FROM reply_markup WHERE id=?`

	rm := repository.ReplyMarkup{ID: id}
	err := s.db.QueryRow(query, id).Scan(&rm.MessageGroupID, &rm.ReplyMessageGroupID)
	if err != nil {
		return repository.ReplyMarkup{}, fmt.Errorf("get replyMarkup by id failed: %w", err)
	}

	return rm, nil
}

func (s *storage) GetReplyMarkupsOfState(messageGroupID int) ([]repository.ReplyMarkup, error) {
	query := `SELECT id, reply_msg_group_id FROM reply_markup WHERE msg_group_id=?`

	rows, err := s.db.Query(query, messageGroupID)
	if err != nil {
		return nil, fmt.Errorf("get replyMarkups of state failed: %w", err)
	}

	rms := make([]repository.ReplyMarkup, 0)
	defer rows.Close()

	for rows.Next() {
		rm := repository.ReplyMarkup{MessageGroupID: messageGroupID}

		err = rows.Scan(&rm.ID, &rm.ReplyMessageGroupID)
		if err != nil {
			return nil, fmt.Errorf("get replyMarkups of state failed: %w", err)
		}
		rms = append(rms, rm)
	}

	return rms, nil
}

func (s storage) DeleteReplyMarkup(replyMarkupID int) error {
	query := `DELETE FROM reply_markup WHERE id=?`

	//you should check how many rows affected
	//if affected rows is zero, return custom error which describe it
	_, err := s.db.Exec(query, replyMarkupID)
	if err != nil {
		return fmt.Errorf("delete replyMarkup failed: %w", err)
	}

	return nil
}
