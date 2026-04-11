package httpclient

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

// WSConn represents a WebSocket connection
type WSConn struct {
	conn *websocket.Conn
}

// DialWS establishes a WebSocket connection
func DialWS(rawURL string, origin string, timeout int) (*WSConn, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	if u.Scheme != "ws" && u.Scheme != "wss" {
		return nil, fmt.Errorf("invalid WebSocket scheme: %s (use ws:// or wss://)", u.Scheme)
	}

	header := http.Header{}
	if origin != "" {
		header.Set("Origin", origin)
	}

	conn, resp, err := websocket.DefaultDialer.Dial(rawURL, header)
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("WebSocket handshake failed: %s %s", resp.Status, resp.Status)
		}
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &WSConn{conn: conn}, nil
}

// ReadMessage reads the next WebSocket message
// Returns a string that is either the text content or a binary size indicator
func (ws *WSConn) ReadMessage() (string, error) {
	msgType, reader, err := ws.conn.NextReader()
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	if msgType == websocket.BinaryMessage {
		return fmt.Sprintf("[binary data of length %d]", len(data)), nil
	}
	return string(data), nil
}

// WriteMessage sends a WebSocket text message
func (ws *WSConn) WriteMessage(msg string) error {
	return ws.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

// Close gracefully closes the WebSocket connection
func (ws *WSConn) Close() error {
	return ws.conn.Close()
}

// RunWebSocket runs the WebSocket interactive session
func RunWebSocket(rawURL string, origin string, timeout int) error {
	ws, err := DialWS(rawURL, origin, timeout)
	if err != nil {
		return err
	}
	defer ws.Close()

	// Channel for incoming messages
	errChan := make(chan error, 1)

	// Goroutine to read messages
	go func() {
		for {
			msg, err := ws.ReadMessage()
			if err != nil {
				if err == io.EOF || websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					errChan <- nil
					return
				}
				errChan <- err
				return
			}
			fmt.Printf("<<< %s\n", msg)
		}
	}()

	// Read from stdin and send
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		fmt.Printf(">>> %s\n", line)
		if err := ws.WriteMessage(line); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stdin error: %w", err)
	}

	// Close connection on EOF
	return ws.Close()
}

// IsWebSocketURL checks if the URL is a WebSocket URL
func IsWebSocketURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return u.Scheme == "ws" || u.Scheme == "wss"
}

// HasWebSocketScheme checks if URL starts with ws:// or wss://
func HasWebSocketScheme(rawURL string) bool {
	return strings.HasPrefix(rawURL, "ws://") || strings.HasPrefix(rawURL, "wss://")
}
