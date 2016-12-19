package send

import (
	"github.com/tychoish/grip/level"
	"github.com/tychoish/grip/message"
)

// internalSender implements a Sender object that makes it possible to
// access logging messages, in the InternalMessage format without
// logging to an output method. The Send method does not filter out
// under-priority and unloggable messages. Used  for testing
// purposes.
type internalSender struct {
	name   string
	level  LevelInfo
	output chan *internalMessage
}

// InternalMessage provides a complete representation of all
// information associated with a logging event.
type internalMessage struct {
	Message  message.Composer
	Level    LevelInfo
	Logged   bool
	Priority level.Priority
	Rendered string
}

// NewInternalLogger creates and returns a Sender implementation that
// does not log messages, but converts them to the InternalMessage
// format and puts them into an internal channel, that allows you to
// access the massages via the extra "GetMessage" method. Useful for
// testing.
func NewInternalLogger(l LevelInfo) (*internalSender, error) {
	s := &internalSender{
		output: make(chan *internalMessage, 100),
	}

	if err := s.SetLevel(l); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *internalSender) Name() string     { return s.name }
func (s *internalSender) SetName(n string) { s.name = n }
func (s *internalSender) Close()           { close(s.output) }
func (s *internalSender) Type() SenderType { return Internal }
func (s *internalSender) Level() LevelInfo { return s.level }

func (s *internalSender) SetLevel(l LevelInfo) error {
	s.level = l
	return nil
}
func (s *internalSender) GetMessage() *internalMessage {
	return <-s.output
}

func (s *internalSender) Send(m message.Composer) {
	s.output <- &internalMessage{
		Message:  m,
		Priority: m.Priority(),
		Rendered: m.Resolve(),
		Logged:   s.level.ShouldLog(m),
	}
}
