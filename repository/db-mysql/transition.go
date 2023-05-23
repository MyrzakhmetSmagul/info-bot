package db_mysql

import (
	"fmt"
	"tg-bot/repository"
)

func (s storage) CreateTransition(transition *repository.Transition) error {
	query := `INSERT INTO transition (msg_group_id, from_state, to_state)
SELECT ?, ?, ?
FROM dual
WHERE NOT EXISTS (
    SELECT 1
    FROM transition
    WHERE msg_group_id = ?
        AND from_state = ?
        AND to_state = ?
)`

	exec, err := s.db.Exec(query, transition.MsgGroup.ID, transition.FromState.ID, transition.ToState.ID, transition.MsgGroup.ID, transition.FromState.ID, transition.ToState.ID)
	if err != nil {
		return fmt.Errorf("create transition failed: %w", err)
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return fmt.Errorf("create transition failed: %w", err)
	}

	transition.ID = int(id)

	return nil
}

func (s storage) GetTransition(fromStateID int, msgGroupID int) (repository.Transition, error) {
	query := `SELECT id, to_state FROM transition WHERE from_state=? and msg_trigger=?`

	transition := repository.Transition{FromState: repository.State{ID: fromStateID}, MsgGroup: repository.MessageGroup{ID: msgGroupID}}
	err := s.db.QueryRow(query, fromStateID, msgGroupID).Scan(&transition.ID, &transition.ToState.ID)
	if err != nil {
		return repository.Transition{}, fmt.Errorf("get transition failed: %w", err)
	}

	return transition, nil
}

func (s storage) GetAllTransitions() ([]repository.Transition, error) {
	query := `SELECT id, msg_trigger, from_state, to_state FROM transition`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("get all transitions failed: %w", err)
	}

	transitions := make([]repository.Transition, 0)
	defer rows.Close()

	for rows.Next() {
		transition := repository.Transition{}

		err = rows.Scan(&transition.ID, &transition.MsgGroup.ID, &transition.FromState.ID, &transition.ToState.ID)
		if err != nil {
			return nil, fmt.Errorf("get all transitions failed: %w", err)
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
		return fmt.Errorf("delete transition failed: %w", err)
	}

	return nil
}
