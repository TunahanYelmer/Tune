package daemon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/tunahanyelmer/Tune/internal/providers"
)

// Serve starts the daemon loop, owning the given provider exclusively
// for the lifetime of the process. Blocks until the listener fails.
func Serve(provider providers.Provider) error {
	path := SocketPath()
	os.Remove(path) // clean up a stale socket from a previous crashed run

	ln, err := net.Listen("unix", path)
	if err != nil {
		return fmt.Errorf("listening on socket: %w", err)
	}
	defer ln.Close()
	defer os.Remove(path)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn, provider)
	}
}

func handleConn(conn net.Conn, provider providers.Provider) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return
	}

	var req Request
	if err := json.Unmarshal(line, &req); err != nil {
		writeResponse(conn, Response{OK: false, Error: "invalid request"})
		return
	}

	writeResponse(conn, handle(provider, req))
}

func handle(p providers.Provider, req Request) Response {
	switch req.Action {
	case ActionLogin:
		if err := p.Login(); err != nil {
			return errResp(err)
		}
		return Response{OK: true}

	case ActionPlay:
		if err := p.Play(req.Query); err != nil {
			return errResp(err)
		}
		return Response{OK: true}

	case ActionPause:
		if err := p.Pause(); err != nil {
			return errResp(err)
		}
		return Response{OK: true}

	case ActionNext:
		if err := p.Next(); err != nil {
			return errResp(err)
		}
		return Response{OK: true}

	case ActionPrev:
		if err := p.Previous(); err != nil {
			return errResp(err)
		}
		return Response{OK: true}

	case ActionCurrent:
		track, err := p.Current()
		if err != nil {
			return errResp(err)
		}
		return Response{OK: true, Track: track}

	case ActionSearch:
		tracks, err := p.Search(req.Query)
		if err != nil {
			return errResp(err)
		}
		return Response{OK: true, Tracks: tracks}

	case ActionVolume:
		if err := p.SetVolume(req.Level); err != nil {
			return errResp(err)
		}
		return Response{OK: true}

	default:
		return Response{OK: false, Error: fmt.Sprintf("unknown action: %s", req.Action)}
	}
}

func errResp(err error) Response {
	return Response{OK: false, Error: err.Error()}
}

func writeResponse(conn net.Conn, resp Response) {
	data, _ := json.Marshal(resp)
	conn.Write(append(data, '\n'))
}