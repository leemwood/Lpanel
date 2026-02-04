package handlers

import (
	"Lpanel/backend/config"
	"Lpanel/backend/models"
	"Lpanel/backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListSites(c *gin.Context) {
	var sites []models.Site
	config.DB.Find(&sites)
	c.JSON(http.StatusOK, sites)
}

func CreateSite(c *gin.Context) {
	var site models.Site
	if err := c.ShouldBindJSON(&site); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	site.Status = "stopped"
	site.NginxConf = utils.GenerateNginxConfig(site.Domain, site.Port, site.Path, site.ProxyPass)

	if err := config.DB.Create(&site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Refresh runner
	var allSites []models.Site
	config.DB.Find(&allSites)
	utils.Runner.UpdateSites(allSites)

	// Save Nginx config file
	if err := utils.SaveNginxConfig(site.Name, site.NginxConf); err != nil {
		// Log error but continue
	}

	c.JSON(http.StatusOK, site)
}

func DeleteSite(c *gin.Context) {
	id := c.Param("id")
	var site models.Site
	if err := config.DB.First(&site, id).Error; err == nil {
		utils.RemoveNginxConfig(site.Name)
	}

	if err := config.DB.Delete(&models.Site{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Refresh runner
	var allSites []models.Site
	config.DB.Find(&allSites)
	utils.Runner.UpdateSites(allSites)

	c.JSON(http.StatusOK, gin.H{"message": "Site deleted"})
}

func ToggleSite(c *gin.Context) {
	id := c.Param("id")
	var site models.Site
	if err := config.DB.First(&site, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	if site.Status == "running" {
		site.Status = "stopped"
	} else {
		site.Status = "running"
	}

	config.DB.Save(&site)

	// Refresh runner
	var allSites []models.Site
	config.DB.Find(&allSites)
	utils.Runner.UpdateSites(allSites)

	c.JSON(http.StatusOK, site)
}
