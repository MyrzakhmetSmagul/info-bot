package db_mysql

import (
	"fmt"
	"tg-bot/repository"
)

func (s storage) CreateState(state *repository.State) error {
	query := `
INSERT INTO state(name)
SELECT ?
FROM dual
WHERE NOT EXISTS (
    SELECT 1
    FROM state
    WHERE name = ?
)`

	exec, err := s.db.Exec(query, state.Name, state.Name)
	if err != nil {
		return fmt.Errorf("create state failed: %w", err)
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return fmt.Errorf("create state failed: %w", err)
	}

	state.ID = int(id)
	return nil
}

func (s storage) GetState(id int) (repository.State, error) {
	query := `SELECT name FROM state WHERE id=?`
	state := repository.State{ID: id}
	err := s.db.QueryRow(query, id).Scan(&state.Name)
	if err != nil {
		return repository.State{}, fmt.Errorf("get state failed: %w", err)
	}

	return state, nil
}

func (s storage) GetAllStates() ([]repository.State, error) {
	query := `SELECT id, name FROM state`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("get all states failed: %w", err)
	}

	states := make([]repository.State, 0)
	defer rows.Close()

	for rows.Next() {
		state := repository.State{}

		err = rows.Scan(&state.ID, &state.Name)
		if err != nil {
			return nil, fmt.Errorf("get all states failed: %w", err)
		}
		states = append(states, state)
	}

	return states, nil
}

func (s storage) DeleteStates(stateID int) error {
	query := `DELETE FROM state WHERE id=?`

	_, err := s.db.Exec(query, stateID)
	if err != nil {
		return fmt.Errorf("delete state failed: %w", err)
	}

	return nil
}
