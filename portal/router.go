package portal

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (p *Portal) Run() error {
	router := gin.Default()
	router.LoadHTMLGlob("./templates/*")
	router.Static("/assets", "./assets")
	create := router.Group("/create")
	{
		create.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "create.html", nil)
		})
		create.GET("/message-group", func(c *gin.Context) {
			c.HTML(http.StatusOK, "createMsgGroup.html", nil)
		})
		create.GET("/state", func(c *gin.Context) {
			c.HTML(http.StatusOK, "createState.html", nil)
		})
		create.GET("/transition", func(c *gin.Context) {
			states, err := p.getAllStates()
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			c.HTML(http.StatusOK, "createTransition.html", gin.H{
				"States": states,
			})
		})
		create.GET("/reply-markup", func(c *gin.Context) {
			states, err := p.getAllStates()
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			messageGroups, err := p.getAllMessageGroup()
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			c.HTML(http.StatusOK, "createReplyMarkup.html", gin.H{
				"States":        states,
				"MessageGroups": messageGroups,
			})

		})
	}
	{
		create.POST("/message-group", p.createMsgGroup)
		create.POST("/state", p.createState)
		create.POST("/transition", p.createTransition)
		create.POST("/reply-markup", p.createReplyMarkup)
	}
	router.Run()
	return nil
}
