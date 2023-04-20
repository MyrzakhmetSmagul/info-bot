package telegram

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-bot/client/telegram"
	"tg-bot/events"
	"tg-bot/repository"
	file_manager "tg-bot/repository/file-manager"
)

type Processor struct {
	tg          telegram.Client
	storage     repository.Repository
	fileManager file_manager.FileManager
	offset      int
}

type Meta struct {
	ChatID   int64
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client telegram.Client, storage repository.Repository, manager file_manager.FileManager) *Processor {
	return &Processor{
		tg:          client,
		storage:     storage,
		fileManager: manager,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get events: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, v := range updates {
		res = append(res, event(v))
	}

	p.offset = updates[len(updates)-1].UpdateID + 1

	return res, nil
}

func (p Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return fmt.Errorf("can't process message: %w", ErrUnknownEventType)
	}
}

func (p Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}

	if err := p.middleware(meta.ChatID, event.Text, meta.Username); err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("can't get meta: %w", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd tgbotapi.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.UserName,
		}
	}

	return res
}

func fetchType(upd tgbotapi.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(upd tgbotapi.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}
