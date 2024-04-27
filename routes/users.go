package routes

import (
	"guardian/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsersRoute(c *gin.Context) {
	users, err := data.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func UpdateRightsRoute(c *gin.Context) {
	type Payload struct {
		Username  string `json:"username" binding:"required"`
		NewRights int    `json:"newRights" binding:"required"`
	}
	var payload Payload
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := data.UpdateUserRights(payload.Username, payload.NewRights)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user rights updated"})
}
