package logging

import (
	"os"

	"github.com/tychoish/grip/level"
	"github.com/tychoish/grip/message"
	"github.com/tychoish/grip/send"
)

// Grip provides the core implementation of the Logging interface. The
// interface is mirrored in the "grip" package's public interface, to
// provide a single, global logging interface that requires minimal
// configuration.
type Grip struct {
	sender send.Sender
}

// NewGrip takes the name for a logging instance and creates a new
// Grip instance with configured with a Bootstrap logging
// instance. The default level is "Notice" and the threshold level is
// "info."
func NewGrip(name string) *Grip {
	return &Grip{send.NewBootstrapLogger(name,
		send.LevelInfo{
			Threshold: level.Info,
			Default:   level.Notice,
		}),
	}
}

// Name of the logger instance
func (g *Grip) Name() string {
	return g.sender.Name()
}

// SetName declare a name string for the logger, including in the logging
// message. Typically this is included on the output of the command.
func (g *Grip) SetName(name string) {
	g.sender.SetName(name)
}

// Internal

// For sending logging messages, in most cases, use the
// Journaler.sender.Send() method, but we have a couple of methods to
// use for the Panic/Fatal helpers.
func (g *Grip) sendPanic(m message.Composer) {
	// the Send method in the Sender interface will perform this
	// check but to add fatal methods we need to do this here.
	if g.sender.Level().ShouldLog(m) {
		return
	}

	g.sender.Send(m)
	panic(m.Resolve())
}

func (g *Grip) sendFatal(m message.Composer) {
	// the Send method in the Sender interface will perform this
	// check but to add fatal methods we need to do this here.
	if g.sender.Level().ShouldLog(m) {
		return
	}

	g.sender.Send(m)
	os.Exit(1)
}
