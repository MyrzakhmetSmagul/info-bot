package db_mysql

import (
	"fmt"
	"tg-bot/repository"
)

func (s storage) CreateChat(chat repository.Chat) error {
	query := `INSERT INTO chat(chat_id, active, lang, state_id, cmd) VALUES (?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, chat.ChatID, chat.Active, chat.Lang, chat.StateID, chat.CMD)
	if err != nil {
		return fmt.Errorf("create chat was failed: %w", err)
	}

	return nil
}

func (s storage) EnableChat(chatID int64) error {
	query := `UPDATE chat SET active=true WHERE chat_id=?`

	//you should check how many rows affected
	//if affected rows is zero, return custom error which describe it
	_, err := s.db.Exec(query, chatID)
	if err != nil {
		return fmt.Errorf("enable chat was failed: %w", err)
	}

	return nil
}

func (s storage) DisableChat(chatID int64) error {
	query := `UPDATE chat SET active=false WHERE chat_id=?`

	//you should check how many rows affected
	//if affected rows is zero, return custom error which describe it
	_, err := s.db.Exec(query, chatID)
	if err != nil {
		return fmt.Errorf("disable chat was failed: %w", err)
	}

	return nil
}

func (s storage) EnableCmd(chatID int64) error {
	query := `UPDATE chat SET cmd=true WHERE chat_id=?`

	//you should check how many rows affected
	//if affected rows is zero, return custom error which describe it
	_, err := s.db.Exec(query, chatID)
	if err != nil {
		return fmt.Errorf("enable cmd was failed: %w", err)
	}

	return nil
}

func (s storage) DisableCmd(chatID int64) error {
	query := `UPDATE chat SET cmd=false WHERE chat_id=?`

	_, err := s.db.Exec(query, chatID)
	if err != nil {
		return fmt.Errorf("disable cmd was failed: %w", err)
	}

	return nil
}

func (s storage) ChangeChatLang(chatID int64, lang string) error {
	query := `UPDATE chat SET lang=? WHERE chat_id=?`

	//you should check how many rows affected
	//if affected rows is zero, return custom error which describe it
	_, err := s.db.Exec(query, lang, chatID)
	if err != nil {
		return fmt.Errorf("change chat lang was failed: %w", err)
	}

	return nil
}

func (s storage) ChangeChatState(chatID int64, stateID int) error {
	query := `UPDATE chat SET state_id=? WHERE chat_id=?`

	//you should check how many rows affected
	//if affected rows is zero, return custom error which describe it
	_, err := s.db.Exec(query, stateID, chatID)
	if err != nil {
		return fmt.Errorf("change state was failed: %w", err)
	}

	return nil
}

func (s storage) GetChat(chatID int64) (repository.Chat, error) {
	query := `SELECT active, lang, state_id, cmd FROM chat WHERE chat_id=?`
	chat := repository.Chat{ChatID: chatID}

	err := s.db.QueryRow(query, chatID).Scan(&chat.Active, &chat.Lang, &chat.StateID, &chat.CMD)
	if err != nil {
		return repository.Chat{}, fmt.Errorf("get chat was failed: %w", err)
	}

	return chat, nil
}

func (s storage) GetAllChats() ([]repository.Chat, error) {
	query := `SELECT chat_id, active, lang, state_id, cmd FROM chat`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("get all chats was failed: %w", err)
	}

	chats := make([]repository.Chat, 0)
	defer rows.Close()

	for rows.Next() {
		chat := repository.Chat{}

		err = rows.Scan(&chat.ChatID, &chat.Active, &chat.Lang, &chat.StateID, &chat.CMD)
		if err != nil {
			return nil, fmt.Errorf("get all chats was failed: %w", err)
		}

		chats = append(chats, chat)
	}

	return chats, nil
}

func (s storage) DeleteChat(chatID int64) error {
	query := `DELETE FROM chat WHERE chat_id=?`

	//you should check how many rows affected
	//if affected rows is zero, return custom error which describe it
	_, err := s.db.Exec(query, chatID)
	if err != nil {
		return fmt.Errorf("delete chat was failed: %w", err)
	}

	return nil
}
