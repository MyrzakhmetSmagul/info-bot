package db_mysql

import (
	"fmt"
	"tg-bot/repository"
)

func (s storage) CreateMessageGroup(msgGroup *repository.MessageGroup) error {
	query := `INSERT INTO message_group(kz_msg_id, ru_msg_id, en_msg_id) VALUES (?, ?, ?)`

	exec, err := s.db.Exec(query, msgGroup.KzMsg.ID, msgGroup.RuMsg.ID, msgGroup.EnMsg.ID)
	if err != nil {
		return fmt.Errorf("create message group failed: %w", err)
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return fmt.Errorf("create message group failed: %w", err)
	}

	msgGroup.ID = int(id)
	return nil
}

func (s storage) GetMessageGroup(msgID int, lang string) (repository.MessageGroup, error) {
	query := fmt.Sprintf("SELECT id, kz_msg_id, ru_msg_id, en_msg_id FROM message_group WHERE %s_msg_id=?", lang)

	row := s.db.QueryRow(query, msgID)

	msgGroup := repository.MessageGroup{}

	err := row.Scan(&msgGroup.ID, &msgGroup.KzMsg.ID, &msgGroup.RuMsg.ID, &msgGroup.EnMsg.ID)
	if err != nil {
		return repository.MessageGroup{}, fmt.Errorf("get message group failed: %w", err)
	}

	return msgGroup, nil
}

func (s storage) GetMessageGroupByID(msgGroupID int) (repository.MessageGroup, error) {
	query := "SELECT id, kz_msg_id, ru_msg_id, en_msg_id FROM message_group WHERE id=?"

	row := s.db.QueryRow(query, msgGroupID)

	msgGroup := repository.MessageGroup{}

	err := row.Scan(&msgGroup.ID, &msgGroup.KzMsg.ID, &msgGroup.RuMsg.ID, &msgGroup.EnMsg.ID)
	if err != nil {
		return repository.MessageGroup{}, fmt.Errorf("get message group failed: %w", err)
	}

	return msgGroup, nil
}

func (s storage) GetAllMessageGroups() ([]repository.MessageGroup, error) {
	query := `SELECT id, kz_msg_id, ru_msg_id, en_msg_id FROM message_group`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("get all messageGroup failed:%w", err)
	}

	msgGroups := make([]repository.MessageGroup, 0)
	defer rows.Close()

	for rows.Next() {
		msgGroup := repository.MessageGroup{}
		err = rows.Scan(&msgGroup.ID, &msgGroup.KzMsg.ID, &msgGroup.RuMsg.ID, &msgGroup.EnMsg.ID)
		if err != nil {
			return nil, fmt.Errorf("get all messageGroup failed:%w", err)
		}
		msgGroups = append(msgGroups, msgGroup)
	}

	return msgGroups, nil
}

func (s storage) DeleteMessageGroup(msgGroupID int) error {
	query := `DELETE FROM message_group WHERE id=?`

	//you should check how many rows affected
	//if affected rows is zero, return custom error which describe it
	_, err := s.db.Exec(query, msgGroupID)
	if err != nil {
		return fmt.Errorf("delete messageGroup failed:%w", err)
	}

	return nil
}
