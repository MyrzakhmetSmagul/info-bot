package repository

import "errors"

type ChatInfo struct {
	ChatID    int64
	Active    bool
	Lang      Lang
	State     State
	PrevState State
	CMD       bool
}

type Message struct {
	ID        int64
	Trigger   string
	Text      string
	Lang      Lang
	State     State
	PrevState State
}

type Keyboard struct {
	ID   int64
	Text string
	Lang Lang
}

type FileInfo struct {
	ID        int64
	MessageID int64
	Name      string
	Type      File
	Content   []byte
}

type Lang string

const (
	Kz Lang = "kz"
	Ru Lang = "ru"
	En Lang = "en"
)

type File string

const (
	Photo    File = "photo"
	Document File = "doc"
	Video    File = "video"
)

var (
	ErrUndefinedFileType = errors.New("undefined file type")
	ErrUndefinedLanguage = errors.New("undefined language")
	ErrUndefinedState    = errors.New("undefined state")
)

type State struct {
	ID   int
	Name string
}

func (s State) Equals(state State) bool {
	return s.Name == state.Name
}
