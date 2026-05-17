package network

import (
	"context"
	"net"
	"time"
)

type ReachabilityResult struct {
	OK  bool
	Err error
}

func Reachable(ctx context.Context, address string, timeout time.Duration) ReachabilityResult {
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return ReachabilityResult{Err: err}
	}
	_ = conn.Close()
	return ReachabilityResult{OK: true}
}
