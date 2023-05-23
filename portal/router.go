package portal

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *Portal) Run(port string) error {
	router := gin.Default()
	router.LoadHTMLGlob("./templates/*")
	router.Static("/assets", "./assets")
	create := router.Group("/create")
	// GET methods
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
			messageGroups, err := p.getAllMessageGroup()
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			c.HTML(http.StatusOK, "createTransition.html", gin.H{
				"States":        states,
				"MessageGroups": messageGroups,
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
		create.GET("/add/file", func(c *gin.Context) {
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
			c.HTML(http.StatusOK, "addFile.html", gin.H{
				"States":        states,
				"MessageGroups": messageGroups,
			})
		})
	}
	// POST methods
	{
		create.POST("/message-group", p.createMsgGroup)
		create.POST("/state", p.createState)
		create.POST("/transition", p.createTransition)
		create.POST("/reply-markup", p.createReplyMarkup)
		create.POST("/add/file", p.addFileToMsgGroup)
	}
	router.Run(":" + port)
	return nil
}
