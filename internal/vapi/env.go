package vapi

import (
	"bufio"
	"os"
	"strings"
)

// LoadDotEnv reads .env then .env.local. Existing process env wins.
func LoadDotEnv(paths ...string) {
	if len(paths) == 0 {
		paths = []string{".env", ".env.local"}
	}
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			k, v, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}
			k = strings.TrimSpace(k)
			v = strings.TrimSpace(v)
			v = strings.Trim(v, `"'`)
			if k == "" || os.Getenv(k) != "" {
				continue
			}
			_ = os.Setenv(k, v)
		}
		_ = f.Close()
	}
}
