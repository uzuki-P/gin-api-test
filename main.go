package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
)

var db = make(map[string]string)

type ItemEntity struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ItemEntities []ItemEntity

type AuthTokenEntity struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func setupRouter() *gin.Engine {
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/string", func(c *gin.Context) {
		c.JSON(http.StatusOK, faker.Name())
	})

	r.GET("/string/error", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing string"})
	})

	r.GET("/list-item-entity", func(c *gin.Context) {
		res := []ItemEntity{
			{ID: 1, Name: faker.Name()},
			{ID: 2, Name: faker.Name()},
			{ID: 3, Name: faker.Name()},
		}

		c.JSON(http.StatusOK, res)
	})

	r.GET("/list-zero", func(c *gin.Context) {
		res := []ItemEntity{}

		c.JSON(http.StatusOK, res)
	})

	r.GET("/item-entity", func(c *gin.Context) {
		c.JSON(http.StatusOK, ItemEntity{ID: 1, Name: "test1"})
	})

	r.GET("/records-list-item-entity", func(c *gin.Context) {
		res := []ItemEntity{
			{ID: 1, Name: faker.Name()},
			{ID: 2, Name: faker.Name()},
			{ID: 3, Name: faker.Name()},
		}

		c.Header("records-page", "1")
		c.Header("records-total-page", "10")
		c.Header("records-per-page", "5")

		c.JSON(http.StatusOK, res)
	})

	r.GET("/error-503", func(c *gin.Context) {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Service Unavailable"})
	})

	r.POST("/login", func(c *gin.Context) {
		user := c.PostForm("user")
		pass := c.PostForm("pass")

		if user == "" || pass == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "missing user or pass"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  "acc Token",
			"refresh_token": "ref Token",
			"expires_in":    100,
		})
	})

	r.POST("/refresh-token", func(c *gin.Context) {
		refreshToken := c.PostForm("refresh_token")

		if refreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "missing token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  "acc Token",
			"refresh_token": "ref Token",
			"expires_in":    100,
		})
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
