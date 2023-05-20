package db_mysql

import (
	"fmt"
	"tg-bot/repository"
)

func (s storage) CreateTransition(transition *repository.Transition) error {
	query := `INSERT INTO transition(msg_trigger, from_state, toState) VALUES (?, ?, ?)`

	exec, err := s.db.Exec(query, transition.MsgTrigger, transition.FromStateID, transition.ToStateID)
	if err != nil {
		return fmt.Errorf("create transition was failed: %w", err)
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return fmt.Errorf("create transition was failed: %w", err)
	}

	transition.ID = int(id)

	return nil
}

func (s storage) GetTransition(fromStateID int, msgTrigger string) (repository.Transition, error) {
	query := `SELECT id, to_state FROM transition WHERE from_state=? and msg_trigger=?`

	transition := repository.Transition{FromStateID: fromStateID, MsgTrigger: msgTrigger}
	err := s.db.QueryRow(query, fromStateID, msgTrigger).Scan(&transition.ID, &transition.ToStateID)
	if err != nil {
		return repository.Transition{}, fmt.Errorf("get transition was failed: %w", err)
	}

	return transition, nil
}

func (s storage) GetAllTransitions() ([]repository.Transition, error) {
	query := `SELECT id, msg_trigger, from_state, to_state FROM transition`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("get all transitions was failed: %w", err)
	}

	transitions := make([]repository.Transition, 0)
	defer rows.Close()

	for rows.Next() {
		transition := repository.Transition{}

		err = rows.Scan(&transition.ID, &transition.MsgTrigger, &transition.FromStateID, &transition.ToStateID)
		if err != nil {
			return nil, fmt.Errorf("get all transitions was failed: %w", err)
		}

		transitions = append(transitions, transition)
	}

	return transitions, nil
}

func (s storage) DeleteTransition(transitionID int) error {
	query := `DELETE FROM transition WHERE id=?`

	//you should check how many rows affected
	//if affected rows is zero, return custom error which describe it
	_, err := s.db.Exec(query, transitionID)
	if err != nil {
		return fmt.Errorf("delete transition was failed: %w", err)
	}

	return nil
}
