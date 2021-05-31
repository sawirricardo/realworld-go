package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	initConfig()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Real world application implemented with Go(Golang) By Ricardo Sawir")
	})
	r.GET("articles", getArticles)
	r.POST("articles", createArticle)
	r.GET("articles/:slug", showArticle)
	r.PUT("articles/:slug", updateArticle)
	r.DELETE("articles/:slug", deleteArticle)

	r.POST("articles/:slug/favorite", favoriteArticle)
	r.DELETE("articles/:slug/favorite", unfavoriteArticle)

	r.GET("articles/:slug/comments", getComments)
	r.POST("articles/:slug/comments", createComment)
	r.DELETE("articles/:slug/comments/:commentId", deleteComment)

	r.POST("users/login", loginUser)
	r.POST("users", createUser)
	r.GET("user", showUser)
	r.PUT("user", updateUser)

	r.GET("profiles/:username", showProfile)
	r.POST("profiles/:username/follow", followProfile)
	r.DELETE("profiles/:username/follow", unfollowProfile)

	r.GET("tags", getTags)
	r.Run(":9001") // listen and serve on 0.0.0.0:8080
}

func initConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getDB() *gorm.DB {

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")), // data source name
		DefaultStringSize:         191,                                                                                                                                                                                                // default size for string fields
		DisableDatetimePrecision:  true,                                                                                                                                                                                               // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                                                                                                                                               // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                                                                                                                                               // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                                                                                                                                              // auto configure based on currently MySQL version
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

type User struct {
	ID              uint64    `json:"-"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	Password        string    `json:"-"`
	Bio             string    `json:"bio"`
	Image           string    `json:"image"`
	EmailVerifiedAt time.Time `json:"-"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

type LoginRequest struct {
	Data LoginUser `json:"user" binding:"required"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required`
}

type UserResource struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

type Profile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type Article struct {
	ID          uint64    `json:"-"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Description string    `json:"description"`
	Slug        string    `json:"slug"`
	UserID      uint64    `json:"-"`
	User        User      `json:"author"`
	Tags        []Tag     `json:"tagList" gorm:"many2many:article_tag"`
	Favoriters  []User    `json:"-" gorm:"many2many:favoriters"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ArticleResource struct {
	Title          string   `json:"title"`
	Body           string   `json:"body"`
	Description    string   `json:"description"`
	Slug           string   `json:"slug"`
	User           User     `json:"author"`
	Tags           []string `json:"tagList"`
	Favorited      bool     `json:"favorited"`
	FavoritesCount uint     `json:"favoritesCount"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

type Comment struct {
	ID        uint64    `json:"id"`
	Body      string    `json:"body"`
	UserID    uint64    `json:"author"`
	ArticleID uint64    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Tag struct {
	ID        uint64    `json:"-"`
	Name      string    `json:"name"`
	Slug      string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func getArticles(c *gin.Context) {
	db := getDB()
	var articles []Article
	var articleCollection []ArticleResource
	db.Preload("User").Preload("Tags").Find(&articles)
	for _, article := range articles {
		var tags []string

		for _, tag := range article.Tags {
			tags = append(tags, tag.Name)
		}

		articleResource := ArticleResource{
			Title:       article.Title,
			Description: article.Description,
			Tags:        tags,
			Body:        article.Body,
			User:        article.User,
			CreatedAt:   article.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   article.UpdatedAt.Format(time.RFC3339),
		}

		articleCollection = append(articleCollection, articleResource)
	}

	c.JSON(200, gin.H{
		"articles":      articleCollection,
		"articlesCount": len(articles),
	})
}

func createArticle(c *gin.Context) {
	c.JSON(200, gin.H{
		"article": "pong",
	})
}

func deleteArticle(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func showArticle(c *gin.Context) {
	slug := c.Param("slug")
	db := getDB()
	var article Article
	db.Where("slug=?", slug).Preload("Tags").First(&article)
	c.JSON(200, gin.H{
		"article": article,
	})
}

func updateArticle(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func favoriteArticle(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func unfavoriteArticle(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func getComments(c *gin.Context) {
	c.JSON(200, gin.H{
		"comments": "pong",
	})
}

func createComment(c *gin.Context) {
	c.JSON(200, gin.H{
		"comment": "pong",
	})
}

func deleteComment(c *gin.Context) {
	c.JSON(200, gin.H{
		"comment": "pong",
	})
}

func loginUser(c *gin.Context) {
	var loginRequest LoginRequest
	if c.BindJSON(&loginRequest) != nil {
		panic("Wrong data")
	}

	var user User
	getDB().First(&user, "email = ?", loginRequest.Data.Email)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Data.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "wrong username or password",
		})
	}
	token, err := CreateToken(user.ID)
	if err != nil {
		panic("Server error")
	}

	c.JSON(200, gin.H{
		"user": &UserResource{
			Username: user.Username,
			Email:    user.Email,
			Bio:      user.Bio,
			Image:    user.Image,
			Token:    token,
		},
	})
}

func CreateToken(userid uint64) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func createUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"comment": "pong",
	})
}

func updateUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"comment": "pong",
	})
}

func showUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"comment": "pong",
	})
}

func showProfile(c *gin.Context) {
	username := c.Param("username")
	var user User
	db := getDB()
	db.Where("username = ?", username).First(&user)
	profile := &Profile{Username: user.Username, Bio: user.Bio, Image: user.Image}
	c.JSON(200, gin.H{
		"profile": profile,
	})
}

func followProfile(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func unfollowProfile(c *gin.Context) {
	c.JSON(200, gin.H{"token": authApi(c)})
	username := c.Param("username")
	var user User
	getDB().First(&user, "username=?", username)
	type Result struct {
		Count uint
	}
	var result Result
	getDB().Raw("SELECT COUNT(*) FROM followers WHERE user_id=? AND follower_id=?", user.ID).Scan(&result)
	if result.Count == 0 {
		c.JSON(402, gin.H{"error": "Unathorized"})
		c.Abort()
	}
	getDB().Raw("DELETE FROM followers WHERE user_id=? AND follower_id=?", user.ID)
	// c.JSON(200, gin.H{
	// 	"profile": Profile{},
	// })
}

func getTags(c *gin.Context) {
	var tags []Tag
	getDB().Find(&tags)
	c.JSON(200, gin.H{
		"tags": tags,
	})
}

func authApi(c *gin.Context) *jwt.Token {
	tokenString := c.Request.Header.Get("Authorization")
	type MyCustomClaims struct {
		UserId uint64 `json:"user_id"`
		jwt.StandardClaims
	}
	bearerToken := strings.Split(tokenString, " ")

	token, err := jwt.ParseWithClaims(bearerToken[1], &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if token == nil || err == nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "not authorized",
			"error":   err.Error(),
		})
	}

	return token
}
