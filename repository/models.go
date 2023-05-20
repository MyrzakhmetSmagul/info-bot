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
	ErrUndefinedLanguage = errors.New("undefined language")
)

type State struct {
	ID   int
	Name string
}

type Message struct {
	ID         int
	MsgTrigger string
	Text       string
	Lang       string
}

type MessageGroup struct {
	ID    int
	KzMsg Message
	RuMsg Message
	EnMsg Message
}

type ReplyMarkup struct {
	ID      int
	MsgID   int
	StateID int
}

type Command struct {
	ID          int
	Name        string
	Description string
}

type File struct {
	ID         int
	MsgGroupID int
	FileName   string
	FileType   string
	Content    []byte
}

type Chat struct {
	ChatID  int64
	Active  bool
	Lang    string
	StateID int
	CMD     bool
}

type Transition struct {
	ID          int
	FromStateID int
	ToStateID   int
	MsgTrigger  string
}
