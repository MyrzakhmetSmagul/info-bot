package db_mysql

import (
	"fmt"
	"tg-bot/repository"
)

func (s storage) CreateMessage(msg *repository.Message) error {
	query := `INSERT INTO message(msg_trigger, text, lang) VALUES (?, ?, ?)`

	exec, err := s.db.Exec(query, msg.MsgTrigger, msg.Text, msg.Lang)
	if err != nil {
		return fmt.Errorf("create message was failed: %w", err)
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return fmt.Errorf("create message was failed: %w", err)
	}

	msg.ID = int(id)
	return nil
}

func (s storage) UpdateMessage(msg repository.Message) error {
	query := `UPDATE message SET msg_trigger=?, text=? WHERE id=?`

	//you should check how many rows affected
	//if affected rows is zero, return custom error which describe it
	_, err := s.db.Exec(query, msg.MsgTrigger, msg.Text, msg.ID)
	if err != nil {
		return fmt.Errorf("update message was failed: %w", err)
	}

	return nil
}

func (s storage) GetMessage(trigger, lang string) (repository.Message, error) {
	query := `SELECT id, text FROM message WHERE msg_trigger=? AND lang=?`

	msg := repository.Message{MsgTrigger: trigger, Lang: lang}
	err := s.db.QueryRow(query, trigger, lang).Scan(&msg.ID, &msg.Text)
	if err != nil {
		return repository.Message{}, fmt.Errorf("get message was failed: %w", err)
	}

	return msg, nil
}

func (s storage) GetAllMessages() ([]repository.Message, error) {
	query := `SELECT id, msg_trigger, text, lang FROM message`

	msgs := make([]repository.Message, 0)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("get all messages was failed: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		msg := repository.Message{}
		err = rows.Scan(&msg.ID, &msg.MsgTrigger, &msg.Text, &msg.Lang)
		if err != nil {
			return nil, fmt.Errorf("get all messages was failed: %w", err)
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (s storage) DeleteMessage(msgID int) error {
	query := `DELETE FROM message WHERE id=?`

	//you should check how many rows affected
	//if affected rows is zero, return custom error which describe it
	_, err := s.db.Exec(query, msgID)
	if err != nil {
		return fmt.Errorf("delete message was failed: %w", err)
	}

	return nil
}
