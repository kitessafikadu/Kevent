package main

import (
	"net/http"
	"time"
	"github.com/kitessafikadu/kevent/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"	
)

type registerRequest struct {
	Email		string `json:"email" binding:"required,email"`
	Password	string `json:"password" binding:"required,min=8"`
	Name		string `json:"name" binding:"required,min=3"`
}

type loginRequest struct{ 
	Email		string `json:"email" binding:"required,email"`
	Password	string `json:"password" binding:"required,min=8"`
}

type loginResponse struct {
	Token string `json:"token"`
	UserId int `json:"userId"`
}

func (app *application) login(c *gin.Context) {
	var auth loginRequest
	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := app.models.Users.GetByEmail(auth.Email)
	if existingUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": existingUser.Id,
		"expr":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(app.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: tokenString, UserId: existingUser.Id})

}

func (app *application) registerUser(c *gin.Context){
	var register registerRequest
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword,err:=bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create the user in the database
	register.Password = string(hashedPassword)
	user := database.User{
		Email:    register.Email,
		Password: register.Password,
		Name:     register.Name,
	}
	if err := app.models.Users.Insert(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
