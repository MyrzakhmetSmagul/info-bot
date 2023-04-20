package db_mysql

import (
	"database/sql"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"tg-bot/lib/e"
	"tg-bot/repository"
)

func (s *storage) GetLangMessage(trigger string, lang repository.Lang) (*repository.MessageWithLang, error) {
	msg := &repository.MessageWithLang{}
	query := fmt.Sprintf("SELECT id, %s_text, state FROM message WHERE message_trigger=?", lang)

	if err := s.db.QueryRow(query, trigger).Scan(&msg.ID, &msg.Text, &msg.State); err != nil {
		return nil, e.Wrap("can't get lang message", err)
	}

	return msg, nil
}

func (s *storage) GetMessage(trigger string) (*repository.Message, error) {
	msg := &repository.Message{}
	query := `SELECT id, kz_text, ru_text, en_text, state FROM message WHERE message_trigger=?`

	if err := s.db.QueryRow(query, trigger).Scan(&msg.ID, &msg.KzText, &msg.RuText, &msg.EnText, &msg.State); err != nil {
		return nil, e.Wrap("can't get message", err)
	}

	return msg, nil
}

func (s *storage) CreateMessage(msg repository.Message) error {
	query := `INSERT INTO message(message_trigger, kz_text, ru_text, en_text, state) VALUES (?,?,?,?)`

	if _, err := s.db.Exec(query, msg.MessageTrigger, msg.KzText, msg.RuText, msg.EnText, msg.State); err != nil {
		return e.Wrap("can't create message", err)
	}

	return nil
}

func (s *storage) UpdateMessage(msg repository.Message) (err error) {
	errMsg := "can't updating message"
	defer func() { err = e.WrapIfErr(errMsg, err) }()
	// Проверяем, что сообщение не пустое
	if msg.MessageTrigger == "" {
		return fmt.Errorf("empty message trigger")
	}

	// Проверяем, что хотя бы одно значение текста сообщения было передано
	if msg.KzText == "" && msg.RuText == "" && msg.EnText == "" {
		return fmt.Errorf("empty message text")
	}

	// Строим SQL-запрос для редактирования сообщения
	query := "UPDATE message SET "
	args := make([]interface{}, 0)
	if msg.KzText != "" {
		query += "kz_text=?, "
		args = append(args, msg.KzText)
	}
	if msg.RuText != "" {
		query += "ru_text=?, "
		args = append(args, msg.RuText)
	}
	if msg.EnText != "" {
		query += "en_text=?, "
		args = append(args, msg.EnText)
	}
	query = strings.TrimSuffix(query, ", ")
	query += " WHERE message_trigger=?"
	args = append(args, msg.MessageTrigger)

	// Выполняем SQL-запрос
	if _, err = s.db.Exec(query, args...); err != nil {
		return err
	}

	return nil
}

func (s *storage) GetReplyMarkup(messageID int64, lang repository.Lang) (*tgbotapi.ReplyKeyboardMarkup, error) {
	replyMarkup := &tgbotapi.ReplyKeyboardMarkup{ResizeKeyboard: true}
	query := fmt.Sprintf("SELECT kb.%s_text FROM reply_markup rm INNER JOIN keyboard kb ON kb.id=rm.keyboard_id WHERE rm.message_id=?", lang)

	result, err := s.db.Query(query, messageID)
	if err != nil {
		return nil, e.Wrap("can't get reply markup", err)
	}

	keyboard, err := s.getKeyboardButtons(result)
	if err != nil {
		return nil, e.Wrap("can't get reply markup", err)
	}
	replyMarkup.Keyboard = s.formatKeyboardButtons(keyboard)

	return replyMarkup, nil
}

func (s *storage) GetKeyboardButtons(messageID int64, lang repository.Lang) ([]tgbotapi.KeyboardButton, error) {
	query := fmt.Sprintf("SELECT kb.%s_text FROM reply_markup rm INNER JOIN keyboard kb ON kb.id=rm.keyboard_id WHERE rm.message_id=?", lang)

	result, err := s.db.Query(query, messageID)
	if err != nil {
		return nil, e.Wrap("can't get reply markup", err)
	}

	return s.getKeyboardButtons(result)
}

func (s *storage) getKeyboardButtons(result *sql.Rows) ([]tgbotapi.KeyboardButton, error) {
	buttons := make([]tgbotapi.KeyboardButton, 0)

	for result.Next() {
		temp := tgbotapi.KeyboardButton{}
		err := result.Scan(&temp.Text)
		if err != nil {
			return nil, e.Wrap("can't get keyboard buttons", err)
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

func (s *storage) CreateKeyboard(keyboard repository.Keyboard) (err error) {
	errMsg := "can't create keyboard"

	defer func() { err = e.WrapIfErr(errMsg, err) }()

	if keyboard.KzText == "" || keyboard.RuText == "" || keyboard.EnText == "" {
		return fmt.Errorf("empty keyboard text")
	}

	query := `INSERT INTO keyboard (kz_text,ru_text,en_text) VALUES (?,?,?)`
	if _, err = s.db.Exec(query, keyboard.KzText, keyboard.RuText, keyboard.EnText); err != nil {
		return err
	}

	return nil
}

func (s *storage) CreateReplyMarkup(messageID, keyboardID int64) error {
	query := `INSERT INTO reply_markup (message_id, keyboard_id) VALUES (?, ?)`
	if _, err := s.db.Exec(query, messageID, keyboardID); err != nil {
		return e.Wrap("can't create reply markup", err)
	}

	return nil
}

func (s *storage) AddFileForMessage(fileInfo repository.FileInfo) error {
	query := `INSERT INTO file (message_id, file_name, file_type) VALUES (?,?,?)`
	if _, err := s.db.Exec(query, fileInfo.MessageID, fileInfo.Name, string(fileInfo.Type)); err != nil {
		return e.Wrap("can't add file for message", err)
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
		return nil, fmt.Errorf("can't get files of message: %w", err)
	}

	return s.getFilesInfoOfMessage(result)
}

func (s *storage) getFilesInfoOfMessage(result *sql.Rows) ([]repository.FileInfo, error) {
	filesInfo := make([]repository.FileInfo, 0)
	for result.Next() {
		temp := repository.FileInfo{}
		err := result.Scan(&temp.Name, &temp.Type)
		if err != nil {
			return nil, e.Wrap("can't get files info of message", err)
		}
		filesInfo = append(filesInfo, temp)
	}
	return filesInfo, nil
}
