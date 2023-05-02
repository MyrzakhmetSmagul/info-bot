package db_mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"tg-bot/repository"
)

func (s *storage) CreateMessage(message *repository.Message) error {
	query := `INSERT INTO message(message_trigger, text, lang, state_id)
				VALUES (?, ?, ?, ?)`

	_, err := s.db.Exec(query, message.MessageTrigger, message.Text, message.Language, message.StateID)
	if err != nil {
		return fmt.Errorf("can't create message: %w", err)
	}
	return nil
}

func (s *storage) UpdateMessage(message *repository.Message) error {
	query := `UPDATE message SET message_trigger=?, text=?, lang=?, state=? WHERE id=?`

	_, err := s.db.Exec(query, message.MessageTrigger, message.Text, message.Language, message.StateID, message.ID)
	if err != nil {
		return fmt.Errorf("can't update message: %w", err)
	}
	return nil
}

func (s *storage) DeleteMessage(messageID int) error {
	query := `DELETE FROM message WHERE id=?`

	if _, err := s.db.Exec(query, messageID); err != nil {
		return fmt.Errorf("can't delete message: %w", err)
	}

	return nil
}

func (s *storage) GetMessageByID(messageID int) (*repository.Message, error) {
	query := `SELECT message_trigger, text, lang, state_id FROM message WHERE id=?`
	message := &repository.Message{ID: messageID}
	err := s.db.QueryRow(query, messageID).Scan(&message.MessageTrigger, &message.Text, &message.Language, &message.StateID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get message by id: %w", err)
	}
	return message, nil
}

func (s *storage) GetMessageByTriggerAndState(messageTrigger string, stateID int) (*repository.Message, error) {
	query := `SELECT id, text, lang FROM message WHERE message_trigger=? AND state_id=?`
	message := &repository.Message{MessageTrigger: messageTrigger, StateID: stateID}
	err := s.db.QueryRow(query, messageTrigger, stateID).Scan(&message.ID, &message.Text, &message.Language)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get message by trigger: %w", err)
	}
	return message, nil
}

func (s *storage) GetAllMessages() ([]*repository.Message, error) {
	query := `SELECT id, message_trigger, text, lang, state_id FROM message`

	messages := make([]*repository.Message, 0)
	rows, err := s.db.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		message := &repository.Message{}
		err = rows.Scan(&message.ID, &message.MessageTrigger, &message.Text, &message.Language, &message.StateID)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
