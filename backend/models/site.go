package models

import (
	"gorm.io/gorm"
)

type Site struct {
	gorm.Model
	Name       string `json:"name" gorm:"uniqueIndex"`
	Domain     string `json:"domain" gorm:"uniqueIndex"`
	Path       string `json:"path"`
	Port       int    `json:"port"`
	ProxyPass  string `json:"proxy_pass"` // For reverse proxy
	Status     string `json:"status"`     // running, stopped
	NginxConf  string `json:"nginx_conf"`
}
