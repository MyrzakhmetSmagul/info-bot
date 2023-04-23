package db_mysql

import (
	"database/sql"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-bot/repository"
)

func (s *storage) GetMessage(trigger string, lang repository.Lang, prevState repository.State) (*repository.Message, error) {
	msg := &repository.Message{Trigger: trigger, Lang: lang, PrevState: prevState}
	query := `SELECT id, text, state_id FROM message WHERE prev_state_id=? and message_trigger=? and lang=?`

	if err := s.db.QueryRow(query, prevState.ID, trigger, lang).Scan(&msg.ID, &msg.Text, &msg.State.ID); err != nil {
		return nil, fmt.Errorf("can't get lang message: %w", err)
	}

	return msg, nil
}

func (s *storage) CreateMessage(msg repository.Message) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// insert messages
	query := `INSERT INTO message(message_trigger, text, lang, state_id, prev_state_id) VALUES (?, ?, ?, ?, ?)`
	if msg.Text == "" || msg.Trigger == "" {
		return fmt.Errorf("error creating message: msg data is empty")
	}

	if _, err = s.db.Exec(query, msg.Trigger, msg.Text, msg.Lang, msg.State.ID, msg.PrevState.ID); err != nil {
		return fmt.Errorf("error creating message: %w", err)
	}

	return nil
}

func (s *storage) GetReplyMarkup(msgID int64, lang repository.Lang) (*tgbotapi.ReplyKeyboardMarkup, error) {
	replyMarkup := &tgbotapi.ReplyKeyboardMarkup{ResizeKeyboard: true}
	query := `SELECT k.text FROM keyboard k 
    			INNER JOIN reply_markup rm ON k.id=rm.keyboard_id
    			WHERE rm.message_id=? AND k.lang=?`

	result, err := s.db.Query(query, msgID, lang)
	if err != nil {
		return nil, fmt.Errorf("can't get reply markup: %w", err)
	}

	keyboard, err := s.getKeyboardButtons(result)
	if err != nil {
		return nil, fmt.Errorf("can't get reply markup: %w", err)
	}
	replyMarkup.Keyboard = s.formatKeyboardButtons(keyboard)

	return replyMarkup, nil
}

func (s *storage) GetKeyboardButtons(msgID int64, lang repository.Lang) ([]tgbotapi.KeyboardButton, error) {
	query := `SELECT k.text FROM keyboard k
    			INNER JOIN reply_markup rm ON k.id=rm.keyboard_id 
              	WHERE rm.message_group_id=? AND k.lang=?`

	result, err := s.db.Query(query, msgID, lang)
	if err != nil {
		return nil, fmt.Errorf("can't get reply markup: %w", err)
	}

	return s.getKeyboardButtons(result)
}

func (s *storage) getKeyboardButtons(result *sql.Rows) ([]tgbotapi.KeyboardButton, error) {
	buttons := make([]tgbotapi.KeyboardButton, 0)

	for result.Next() {
		temp := tgbotapi.KeyboardButton{}
		err := result.Scan(&temp.Text)
		if err != nil {
			return nil, fmt.Errorf("can't get keyboard buttons: %w", err)
		}
		buttons = append(buttons, temp)
	}

	return buttons, nil
}

func (s *storage) formatKeyboardButtons(buttons []tgbotapi.KeyboardButton) [][]tgbotapi.KeyboardButton {
	var rows [][]tgbotapi.KeyboardButton
	var row []tgbotapi.KeyboardButton

	for i, button := range buttons {
		if i > 0 && i%3 == 0 {
			rows = append(rows, row)
			row = []tgbotapi.KeyboardButton{}
		}
		row = append(row, button)
	}

	if len(row) > 0 {
		rows = append(rows, row)
	}

	return rows
}

func (s *storage) CreateKeyboard(keyboard *repository.Keyboard) error {
	query := `INSERT INTO keyboard(text, lang) VALUES(?, ?)`
	exec, err := s.db.Exec(query, keyboard.Text, keyboard.Lang)
	if err != nil {
		return fmt.Errorf("creating keyboard error: %w", err)
	}

	id, err := exec.LastInsertId()
	if err != nil {
		fmt.Errorf("creating keyboard error: %w", err)
	}

	keyboard.ID = id
	return nil
}

func (s *storage) CreateKeyboards(keyboards []*repository.Keyboard) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("can't begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	insertKeyboardQuery := `INSERT INTO keyboard (text, lang) VALUES (?,?)`
	for _, v := range keyboards {
		if v.Text == "" || v.Lang == "" {
			return fmt.Errorf("error creeating keyboard group: keyboard.Text is empty")
		}

		exec, err := tx.Exec(insertKeyboardQuery, v.Text, v.Lang)
		if err != nil {
			return fmt.Errorf("insert keyboard was failed: %w", err)
		}

		keyboardId, err := exec.LastInsertId()
		if err != nil {
			return fmt.Errorf("create keyboard group can't get last insert index: %w", err)
		}
		v.ID = keyboardId
	}

	return nil
}

func (s *storage) CreateReplyMarkup(messageID, keyboardID int64) error {
	query := `INSERT INTO reply_markup (message_id, keyboard_id) VALUES (?, ?)`
	if _, err := s.db.Exec(query, messageID, keyboardID); err != nil {
		return fmt.Errorf("can't create reply markup: %w", err)
	}

	return nil
}

func (s *storage) AddFileForMessage(fileInfo repository.FileInfo) error {
	query := `INSERT INTO file (message_id, file_name, file_type) VALUES (?,?,?)`
	if _, err := s.db.Exec(query, fileInfo.MessageID, fileInfo.Name, string(fileInfo.Type)); err != nil {
		return fmt.Errorf("can't add file for message: %w", err)
	}

	return nil
}

func (s *storage) GetFilesInfoOfMessage(messageID int64) ([]repository.FileInfo, error) {
	query := `SELECT f.file_name, f.file_type FROM file f INNER JOIN message ON message.id = f.message_id WHERE message.id=?`
	result, err := s.db.Query(query, messageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get files of message: %w: %w", err)
	}

	return s.getFilesInfoOfMessage(result)
}

func (s *storage) getFilesInfoOfMessage(result *sql.Rows) ([]repository.FileInfo, error) {
	filesInfo := make([]repository.FileInfo, 0)
	for result.Next() {
		temp := repository.FileInfo{}
		err := result.Scan(&temp.Name, &temp.Type)
		if err != nil {
			return nil, fmt.Errorf("can't get files info of message: %w", err)
		}
		filesInfo = append(filesInfo, temp)
	}
	return filesInfo, nil
}

func (s *storage) UpdateMessage(msg *repository.Message, text string) (err error) {
	query := `UPDATE message SET text=? WHERE id=?`

	if _, err = s.db.Exec(query, text, msg.ID); err != nil {
		return err
	}

	msg.Text = text
	return nil
}
