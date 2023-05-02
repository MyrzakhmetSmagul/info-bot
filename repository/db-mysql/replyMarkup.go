package db_mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"tg-bot/repository"
)

func (s *storage) CreateReplyMarkup(replyMarkup *repository.ReplyMarkup) error {
	query := `INSERT INTO reply_markup(message_id, keyboard_id)
				VALUES (?, ?)`

	exec, err := s.db.Exec(query, replyMarkup.MessageID, replyMarkup.KeyboardID)
	if err != nil {
		return fmt.Errorf("can't create reply markup: %w", err)
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return fmt.Errorf("can't create reply markup: %w", err)
	}

	replyMarkup.ID = int(id)
	return nil
}

func (s *storage) UpdateReplyMarkup(replyMarkup *repository.ReplyMarkup) error {
	//TODO implement me
	panic("implement me")
}

func (s *storage) DeleteReplyMarkup(replyMarkupID int) error {
	query := `DELETE FROM reply_markup WHERE id = ?`
	_, err := s.db.Exec(query, replyMarkupID)
	if err != nil {
		return fmt.Errorf("can't delete reply markup: %w", err)
	}

	return nil
}

func (s *storage) GetReplyMarkupByID(replyMarkupID int) (*repository.ReplyMarkup, error) {
	query := `SELECT message_id, keyboard_id FROM reply_markup WHERE id=?`

	replyMarkup := &repository.ReplyMarkup{ID: replyMarkupID}
	err := s.db.QueryRow(query, replyMarkupID).Scan(&replyMarkup.MessageID, &replyMarkup.KeyboardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("can't get reply markup by id: %w", err)
	}

	return nil, nil
}

func (s *storage) GetAllReplyMarkups() ([]*repository.ReplyMarkup, error) {
	query := `SELECT id, message_id, keyboard_id FROM reply_markup`

	replyMarkups := make([]*repository.ReplyMarkup, 0)
	rows, err := s.db.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get all reply markups: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		replyMarkup := &repository.ReplyMarkup{}
		err = rows.Scan(&replyMarkup.ID, &replyMarkup.MessageID, &replyMarkup.KeyboardID)
		if err != nil {
			return nil, fmt.Errorf("can't get all reply markups: %w", err)
		}
		replyMarkups = append(replyMarkups, replyMarkup)
	}

	return replyMarkups, nil
}
