package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Note represents a note in the app
type Note struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var notes []Note
var noteID int

func main() {
	r := gin.Default()

	// Set up sessions
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// Routes
	r.POST("/signup", signup)
	r.POST("/login", login)

	// Protected routes
	auth := r.Group("/api")
	auth.Use(authMiddleware())
	{
		auth.POST("/notes", createNote)
		auth.GET("/notes", getNotes)
		auth.GET("/notes/:id", getNote)
		auth.PUT("/notes/:id", updateNote)
		auth.DELETE("/notes/:id", deleteNote)
	}

	r.Run(":8080")
}

// Middleware to check if the user is authenticated
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}

		// Extract the token from the "Authorization" header
		tokenString := authHeader[len("Bearer "):]

		// Verify the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Provide the same secret key used for token generation
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Token is valid, proceed to the next handler
		c.Next()
	}
}

// Sign up handler
func signup(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate username and password

	// Create user session
	session := sessions.Default(c)
	session.Set("user", username)
	session.Set("password", password)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Sign up successful"})
}

// Login handler

// Login handler
func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate username and password

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["password"] = password
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time (1 day)

	// Generate encoded token string
	tokenString, err := token.SignedString([]byte("your-secret-key")) // Replace with your secret key
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}

// Middleware to check if the user is authenticated

// Create a new note
func createNote(c *gin.Context) {
	var note Note
	if err := c.BindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note.ID = noteID
	noteID++
	notes = append(notes, note)

	c.JSON(http.StatusCreated, note)
}

// Get all notes
func getNotes(c *gin.Context) {
	c.JSON(http.StatusOK, notes)
}

// Get a specific note
func getNote(c *gin.Context) {
	id := c.Param("id")

	for _, note := range notes {
		if strconv.Itoa(note.ID) == id {
			c.JSON(http.StatusOK, note)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
}

// Update a note
func updateNote(c *gin.Context) {
	id := c.Param("id")

	var updatedNote Note
	if err := c.BindJSON(&updatedNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, note := range notes {
		if strconv.Itoa(note.ID) == id {
			notes[i].Title = updatedNote.Title
			notes[i].Body = updatedNote.Body
			c.JSON(http.StatusOK, notes[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
}

// Delete a note
func deleteNote(c *gin.Context) {
	id := c.Param("id")

	for i, note := range notes {
		if strconv.Itoa(note.ID) == id {
			notes = append(notes[:i], notes[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
}
