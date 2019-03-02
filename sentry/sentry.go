package sentry

import (
	raven "github.com/getsentry/raven-go"
)

const (
	_platform = "go"
)

var defaultOption = &options{
	logger: "root",
}

type Tags = raven.Tags
type Tag = raven.Tag
type Extra = raven.Extra

type User struct {
	ID     string   `json:"id,omitempty"`
	Email  string   `json:"email,omitempty"`
	Groups []string `json:"groups,omitempty"`
	IP     string   `json:"ip,omitempty"`
}

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

func SetOptions(option ...Option) {
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
			Logger:      defaultOption.logger,
			Fingerprint: make([]string, 0, 4),
		},
	}

	p.packet.Fingerprint[0] = defaultOption.logger
	if defaultOption.extra != nil {
		mergeExtra(p.packet.Extra, defaultOption.extra)
	}

	return p
}

func (p *Packet) RawPacket() *raven.Packet {
	return p.packet
}

func (p *Packet) SetCulprit(culprit string) {
	p.packet.Culprit = culprit
}

func (p *Packet) SetFingerprint(fingerprints ...string) {
	p.packet.Fingerprint = append(p.packet.Fingerprint, fingerprints...)
}

func (p *Packet) AddUser(user *User) {
	p.packet.Extra["user"] = user
}

func (p *Packet) AddExtra(extra Extra) {
	if extra != nil {
		mergeExtra(p.packet.Extra, extra)
	}
}

func (p *Packet) AddTag(key, value string) {
	p.packet.Tags = append(p.packet.Tags, Tag{key, value})
}

func (p *Packet) AddStackTrace(stack raven.Interface) {
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
