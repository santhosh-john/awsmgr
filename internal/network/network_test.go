package network

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestReachableSucceedsForListeningAddress(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen() error = %v", err)
	}
	defer listener.Close()

	result := Reachable(context.Background(), listener.Addr().String(), time.Second)
	if !result.OK {
		t.Fatalf("Reachable().OK = false, want true; result=%#v", result)
	}
	if result.Err != nil {
		t.Fatalf("Reachable().Err = %v, want nil", result.Err)
	}
}

func TestReachableFailsForClosedAddress(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen() error = %v", err)
	}
	address := listener.Addr().String()
	listener.Close()

	result := Reachable(context.Background(), address, 100*time.Millisecond)
	if result.OK {
		t.Fatal("Reachable().OK = true, want false")
	}
	if result.Err == nil {
		t.Fatal("Reachable().Err = nil, want error")
	}
}
