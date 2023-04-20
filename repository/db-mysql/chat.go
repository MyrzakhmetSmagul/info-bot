package db_mysql

import (
	"fmt"
	"tg-bot/lib/e"
	"tg-bot/repository"
)

func (s *storage) ChatInfo(chatID int64) (repository.ChatInfo, error) {
	query := `SELECT active, lang, state FROM chat WHERE chat_id=?`
	result := repository.ChatInfo{ChatID: chatID}

	if err := s.db.QueryRow(query, chatID).Scan(&result.Active, &result.Lang, &result.State); err != nil {
		return repository.ChatInfo{}, e.Wrap("can't get chat info from db", err)
	}

	return result, nil
}

func (s *storage) SaveChat(info repository.ChatInfo) error {
	query := `INSERT INTO chat(chat_id, active, lang, state) VALUES (?, ?, ?, ?)`

	if _, err := s.db.Exec(query, info.ChatID, info.Active, info.Lang, info.State); err != nil {
		return e.Wrap("can't save chat info", err)
	}

	return nil
}

func (s *storage) EnableChat(chatID int64) error {
	query := `UPDATE chat SET active=true WHERE chat_id=?`

	if _, err := s.db.Exec(query, chatID); err != nil {
		return e.Wrap("can't enable chat", err)
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

func (s *storage) ChangeLang(info repository.ChatInfo) error {
	query := `UPDATE chat SET lang=? WHERE chat_id=?`
	if _, err := s.db.Exec(query, info.Lang, info.ChatID); err != nil {
		return fmt.Errorf("change state error: %w", err)
	}
	return nil
}

func (s *storage) ChangeState(info repository.ChatInfo) error {
	query := `UPDATE chat SET state=? WHERE chat_id=?`
	if _, err := s.db.Exec(query, info.State, info.ChatID); err != nil {
		return fmt.Errorf("change state error: %w", err)
	}

	return nil
}
