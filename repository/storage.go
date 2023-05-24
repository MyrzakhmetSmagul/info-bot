package repository

type Repository interface {
	MessageRepository
	MessageGroupRepository
	ReplyMarkupRepository
	CommandRepository
	FileRepository
	ChatRepository
	StateRepository
	TransitionRepository
}
type MessageRepository interface {
	CreateMessage(msg *Message) error
	UpdateMessage(msg Message) error
	GetMessage(trigger, lang string) (Message, error)
	GetMessageByID(id int) (Message, error)
	GetAllMessages() ([]Message, error)
	DeleteMessage(msgID int) error
}

type MessageGroupRepository interface {
	CreateMessageGroup(msgGroup *MessageGroup) error
	GetMessageGroup(msgID int, lang string) (MessageGroup, error)
	GetMessageGroupByID(msgGroupID int) (MessageGroup, error)
	GetAllMessageGroups() ([]MessageGroup, error)
	DeleteMessageGroup(msgGroupID int) error
}

type StateRepository interface {
	CreateState(state *State) error
	GetState(id int) (State, error)
	GetAllStates() ([]State, error)
	DeleteStates(stateID int) error
}

type TransitionRepository interface {
	CreateTransition(transition *Transition) error
	GetTransition(fromStateID int, msgGroupID int) (Transition, error)
	GetAllTransitions() ([]Transition, error)
	DeleteTransition(transitionID int) error
}

type ReplyMarkupRepository interface {
	CreateReplyMarkup(replyMarkup *ReplyMarkup) error
	GetReplyMarkupByID(stateID int) (ReplyMarkup, error)
	GetReplyMarkupsOfState(stateID int) ([]ReplyMarkup, error)
	DeleteReplyMarkup(replyMarkupID int) error
}

type FileRepository interface {
	AddFileToMessage(file *File) error
	GetFileByID(fileID int) (File, error)
	GetFilesOfMsgGroup(msgGroupID int) ([]File, error)
	DeleteFile(fileID int) error
}

type ChatRepository interface {
	CreateChat(chat Chat) error
	EnableChat(chatID int64) error
	DisableChat(chatID int64) error
	EnableCmd(chatID int64) error
	DisableCmd(chatID int64) error
	ChangeChatLang(chatID int64, newLang string) error
	ChangeChatState(chatID int64, newState int) error
	GetChat(chatID int64) (Chat, error)
	GetAllChats() ([]Chat, error)
	DeleteChat(chatID int64) error
}

type CommandRepository interface {
	CreateCommand(command *Command) error
	UpdateCommand(command Command) error
	GetCommand(commandID int) (Command, error)
	GetAllCommands() ([]Command, error)
	DeleteCommand(commandID int) error
}
