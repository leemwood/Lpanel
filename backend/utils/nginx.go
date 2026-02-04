package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GenerateNginxConfig(domain string, port int, path string, proxyPass string) string {
	config := fmt.Sprintf(`server {
    listen %d;
    server_name %s;

    location / {
`, port, domain)

	if proxyPass != "" {
		config += fmt.Sprintf(`        proxy_pass %s;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
`, proxyPass)
	} else {
		config += fmt.Sprintf(`        root %s;
        index index.html index.htm;
`, path)
	}

	config += `    }
}
`
	return config
}

func SaveNginxConfig(name string, content string) error {
	confDir := "data/nginx/conf.d"
	if err := os.MkdirAll(confDir, 0755); err != nil {
		return err
	}
	confPath := filepath.Join(confDir, name+".conf")
	return os.WriteFile(confPath, []byte(content), 0644)
}

func RemoveNginxConfig(name string) error {
	confPath := filepath.Join("data/nginx/conf.d", name+".conf")
	return os.Remove(confPath)
}
