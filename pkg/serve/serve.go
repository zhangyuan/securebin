package serve

import (
	"fmt"
	"mime"
	"strings"
	"time"

	"securebin/ui"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CreateMessageRequest struct {
	ExpireAt       time.Time `json:"expire_at"`
	Content        string    `json:"content" binding:"required"`
	Password       string    `json:"password"`
	MaxAccessCount int64     `json:"max_access_count"`
}

type GetMessageRequest struct {
	Password string
}

type CreateMessageResponse struct {
	ID uint `json:"id"`
}

type GetMessageResponse struct {
	ExpireAt       time.Time `json:"expire_at"`
	Content        string    `json:"content"`
	AccessCount    int64     `json:"acccess_count"`
	MaxAccessCount int64     `json:"max_access_count"`
}

type Message struct {
	ExpireAt time.Time
	gorm.Model
	EncryptedPassword string
	Content           string
	MaxAccessCount    int64
	AccessCount       int64
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func BuildUIAssetPath(path string) string {
	return fmt.Sprintf("dist%s", path)
}

func Invoke(addr string) error {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&Message{}); err != nil {
		return err
	}

	r := gin.Default()

	staticFS := ui.StaticFS

	r.Use(func(ctx *gin.Context) {
		requestPath := ctx.Request.URL.Path

		if strings.HasPrefix(requestPath, "/api/") {
			ctx.Next()
			return
		}
		if ctx.Request.Method == "GET" {
			extentions := []string{".css", ".js", ".ico", ".png", ".jpg", ".svg"}
			for _, extension := range extentions {
				if strings.HasSuffix(requestPath, extension) {
					data, err := staticFS.ReadFile(BuildUIAssetPath(requestPath))
					if err != nil {
						_ = ctx.AbortWithError(500, err)
						return
					} else {
						ctx.Data(200, mime.TypeByExtension(extension), data)
						return
					}
				}
			}

			acceptHeader := ctx.Request.Header.Get("Accept")
			if strings.Contains(acceptHeader, "text/html") || strings.Contains(acceptHeader, "*/*") {
				file, err := staticFS.ReadFile(BuildUIAssetPath("/index.html"))
				if err != nil {
					_ = ctx.AbortWithError(500, err)
					return
				}
				ctx.Data(200, "text/html", file)
				return
			}

			return
		}
		ctx.AbortWithStatus(404)
	})

	r.POST("/api/messages", func(c *gin.Context) {
		var request CreateMessageRequest
		if err := c.BindJSON(&request); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}

		encryptedPassword := ""
		if request.Password != "" {
			encryptedPassword, err = HashPassword(request.Password)
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
				return
			}
		}

		message := Message{EncryptedPassword: encryptedPassword,
			Content: request.Content, MaxAccessCount: request.MaxAccessCount, ExpireAt: request.ExpireAt,
		}
		if err := db.Create(&message).Error; err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, CreateMessageResponse{
			ID: message.ID,
		})
	})

	r.POST("/api/messages/:id", func(c *gin.Context) {
		id := c.Param("id")

		var request GetMessageRequest
		if err := c.BindJSON(&request); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}

		var message Message
		if err := db.First(&message, id).Error; err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		if message.EncryptedPassword != "" {
			if !CheckPasswordHash(request.Password, message.EncryptedPassword) {
				c.AbortWithStatusJSON(400, gin.H{
					"error": "invalid password",
				})
				return
			}
		}

		if (message.ExpireAt != time.Time{}) && (message.ExpireAt.UnixMilli() >= time.Now().UnixMilli()) {
			if err := db.Delete(&Message{}, message.ID).Error; err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
				return
			}

			c.AbortWithStatusJSON(404, gin.H{"error": "NOT FOUND"})
			return
		}

		if err := db.Model(&message).Update("access_count", gorm.Expr("access_count + ?", 1)).Error; err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err})
			return
		}

		if err := db.First(&message, id).Error; err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		if (message.MaxAccessCount != 0) && (message.AccessCount >= message.MaxAccessCount) {
			if err := db.Delete(&Message{}, message.ID).Error; err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
				return
			}

			c.AbortWithStatusJSON(404, gin.H{"error": "NOT FOUND"})
		}

		c.JSON(200, GetMessageResponse{
			Content:        message.Content,
			ExpireAt:       message.ExpireAt,
			AccessCount:    message.AccessCount,
			MaxAccessCount: message.MaxAccessCount,
		})
	})

	return r.Run(addr)
}
