package repository

import "errors"

const (
	Photo    = "photo"
	Video    = "video"
	Document = "doc"
	Kz       = "kz"
	Ru       = "ru"
	En       = "en"
)

var (
	ErrUndefinedFileType = errors.New("undefined file type")
	ErrStateNotFound     = errors.New("state not found")
)

type State struct {
	ID   int
	Name string
}

type Message struct {
	ID             int
	MessageTrigger string
	Text           string
	Language       string
	StateID        int
}

type Keyboard struct {
	ID     int
	KzText string
	RuText string
	EnText string
}

type ReplyMarkup struct {
	ID         int
	MessageID  int
	KeyboardID int
}

type Command struct {
	ID          int
	Name        string
	Description string
}

type File struct {
	ID        int
	MessageID int
	FileName  string
	FileType  string
	Content   []byte
}

type Chat struct {
	ChatID   int64
	Active   bool
	Language string
	StateID  int
}

type Transition struct {
	ID        int
	Name      string
	FromState State
	ToState   State
	Message   Message
}
