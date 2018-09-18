package sentry

import (
	"github.com/getsentry/raven-go"
)

type options struct {
	serverName string
	release    string
	env        string
	tags       raven.Tags
	extra      raven.Extra
}

type option func(o *options)

func WithEnv(env string) option {
	return func(o *options) {
		o.env = env
	}
}

func WithTags(tags Tags) option {
	return func(o *options) {
		o.tags = raven.Tags(tags)
	}
}

func WithDefaultExtra(extra Extra) option {
	return func(o *options) {
		o.extra = raven.Extra(extra)
	}
}

func WithServerName(serverName string) option {
	return func(o *options) {
		o.serverName = serverName
	}
}

func WithRelease(release string) option {
	return func(o *options) {
		o.release = release
	}
}
