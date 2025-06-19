package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	Email		string `json:"email" binding:"required,email"`
	Password	string `json:"password" binding:"required,min=8"`
	Name		string `json:"name" binding:"required,min=3"`
}


