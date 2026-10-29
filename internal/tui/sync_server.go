package tui

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SyncServer struct {
	listener     net.Listener
	port         int
	clients      map[net.Conn]struct{}
	clientsMu    sync.Mutex
	running      bool
	currentSlide int
	// gotos carries the slides clients ask for, numbered from 1.
	gotos chan int
}

const port = 34622

func NewSyncServer() (*SyncServer, error) {
	return listenSync(fmt.Sprintf("127.0.0.1:%d", port))
}

func listenSync(address string) (*SyncServer, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port: %w", err)
	}

	server := &SyncServer{
		listener:     listener,
		port:         listener.Addr().(*net.TCPAddr).Port,
		clients:      make(map[net.Conn]struct{}),
		currentSlide: 0,
		gotos:        make(chan int),
	}

	return server, nil
}

func (s *SyncServer) GetPort() int {
	return s.port
}

// GoTos delivers the slides clients ask for, numbered from 1.
func (s *SyncServer) GoTos() <-chan int {
	return s.gotos
}

func (s *SyncServer) Start() {
	s.running = true

	go s.acceptConnections()
}

func (s *SyncServer) Stop() {
	s.running = false

	if s.listener != nil {
		s.listener.Close()
	}

	s.clientsMu.Lock()
	for client := range s.clients {
		client.Close()
	}
	s.clients = make(map[net.Conn]struct{})
	s.clientsMu.Unlock()

	slog.Info("Sync server stopped")
}

func (s *SyncServer) BroadcastSlideChange(slideNumber int) {
	s.currentSlide = slideNumber

	message := fmt.Sprintf("SLIDE:%d\n", slideNumber)

	s.clientsMu.Lock()
	for client := range s.clients {
		_, err := client.Write([]byte(message))
		if err != nil {
			delete(s.clients, client)
			client.Close()
		}
	}
	s.clientsMu.Unlock()
}

func (s *SyncServer) acceptConnections() error {
	for {
		running := s.running

		if !running {
			break
		}

		conn, err := s.listener.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %w", err)
		}

		s.clientsMu.Lock()
		s.clients[conn] = struct{}{}
		s.clientsMu.Unlock()

		currentSlide := s.currentSlide

		message := fmt.Sprintf("SLIDE:%d\n", currentSlide)
		_, err = conn.Write([]byte(message))
		if err != nil {
			return fmt.Errorf("failed to send current slide to new client: %w", err)
		}

		go s.serve(conn)
	}
	return nil
}

// serve reads what a client asks for until it leaves: "GOTO:5" for the fifth
// slide, as `kyma goto 5` sends.
func (s *SyncServer) serve(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if slide, ok := strings.CutPrefix(line, "GOTO:"); ok {
			if slideNum, err := strconv.Atoi(slide); err == nil {
				s.gotos <- slideNum
			}
		}
	}

	s.clientsMu.Lock()
	delete(s.clients, conn)
	s.clientsMu.Unlock()
	conn.Close()
}

type SyncClient struct {
	conn net.Conn
	port int
}

func NewSyncClient() (*SyncClient, error) {
	return dialSync(fmt.Sprintf("127.0.0.1:%d", port))
}

func dialSync(address string) (*SyncClient, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sync server: %w", err)
	}

	client := &SyncClient{
		conn: conn,
		port: port,
	}

	return client, nil
}

func (c *SyncClient) ListenForSlideChanges(slideChangeChan chan<- int) {
	scanner := bufio.NewScanner(c.conn)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "SLIDE:") {
			slideNumStr := strings.TrimPrefix(line, "SLIDE:")
			if slideNum, err := strconv.Atoi(slideNumStr); err == nil {
				slideChangeChan <- slideNum
			}
		}
	}
}

// GoTo asks the presentation to show a slide, numbered from 1.
func (c *SyncClient) GoTo(slideNum int) error {
	if err := c.conn.SetDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return err
	}

	// The server first sends the slide it is on. Closing the connection with
	// that still unread would reset it, which can drop the request.
	if _, err := bufio.NewReader(c.conn).ReadString('\n'); err != nil {
		return fmt.Errorf("failed to read from sync server: %w", err)
	}

	if _, err := fmt.Fprintf(c.conn, "GOTO:%d\n", slideNum); err != nil {
		return fmt.Errorf("failed to send slide to sync server: %w", err)
	}
	return nil
}

func (c *SyncClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
