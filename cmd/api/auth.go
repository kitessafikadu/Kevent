package main

import (
	"net/http"
	"github.com/kitessafikadu/kevent/internal/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"	
)

type registerRequest struct {
	Email		string `json:"email" binding:"required,email"`
	Password	string `json:"password" binding:"required,min=8"`
	Name		string `json:"name" binding:"required,min=3"`
}

type loginRequest struct{
	
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
