package al_plugin

import (
	"errors"
	"io"
	"net"
	"os"
	"testing"
	"time"
)

func TestIOConnCloseUnblocksRead(t *testing.T) {
	reader, writer := io.Pipe()
	conn, _ := NewIOConn(reader, writer, NewIOAddr(reader, writer))
	done := make(chan error, 1)
	go func() { _, err := conn.Read(make([]byte, 1)); done <- err }()
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("read succeeded")
		}
	case <-time.After(time.Second):
		t.Fatal("read blocked after close")
	}
}

func TestIOConnDeadline(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	conn, _ := NewIOConn(reader, writer, NewIOAddr(reader, writer))
	defer conn.Close()
	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	_, err = conn.Read(make([]byte, 1))
	if !errors.Is(err, os.ErrDeadlineExceeded) {
		t.Fatalf("deadline error %v", err)
	}
	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Write([]byte("x")); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Read(make([]byte, 1)); err != nil {
		t.Fatal(err)
	}
}

func TestIOListenerCloseBeforeAcceptAndRepeatedClose(t *testing.T) {
	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()
	listener, _ := NewIOListener(reader, writer)
	listener.Close()
	listener.Close()
	if _, err := listener.Accept(); !errors.Is(err, net.ErrClosed) {
		t.Fatalf("accept error %v", err)
	}
}
