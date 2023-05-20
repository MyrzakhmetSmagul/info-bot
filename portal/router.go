package portal

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (p *Portal) Run() error {
	router := gin.Default()
	create := router.Group("/create")
	{
		create.GET("/message")
	}
	return nil
}

func Run() {
	router := gin.Default()
	router.LoadHTMLGlob("./templates/*")
	router.Static("/assets", "./assets")
	create := router.Group("/create")
	{
		create.GET("/state", func(c *gin.Context) {
			c.HTML(http.StatusOK, "create.html", nil)
		})
		create.GET("/message-group", func(c *gin.Context) {
			c.HTML(http.StatusOK, "create.html", nil)
		})
	}
	router.Run()
}
