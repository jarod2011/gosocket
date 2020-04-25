package client

import (
	"github.com/jarod2011/gosocket/buffers"
	"github.com/jarod2011/gosocket/conn"
	"github.com/jarod2011/gosocket/util"
	"io"
	"time"
)

type OnReadHandle func(reader io.Reader) []byte

var bytesPool = buffers.New()

func NewConnectionHandler(read OnReadHandle, closed OnClosed) ConnectionHandler {
	return func(cc conn.Conn, opt Options) {
		isClosed := false
		var cErr error
		for {
			cErr = nil
			cc.SetReadDeadline(time.Now().Add(opt.ReadTimeout))
			buf := bytesPool.Get()
			n, err := buf.ReadFrom(cc)
			if n > 0 {
				w := read(buf)
				if len(w) > 0 {
					cc.SetWriteDeadline(time.Now().Add(opt.WriteTimeout))
					_, err := cc.Write(w)
					if err != nil {
						cErr = err
						if util.IsClosedConnection(err) || util.IsRemoteClosedError(err) {
							isClosed = true
							goto clean
						}
						if util.IsTimeout(err) {
							goto clean
						}
					}
				}
			}
			if err != nil {
				cErr = err
				if util.IsClosedConnection(err) || util.IsRemoteClosedError(err) {
					isClosed = true
					goto clean
				}
				if util.IsTimeout(err) {
					continue
				}
				goto clean
			}
		}
	clean:
		closed(cc, isClosed, cErr)
	}
}
