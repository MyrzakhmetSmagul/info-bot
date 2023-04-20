package repository

import "errors"

type ChatInfo struct {
	ChatID int64
	Active bool
	Lang   Lang
	State  State
}

type Message struct {
	ID             int64
	MessageTrigger string
	KzText         string
	RuText         string
	EnText         string
	State          State
}

type MessageWithLang struct {
	ID    int64
	Text  string
	State State
}

type Keyboard struct {
	ID     int64
	KzText string
	RuText string
	EnText string
	State  State
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

type State int

const (
	DefaultState     State = 0
	ChangeLangState  State = 1
	AskQuestionState State = 2
)

func (s State) Equals(state State) bool {
	return s == state
}
