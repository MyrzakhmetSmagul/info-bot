package portal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tg-bot/repository"
)

func (p *Portal) createMsgGroup(c *gin.Context) {
	languages := []string{"kz", "ru", "en"}
	messages := make([]repository.Message, 0)
	for _, lang := range languages {
		message := repository.Message{}
		message.MsgTrigger = c.PostForm(lang + "MsgTrigger")
		message.Text = c.PostForm(lang + "Text")
		message.Lang = lang
		err := p.repository.CreateMessage(&message)
		if err != nil {
			log.Printf("portal.CreateMsgGroup failed: %w", err)
			abortErr := p.abortCreateMsgGroup(messages)
			if abortErr != nil {
				err = fmt.Errorf("portal.CreateMsgGroup failed: %w\n%w", err, abortErr)
			}
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		messages = append(messages, message)
	}
	messageGroup := repository.MessageGroup{KzMsg: messages[0], RuMsg: messages[1], EnMsg: messages[2]}
	err := p.repository.CreateMessageGroup(&messageGroup)
	if err != nil {
		log.Printf("portal.CreateMsgGroup failed: %w", err)
		abortErr := p.abortCreateMsgGroup(messages)
		if abortErr != nil {
			err = fmt.Errorf("portal.CreateMsgGroup failed: %w\n%w", err, abortErr)
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, "/create/")
}

func (p *Portal) abortCreateMsgGroup(messages []repository.Message) error {
	for _, message := range messages {
		err := p.repository.DeleteMessage(message.ID)
		if err != nil {
			return fmt.Errorf("portal.AbortCreateMsgGroup failed: %w", err)
		}
	}
	return nil
}

func (p *Portal) createState(c *gin.Context) {
	state := repository.State{}
	state.Name = c.PostForm("stateName")
	if state.Name == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("createState: %s", http.StatusText(http.StatusBadRequest)))
		return
	}
	err := p.repository.CreateState(&state)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, "/create/")
}

func (p *Portal) createTransition(c *gin.Context) {
	var err error

	transition := repository.Transition{}
	transition.MsgTrigger = c.PostForm("MsgTrigger")

	temp := c.PostForm("fromState")
	transition.FromStateID, err = strconv.Atoi(temp)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	temp = c.PostForm("toState")
	transition.ToStateID, err = strconv.Atoi(temp)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = p.repository.CreateTransition(&transition)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusFound, "/create")
}

func (p *Portal) createReplyMarkup(c *gin.Context) {
	var err error
	rm := repository.ReplyMarkup{}
	language := c.PostForm("language")
	temp := c.PostForm(language + "MsgGroup")
	rm.MsgID, err = strconv.Atoi(temp)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	temp = c.PostForm("state")
	rm.StateID, err = strconv.Atoi(temp)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = p.repository.CreateReplyMarkup(&rm)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusFound, "/create/")
}
