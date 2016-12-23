// +build linux

package send

import (
	"log"
	"os"
	"strings"

	"github.com/coreos/go-systemd/journal"
	"github.com/tychoish/grip/level"
	"github.com/tychoish/grip/message"
)

type systemdJournal struct {
	options  map[string]string
	fallback *log.Logger
	*base
}

// NewJournaldLogger creates a Sender object that writes log messages
// to the system's systemd journald logging facility. If there's an
// error with the sending to the journald, messages fallback to
// writing to standard output.
func NewJournaldLogger(name string, l LevelInfo) (Sender, error) {
	s := &systemdJournal{
		options: make(map[string]string),
		base:    newBase(name),
	}

	s.reset = func() {
		s.fallback = log.New(os.Stdout, strings.Join([]string{"[", s.Name(), "] "}, ""), log.LstdFlags)
	}

	if err := s.SetLevel(l); err != nil {
		return nil, err
	}

	s.reset()

	return s, nil
}

func (s *systemdJournal) Close() error     { return nil }
func (s *systemdJournal) Type() SenderType { return Systemd }

func (s *systemdJournal) Send(m message.Composer) {
	if s.level.ShouldLog(m) {
		msg := m.Resolve()
		p := m.Priority()
		err := journal.Send(msg, s.Level().convertPrioritySystemd(p), s.options)
		if err != nil {
			s.fallback.Println("systemd journaling error:", err.Error())
			s.fallback.Printf("[p=%s]: %s", p, msg)
		}
	}
}

func (l LevelInfo) convertPrioritySystemd(p level.Priority) journal.Priority {
	switch p {
	case level.Emergency:
		return journal.PriEmerg
	case level.Alert:
		return journal.PriAlert
	case level.Critical:
		return journal.PriCrit
	case level.Error:
		return journal.PriErr
	case level.Warning:
		return journal.PriWarning
	case level.Notice:
		return journal.PriNotice
	case level.Info:
		return journal.PriInfo
	case level.Debug:
		return journal.PriDebug
	default:
		return l.convertPrioritySystemd(l.Default)
	}
}
