package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "localhost"
	port = "23234"
)

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		ssh.AllocatePty(),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := initialModel(s)
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

type screen int
type focus int

const (
	WelcomeScreen screen = iota
	LobbyScreen
	RoomScreen
)

const (
	None focus = iota
	InputFocus
	MessagesFocus
)

type Message struct{}

type model struct {
	width         int
	height        int
	currentScreen screen
	currentFocus  focus
	userInput     string
	messages      *[]Message
	session       ssh.Session
	style         lipgloss.Style
	errStyle      lipgloss.Style
}

func initialModel(s ssh.Session) model {
	renderer := bubbletea.MakeRenderer(s)
	return model{
		currentScreen: WelcomeScreen,
		currentFocus:  None,
		session:       s,
		style:         renderer.NewStyle().Foreground(lipgloss.Color("8")),
		errStyle:      renderer.NewStyle().Foreground(lipgloss.Color("3")),
	}

}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "c":
		case "j":
		case "s":
		case "u":
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	switch m.currentScreen {
	case WelcomeScreen:
		return welcomeView(m)
	case LobbyScreen:
		return lobbyView(m)
	}
	return ""
}
