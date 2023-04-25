package db_mysql

import (
	"fmt"
	"tg-bot/repository"
)

func (s *storage) CreateQuestion(question *repository.QuestionInfo) error {
	query := `INSERT INTO question(chat_id, question) VALUES (?, ?)`
	exec, err := s.db.Exec(query, question.ChatID, question.Question)
	if err != nil {
		return fmt.Errorf("create question error: %w", err)
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return fmt.Errorf("create question error: %w", err)
	}
	question.ID = int(id)
	return nil
}

func (s *storage) AnswerToQuestion(question *repository.QuestionInfo, answer string) error {
	query := `UPDATE question SET answer=? WHERE id=?`
	if _, err := s.db.Exec(query, answer, question.ID); err != nil {
		return err
	}

	return nil
}
