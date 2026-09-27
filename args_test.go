package main

import (
	"os"
	"testing"
)

func TestParsePluginOptions(t *testing.T) {
	s := "proto=rudp;type=server;encrypt=aes-128;compress=1"
	opts, err := parsePluginOptions(s)
	if err != nil {
		t.Fatalf("parsePluginOptions failed: %v", err)
	}
	if v, ok := opts.Get("proto"); !ok || v != "rudp" {
		t.Fatalf("expected proto=rudp, got %v", v)
	}
	if v, ok := opts.Get("type"); !ok || v != "server" {
		t.Fatalf("expected type=server, got %v", v)
	}
	if v, ok := opts.Get("encrypt"); !ok || v != "aes-128" {
		t.Fatalf("expected encrypt=aes-128, got %v", v)
	}
	if v, ok := opts.Get("compress"); !ok || v != "1" {
		t.Fatalf("expected compress=1, got %v", v)
	}
}

func TestParseEnv(t *testing.T) {
	os.Setenv("SS_REMOTE_HOST", "127.0.0.1")
	os.Setenv("SS_REMOTE_PORT", "8388")
	os.Setenv("SS_LOCAL_HOST", "127.0.0.1")
	os.Setenv("SS_LOCAL_PORT", "1080")
	os.Setenv("SS_PLUGIN_OPTIONS", "proto=kcp;__android_vpn=1")
	defer func() {
		os.Unsetenv("SS_REMOTE_HOST")
		os.Unsetenv("SS_REMOTE_PORT")
		os.Unsetenv("SS_LOCAL_HOST")
		os.Unsetenv("SS_LOCAL_PORT")
		os.Unsetenv("SS_PLUGIN_OPTIONS")
	}()

	opts, err := parseEnv()
	if err != nil {
		t.Fatalf("parseEnv failed: %v", err)
	}
	if v, ok := opts.Get("remoteaddr"); !ok || v != "127.0.0.1:8388" {
		t.Fatalf("expected remoteaddr=127.0.0.1:8388, got %v", v)
	}
	if v, ok := opts.Get("proto"); !ok || v != "kcp" {
		t.Fatalf("expected proto=kcp, got %v", v)
	}
	if _, ok := opts.Get("__android_vpn"); !ok {
		t.Fatalf("expected __android_vpn to be present")
	}
}

func TestIPv6(t *testing.T) {
	if !isIPv6("::1") {
		t.Fatalf("expected ::1 to be IPv6")
	}
	if isIPv6("127.0.0.1") {
		t.Fatalf("expected 127.0.0.1 not to be IPv6")
	}
}
