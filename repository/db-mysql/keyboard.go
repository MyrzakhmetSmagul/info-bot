package db_mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"tg-bot/repository"
)

func (s *storage) CreateKeyboard(keyboard *repository.Keyboard) error {
	query := `INSERT INTO keyboard(kz_text, ru_text, en_text)
			VALUES (?, ?, ?)`

	_, err := s.db.Exec(query, keyboard.KzText, keyboard.RuText, keyboard.EnText)
	if err != nil {
		return fmt.Errorf("can't create keyboard: %w", err)
	}
	return nil
}

func (s *storage) UpdateKeyboard(keyboard *repository.Keyboard) error {
	query := `UPDATE keyboard SET kz_text=?, ru_text=?, en_text=? WHERE id=?`

	_, err := s.db.Exec(query, keyboard.KzText, keyboard.RuText, keyboard.EnText)
	if err != nil {
		return fmt.Errorf("can't update keyboard: %w", err)
	}

	return nil
}

func (s *storage) DeleteKeyboard(keyboardID int) error {
	query := `DELETE FROM keyboard WHERE id=?`
	_, err := s.db.Exec(query, keyboardID)
	if err != nil {
		return fmt.Errorf("can't delete keyboard: %w", err)
	}

	return nil
}

func (s *storage) GetKeyboardByID(keyboardID int) (*repository.Keyboard, error) {
	query := `SELECT kz_text, ru_text, en_text FROM keyboard WHERE id=?`

	keyboard := &repository.Keyboard{ID: keyboardID}
	err := s.db.QueryRow(query, keyboardID).Scan(&keyboard.KzText, &keyboard.RuText, &keyboard.EnText)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get keyboard by id: %w", err)
	}

	return keyboard, nil
}

func (s *storage) GetAllKeyboards() ([]*repository.Keyboard, error) {
	query := `SELECT id, kz_text, ru_text, en_text FROM keyboard`

	keyboards := make([]*repository.Keyboard, 0)
	rows, err := s.db.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get all keyboards: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		keyboard := &repository.Keyboard{}
		err = rows.Scan(&keyboard.ID, &keyboard.KzText, &keyboard.RuText, &keyboard.EnText)
		if err != nil {
			return nil, fmt.Errorf("can't get all keyboards: %w", err)
		}
		keyboards = append(keyboards, keyboard)
	}

	return keyboards, nil
}
