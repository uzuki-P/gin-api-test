package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
)

var db = make(map[string]string)

type ItemEntity struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type OnboardingItemEntity struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	ImageUrl string `json:"image_url"`
}

var fakeCat = "https://loremflickr.com/500/500/cat"

type ItemEntities []ItemEntity

type AuthTokenEntity struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type QueryRecords struct {
	Page  int `form:"page"   default:"1"`
	Limit int `form:"limit"  default:"10"`
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

	r.GET("/null", func(c *gin.Context) {
		c.String(http.StatusOK, "null")
	})

	r.GET("/int", func(c *gin.Context) {
		c.String(http.StatusOK, "10")
	})

	r.GET("/empty-json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	r.GET("/string/error", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing string"})
	})

	r.GET("/list-string", func(c *gin.Context) {
		res := []string{
			"satu", "dua", "tiga",
		}

		c.JSON(http.StatusOK, res)
	})

	r.GET("/list-int", func(c *gin.Context) {
		res := []int{
			1, 2, 3,
		}

		c.JSON(http.StatusOK, res)
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
		var query QueryRecords
		if err := c.ShouldBindQuery(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		page := query.Page
		limit := query.Limit
		totalPage := 13 * limit

		// res := []ItemEntity{
		// 	{ID: 1, Name: faker.Name()},
		// 	{ID: 2, Name: faker.Name()},
		// 	{ID: 3, Name: faker.Name()},
		// }

		res := make([]ItemEntity, limit)
		for i := 0; i < limit; i++ {
			res[i] = ItemEntity{ID: (i + 1) + (limit * (page - 1)), Name: faker.Name()}
		}

		c.Header("records-page", strconv.Itoa(page))
		c.Header("records-total-page", strconv.Itoa(totalPage))
		c.Header("records-per-page", strconv.Itoa(limit))

		c.JSON(http.StatusOK, res)
	})

	r.GET("/onboarding", func(c *gin.Context) {
		res := []OnboardingItemEntity{
			{ID: 1, Title: faker.Name(), Subtitle: faker.Username(), ImageUrl: fakeCat + "?id=1"},
			{ID: 2, Title: faker.Name(), Subtitle: faker.Username(), ImageUrl: fakeCat + "?id=2"},
			{ID: 3, Title: faker.Name(), Subtitle: faker.Username(), ImageUrl: fakeCat + "?id=3"},
			{ID: 4, Title: faker.Name(), Subtitle: faker.Username(), ImageUrl: fakeCat + "?id=4"},
			{ID: 5, Title: faker.Name(), Subtitle: faker.Username(), ImageUrl: fakeCat + "?id=5"},
			{ID: 6, Title: faker.Name(), Subtitle: faker.Username(), ImageUrl: fakeCat + "?id=6"},
		}

		c.JSON(http.StatusOK, res)
	})

	r.GET("/error-503", func(c *gin.Context) {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Service Unavailable"})
	})

	r.POST("/login", func(c *gin.Context) {
		user := c.PostForm("user")
		pass := c.PostForm("pass")

		if user == "" || pass == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing user or pass"})
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
