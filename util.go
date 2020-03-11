package gosocket

import "strings"

// IsRemoteClosedError whether the remote closed connection
func IsRemoteClosedError(err error) bool {
	return strings.Contains(err.Error(), "closed by the remote host")
}

// IsClosedConnection whether the connection is already closed
func IsClosedConnection(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}

// IsTimeout whether is i/o timeout
func IsTimeout(err error) bool {
	return strings.Contains(err.Error(), "i/o timeout")
}
