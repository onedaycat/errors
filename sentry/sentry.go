package sentry

import (
	raven "github.com/getsentry/raven-go"
)

const (
	_platform = "go"
)

var defaultOption = &options{}

type Tags = raven.Tags
type Extra = raven.Extra
type User = raven.User

type Packet struct {
	packet *raven.Packet
}

type Sentry struct {
	dsn     string
	options *options
}

func SetDSN(dsn string) {
	raven.SetDSN(dsn)
}

func SetOptions(option ...option) {
	for _, opt := range option {
		opt(defaultOption)
	}
}

func NewPacket(err error) *Packet {
	p := &Packet{
		packet: &raven.Packet{
			Message:     err.Error(),
			Level:       raven.ERROR,
			Platform:    _platform,
			Extra:       Extra{},
			Tags:        defaultOption.tags,
			Environment: defaultOption.env,
			ServerName:  defaultOption.serverName,
			Release:     defaultOption.release,
		},
	}

	if defaultOption.extra != nil {
		mergeExtra(p.packet.Extra, defaultOption.extra)
	}

	return p
}

func (p *Packet) AddUser(user *User) {
	p.packet.Extra["user"] = user
}

func (p *Packet) AddExtra(extra Extra) {
	if extra != nil {
		mergeExtra(p.packet.Extra, extra)
	}
}

func (p *Packet) AddStackTrace(stack *raven.Stacktrace) {
	if stack != nil {
		p.packet.Interfaces = append(p.packet.Interfaces, stack)
	}
}

func Capture(packet *Packet) {
	_, _ = raven.Capture(packet.packet, nil)
}

func CaptureAndWait(packet *Packet) {
	eventID, ch := raven.Capture(packet.packet, nil)
	if eventID != "" {
		<-ch
	}
}

func mergeExtra(baseExtra Extra, extra Extra) {
	for key, val := range extra {
		baseExtra[key] = val
	}
}
