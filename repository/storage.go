package repository

type Repository interface {
	MessageRepository
	KeyboardRepository
	ReplyMarkupRepository
	CommandRepository
	FileRepository
	ChatRepository
	StateRepository
	TransitionRepository
}
type MessageRepository interface {
	CreateMessage(message *Message) error
	UpdateMessage(message *Message) error
	DeleteMessage(messageID int) error
	GetMessageByID(messageID int) (*Message, error)
	GetMessageByTriggerAndState(messageTrigger string, stateID int) (*Message, error)
	GetAllMessages() ([]*Message, error)
}

type KeyboardRepository interface {
	CreateKeyboard(keyboard *Keyboard) error
	UpdateKeyboard(keyboard *Keyboard) error
	DeleteKeyboard(keyboardID int) error
	GetKeyboardByID(keyboardID int) (*Keyboard, error)
	GetAllKeyboards() ([]*Keyboard, error)
}

type ReplyMarkupRepository interface {
	CreateReplyMarkup(replyMarkup *ReplyMarkup) error
	UpdateReplyMarkup(replyMarkup *ReplyMarkup) error
	DeleteReplyMarkup(replyMarkupID int) error
	GetReplyMarkupByID(replyMarkupID int) (*ReplyMarkup, error)
	GetAllReplyMarkups() ([]*ReplyMarkup, error)
}

type CommandRepository interface {
	CreateCommand(command *Command) error
	UpdateCommand(command *Command) error
	DeleteCommand(commandID int) error
	GetCommandByID(commandID int) (*Command, error)
	GetCommandByName(commandName string) (*Command, error)
	GetAllCommands() ([]*Command, error)
}

type FileRepository interface {
	AddFile(file *File) error
	DeleteFile(fileID int) error
	GetFileByID(fileID int) (*File, error)
	GetFilesByMessageID(messageID int) ([]*File, error)
	GetAllFiles() ([]*File, error)
}

type ChatRepository interface {
	CreateChat(chat *Chat) error
	EnableChat(chatID int64) error
	DeleteChat(chatID int64) error
	DisableChat(chatID int64) error
	ChangeLang(info Chat) error
	ChangeState(info Chat) error
	GetChatByID(chatID int64) (*Chat, error)
	GetAllChats() ([]*Chat, error)
}

type StateRepository interface {
	// Метод для добавления нового состояния
	AddState(state *State) error
	// Метод для получения состояния по его ID
	GetStateByID(id int) (*State, error)
	// Метод для получения списка всех состояний
	GetAllStates() ([]*State, error)
}

type TransitionRepository interface {
	// Метод для добавления нового перехода
	AddTransition(transition *Transition) error
	// Метод для получения списка всех переходов из указанного состояния
	GetTransitionsFromState(stateID int) ([]*Transition, error)
	// Метод для получения списка всех переходов в указанное состояние
	GetTransitionsToState(stateID int) ([]*Transition, error)
	// Метод для получения списка всех переходов
	GetAllTransitions() ([]*Transition, error)
}
