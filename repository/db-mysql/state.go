package db_mysql

import (
	"database/sql"
	"fmt"
	"tg-bot/repository"
)

func (s *storage) AddState(state *repository.State) error {
	query := "INSERT INTO state (name) VALUES (?)"
	result, err := s.db.Exec(query, state.Name)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	state.ID = int(id)
	return nil
}

func (s *storage) GetStateByID(id int) (*repository.State, error) {
	query := "SELECT id, name FROM state WHERE id = ?"
	state := &repository.State{}

	err := s.db.QueryRow(query, id).Scan(&state.ID, &state.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrStateNotFound
		}
		return nil, err
	}

	return state, nil
}

func (s *storage) GetAllStates() ([]*repository.State, error) {
	rows, err := s.db.Query("SELECT id, name FROM state")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	states := make([]*repository.State, 0)
	for rows.Next() {
		state := &repository.State{}
		if err := rows.Scan(&state.ID, &state.Name); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		states = append(states, state)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return states, nil
}
