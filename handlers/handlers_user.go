package handlers

import (
	"github.com/gin-gonic/gin"
	"golang-cookies/internal/config"
)

type LocalApiConfig struct {
	*config.ApiConfig
}

func (lac *LocalApiConfig) handler_CreateUser(c *gin.Context) {
}
