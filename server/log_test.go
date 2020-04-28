package server

import (
	"bytes"
	log "github.com/jarod2011/go-log"
	"sync"
	"testing"
	"time"
)

func TestDefaultLogger_EnableDebug(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetLevel(log.Info)
	log.D().SetWriter(buf)
	l := new(defaultLogger)
	l.Log(debugPrefix, "123")
	if buf.Len() > 0 {
		t.Error(buf.Len())
	}
	l.EnableDebug()
	l.Logf(debugPrefix+"%v", time.Now().String())
	if buf.Len() == 0 {
		t.Error(buf.Len())
	}
	buf.Reset()
	l.Log(infoPrefix, "11")
	if buf.Len() > 0 {
		t.Error(buf.Len())
	}
	log.I().SetWriter(buf)
	l.Log(infoPrefix, "111")
	if buf.Len() == 0 {
		t.Error(buf.Len())
	}
	buf.Reset()
	l.Log(warnPrefix, "11")
	if buf.Len() > 0 {
		t.Error(buf.Len())
	}
	log.W().SetWriter(buf)
	l.Log(warnPrefix, "111")
	if buf.Len() == 0 {
		t.Error(buf.Len())
	}
	buf.Reset()
	l.Log(errPrefix, "11")
	if buf.Len() > 0 {
		t.Error(buf.Len())
	}
	log.E().SetWriter(buf)
	l.Log(errPrefix, "111")
	if buf.Len() == 0 {
		t.Error(buf.Len())
	}
	buf.Reset()
	var wg sync.WaitGroup
	wg.Add(1)
	log.SetFatalExit(false)
	go func() {
		defer wg.Done()
		defer func() {
			recover()
		}()
		l.Log(fatalPrefix, "11")
		if buf.Len() > 0 {
			t.Error(buf.Len())
		}
		log.E().SetWriter(buf)
		l.Log(fatalPrefix, "111")
		if buf.Len() == 0 {
			t.Error(buf.Len())
		}
		buf.Reset()
	}()
	wg.Wait()
	l.Logf("12334%d", 11)
	l.Log()
}
