package sentry

import (
	"github.com/getsentry/raven-go"
)

type options struct {
	serverName string
	release    string
	env        string
	logger     string
	tags       raven.Tags
	extra      raven.Extra
}

type Option func(o *options)

func WithEnv(env string) Option {
	return func(o *options) {
		o.env = env
	}
}

func WithTags(tags Tags) Option {
	return func(o *options) {
		o.tags = raven.Tags(tags)
	}
}

func WithDefaultExtra(extra Extra) Option {
	return func(o *options) {
		o.extra = raven.Extra(extra)
	}
}

func WithServerName(serverName string) Option {
	return func(o *options) {
		o.serverName = serverName
	}
}

func WithLogger(logger string) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func WithRelease(release string) Option {
	return func(o *options) {
		o.release = release
	}
}
