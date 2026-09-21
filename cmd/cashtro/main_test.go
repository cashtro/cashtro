package main

import "testing"

func TestListenAddrPrefersAzurePort(t *testing.T) {
	t.Setenv("PORT", "8081")
	t.Setenv("WEBSITES_PORT", "9999")
	if got := listenAddr(":8080"); got != ":8081" {
		t.Fatalf("listenAddr = %q, want :8081", got)
	}
}

func TestListenAddrUsesWebsitesPort(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("WEBSITES_PORT", "80")
	if got := listenAddr(":8080"); got != ":80" {
		t.Fatalf("listenAddr = %q, want :80", got)
	}
}

func TestListenAddrFallsBackToFlag(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("WEBSITES_PORT", "")
	if got := listenAddr(":7070"); got != ":7070" {
		t.Fatalf("listenAddr = %q, want :7070", got)
	}
}
