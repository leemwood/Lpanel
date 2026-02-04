package utils

import (
	"Lpanel/backend/models"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type SiteRunner struct {
	servers map[int]*http.Server
	mu      sync.Mutex
	sites   []models.Site
}

var Runner = &SiteRunner{
	servers: make(map[int]*http.Server),
}

func (sr *SiteRunner) UpdateSites(sites []models.Site) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.sites = sites
	sr.refreshListeners()
}

func (sr *SiteRunner) refreshListeners() {
	// Get all unique ports for running sites
	ports := make(map[int]bool)
	for _, site := range sr.sites {
		if site.Status == "running" {
			ports[site.Port] = true
		}
	}

	// Stop servers for ports no longer needed
	for port, server := range sr.servers {
		if !ports[port] {
			log.Printf("Stopping server on port %d", port)
			server.Close()
			delete(sr.servers, port)
		}
	}

	// Start servers for new ports
	for port := range ports {
		if _, exists := sr.servers[port]; !exists {
			sr.startServer(port)
		}
	}
}

func (sr *SiteRunner) startServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sr.mu.Lock()
		defer sr.mu.Unlock()

		host := r.Host
		// Remove port from host if present
		if h, _, err := netSplitHostPort(host); err == nil {
			host = h
		}

		var targetSite *models.Site
		for i := range sr.sites {
			s := &sr.sites[i]
			if s.Port == port && s.Status == "running" && (s.Domain == host || s.Domain == "0.0.0.0" || s.Domain == "*") {
				targetSite = s
				break
			}
		}

		// Fallback to first running site on this port if no domain match
		if targetSite == nil {
			for i := range sr.sites {
				s := &sr.sites[i]
				if s.Port == port && s.Status == "running" {
					targetSite = s
					break
				}
			}
		}

		if targetSite == nil {
			http.Error(w, "Site not found or not running", http.StatusNotFound)
			return
		}

		if targetSite.ProxyPass != "" {
			target, _ := url.Parse(targetSite.ProxyPass)
			proxy := httputil.NewSingleHostReverseProxy(target)
			proxy.ServeHTTP(w, r)
		} else {
			http.FileServer(http.Dir(targetSite.Path)).ServeHTTP(w, r)
		}
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	sr.servers[port] = server

	go func() {
		log.Printf("Starting site listener on port %d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server on port %d failed: %v", port, err)
		}
	}()
}

// Helper because net.SplitHostPort fails if no port is present
func netSplitHostPort(host string) (string, string, error) {
	for i := len(host) - 1; i >= 0; i-- {
		if host[i] == ':' {
			return host[:i], host[i+1:], nil
		}
	}
	return host, "", fmt.Errorf("no port found")
}
