package util

import (
	"net"
	"strings"
)

// IsRemoteClosedError whether the remote closed connection
func IsRemoteClosedError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "closed by the remote host") || strings.Contains(err.Error(), "an unreachable network"))
}

// IsClosedConnection whether the connection is already closed
func IsClosedConnection(err error) bool {
	return err != nil && strings.Contains(err.Error(), "use of closed network connection")
}

// IsTimeout whether is i/o timeout
func IsTimeout(err error) bool {
	if e, ok := err.(net.Error); ok && e != nil {
		return e.Timeout()
	}
	return false
	//return strings.Contains(err.Error(), "i/o timeout")
}
