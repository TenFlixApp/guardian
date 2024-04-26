package routes

import (
	"guardian/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDashboardStatsRoute(c *gin.Context) {
	adminUserCount, err := data.CountAdminUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	disabledUserCount, err := data.CountDisabledUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"adminUserCount":    adminUserCount,
		"disabledUserCount": disabledUserCount,
	})
}
