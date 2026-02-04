package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type FileInfo struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"is_dir"`
	Mod   string `json:"mod_time"`
}

func ListFiles(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = "E:\\"
	}

	// 确保路径存在且是目录
	f, err := os.Stat(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "路径不存在: " + err.Error()})
		return
	}
	if !f.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不是一个目录"})
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取目录: " + err.Error()})
		return
	}

	var files []FileInfo
	for _, entry := range entries {
		info, _ := entry.Info()
		files = append(files, FileInfo{
			Name:  entry.Name(),
			Size:  info.Size(),
			IsDir: entry.IsDir(),
			Mod:   info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, files)
}

func DeleteFile(c *gin.Context) {
	path := c.Query("path")
	if err := os.RemoveAll(path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
