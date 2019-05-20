package sentry

import (
    "strconv"

    "github.com/getsentry/raven-go"
    "github.com/onedaycat/errors"
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

const defaultFingerprint = "{{ default }}"
const inputsKey = "inputs"

type User struct {
    ID     string   `json:"id,omitempty"`
    Email  string   `json:"email,omitempty"`
    Groups []string `json:"groups,omitempty"`
    IP     string   `json:"ip,omitempty"`
}

type stackInput struct {
    Exception string      `json:"exception"`
    Input     interface{} `json:"input"`
}

func (h *User) Class() string { return "user" }

type Packet struct {
    packet *raven.Packet
}

type Sentry struct {
    dsn     string
    options *options
}

func SetDSN(dsn string) {
    _ = raven.SetDSN(dsn)
}

func SetOptions(option ...Option) {
    for _, opt := range option {
        opt(defaultOption)
    }
}

func NewPacket(err errors.Error) *Packet {
    if err == nil {
        return nil
    }

    isPanic := err.IsPanic()
    msg := err.GetMessage()
    root := err.RootError()

    expsList := make([]*raven.Exception, 0, 5)
    inputs := make([]interface{}, 0, 5)

    var errStack *errors.Stacktrace
    for err != nil {
        errStack = err.GetStacktrace()
        if errStack == nil {
            expsList = append(expsList, &raven.Exception{
                Value: err.Error(),
                Type:  err.GetCode(),
            })
        } else {
            stackTrace := &raven.Stacktrace{
                Frames: make([]*raven.StacktraceFrame, len(errStack.Frames)),
            }
            for i, frame := range errStack.Frames {
                stackTrace.Frames[i] = raven.NewStacktraceFrame(0, frame.Function, frame.Filename, frame.Lineno, 3, nil)
            }

            expsList = append(expsList, &raven.Exception{
                Value:      err.Error(),
                Type:       err.GetCode(),
                Stacktrace: stackTrace,
            })
        }

        inputs = append(inputs, err.GetInput())

        xerr := err.Unwrap()
        if xerr == nil {
            break
        }
        err = xerr.(errors.Error)
    }

    extra := Extra{}
    n := len(expsList)
    exps := &raven.Exceptions{
        Values: make([]*raven.Exception, n),
    }

    stackInputs := make([]*stackInput, n)
    for i := 0; i < n; i++ {
        exps.Values[i] = expsList[n-i-1]
        stackInputs[i] = &stackInput{
            Exception: expsList[i].Type + "_" + strconv.Itoa(i+1),
            Input:     inputs[i],
        }
    }

    extra[inputsKey] = stackInputs

    p := &Packet{
        packet: &raven.Packet{
            Message:     msg,
            Platform:    _platform,
            Extra:       extra,
            Tags:        defaultOption.tags,
            Environment: defaultOption.env,
            ServerName:  defaultOption.serverName,
            Release:     defaultOption.release,
            Logger:      defaultOption.logger,
            Culprit:     root.GetMessage(),
            Fingerprint: []string{
                defaultFingerprint,
                defaultOption.logger,
            },
            Interfaces: []raven.Interface{exps},
        },
    }

    if isPanic {
        p.packet.Level = raven.FATAL
    } else {
        p.packet.Level = raven.ERROR
    }

    if defaultOption.extra != nil {
        mergeExtra(p.packet.Extra, defaultOption.extra)
    }

    return p
}

func (p *Packet) RawPacket() *raven.Packet {
    return p.packet
}

func (p *Packet) SetMessage(msg string) {
    p.packet.Message = msg
}

func (p *Packet) SetCulprit(culprit string) {
    p.packet.Culprit = culprit
}

func (p *Packet) AddFingerprint(fingerprints ...string) {
    p.packet.Fingerprint = append(p.packet.Fingerprint, fingerprints...)
}

func (p *Packet) SetUser(user *User) {
    p.packet.Interfaces = append(p.packet.Interfaces, user)
}

func (p *Packet) AddExtra(extra Extra) {
    if extra != nil {
        mergeExtra(p.packet.Extra, extra)
    }
}

func (p *Packet) AddTag(key, value string) {
    p.packet.Tags = append(p.packet.Tags, Tag{Key: key, Value: value})
}

func Capture(packet *Packet) {
    if packet == nil {
        return
    }
    _, _ = raven.Capture(packet.packet, nil)
}

func CaptureAndWait(packet *Packet) {
    if packet == nil {
        return
    }
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
