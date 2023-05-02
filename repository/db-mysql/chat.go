package db_mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"tg-bot/lib/e"
	"tg-bot/repository"
)

func (s *storage) CreateChat(info *repository.Chat) error {
	query := `INSERT INTO chat(chat_id, active, lang, state_id) VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(query, info.ChatID, info.Active, info.Language, info.StateID)
	if err != nil {
		return fmt.Errorf("can't save chat info %w", err)
	}

	return nil
}

func (s *storage) EnableChat(chatID int64) error {
	query := `UPDATE chat SET active=true WHERE chat_id=?`

	if _, err := s.db.Exec(query, chatID); err != nil {
		return fmt.Errorf("can't enable chat: %w", err)
	}

	return nil
}

func (s *storage) DeleteChat(chatID int64) error {
	query := `DELETE FROM chat WHERE chat_id=?`

	if _, err := s.db.Exec(query, chatID); err != nil {
		return fmt.Errorf("can't delete chat %w", err)
	}

	return nil
}

func (s *storage) DisableChat(chatID int64) error {
	query := `UPDATE chat SET active=false WHERE chat_id=?`

	if _, err := s.db.Exec(query, chatID); err != nil {
		return e.Wrap("can't disable chat", err)
	}

	return nil
}

func (s *storage) ChangeLang(info repository.Chat) error {
	query := `UPDATE chat SET lang=? WHERE chat_id=?`
	if _, err := s.db.Exec(query, info.Language, info.ChatID); err != nil {
		return fmt.Errorf("change state error: %w", err)
	}
	return nil
}

func (s *storage) ChangeState(info repository.Chat) error {
	query := `UPDATE chat SET state_id=? WHERE chat_id=?`
	if _, err := s.db.Exec(query, info.StateID, info.ChatID); err != nil {
		return fmt.Errorf("change state error: %w", err)
	}

	return nil
}

func (s *storage) GetChatByID(chatID int64) (*repository.Chat, error) {
	query := `SELECT active, lang, state_id FROM chat WHERE chat_id=?`
	result := &repository.Chat{ChatID: chatID}

	if err := s.db.QueryRow(query, chatID).Scan(&result.Active, &result.Language, &result.StateID); err != nil {
		return nil, e.Wrap("can't get chat info from db", err)
	}

	return result, nil
}

func (s *storage) GetAllChats() ([]*repository.Chat, error) {
	query := `SELECT chat_id, active, lang, state_id FROM chat`

	chats := make([]*repository.Chat, 0)
	rows, err := s.db.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		chat := &repository.Chat{}
		err = rows.Scan(chat.ChatID, chat.Active, chat.Language, chat.StateID)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}
