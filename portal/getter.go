package portal

import (
	"fmt"
	"tg-bot/repository"
)

func (p *Portal) getAllMessageGroup() ([]repository.MessageGroup, error) {
	messageGroups, err := p.repository.GetAllMessageGroups()
	if err != nil {
		return nil, fmt.Errorf("portal.getAllMessageGroup failed: %w", err)
	}

	for i := 0; i < len(messageGroups); i++ {
		messageGroups[i].KzMsg, err = p.getMessageByID(messageGroups[i].KzMsg.ID)
		if err != nil {
			return nil, fmt.Errorf("portal.getAllMessageGroup failed: %w", err)
		}
		messageGroups[i].RuMsg, err = p.getMessageByID(messageGroups[i].RuMsg.ID)
		if err != nil {
			return nil, fmt.Errorf("portal.getAllMessageGroup failed: %w", err)
		}
		messageGroups[i].EnMsg, err = p.getMessageByID(messageGroups[i].EnMsg.ID)
		if err != nil {
			return nil, fmt.Errorf("portal.getAllMessageGroup failed: %w", err)
		}
	}

	return messageGroups, nil
}

func (p *Portal) getMessageByID(id int) (repository.Message, error) {
	msg, err := p.repository.GetMessageByID(id)
	if err != nil {
		return repository.Message{}, fmt.Errorf("portal.getMessageByID failed")
	}

	return msg, nil
}

func (p *Portal) getAllStates() ([]repository.State, error) {
	states, err := p.repository.GetAllStates()
	if err != nil {
		return nil, fmt.Errorf("portal.getAllStates failed: %w", err)
	}

	return states, nil
}
