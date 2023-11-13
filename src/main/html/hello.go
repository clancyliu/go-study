package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"github.com/go-session/session"
	"net/http"
)

func main() {
	//r := gin.Default()
	//
	//// 设置 Session 中间件
	//store := cookie.NewStore([]byte("your-secret-key"))
	//r.Use(sessions.Sessions("mysession", store))
	//
	//// 配置跨域中间件
	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:63342"} // 允许的域名
	//config.AllowCredentials = true                           // 允许发送和接收 Cookie
	//r.Use(cors.New(config))
	//
	//r.GET("/set-session", func(c *gin.Context) {
	//	session := sessions.Default(c)
	//	session.Set("user", "exampleUser")
	//	session.Save()
	//	c.String(http.StatusOK, "Session set")
	//})
	//
	//r.GET("/get-session", func(c *gin.Context) {
	//	session := sessions.Default(c)
	//	user := session.Get("user")
	//	c.String(http.StatusOK, "User: %v", user)
	//})
	//
	//r.Run(":8080")

	r := gin.Default()
	r.Use(ginsession.New(session.SetEnableSetCookie(true), session.SetSecure(true)))
	// 配置跨域中间件
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:63342"} // 允许的域名
	config.AllowCredentials = true                           // 允许发送和接收 Cookie
	r.Use(cors.New(config))

	r.GET("/set-session", func(c *gin.Context) {
		store := ginsession.FromContext(c)
		store.Set("foo", "bar")
		err := store.Save()
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.String(http.StatusOK, "Session set")
	})

	r.GET("/get-session", func(c *gin.Context) {
		store := ginsession.FromContext(c)
		user, _ := store.Get("foo")
		c.String(http.StatusOK, "User: %v", user)
	})

	r.Run(":8080")
}
