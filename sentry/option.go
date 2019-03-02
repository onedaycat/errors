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

func WithServiceName(serviceName string) Option {
	return func(o *options) {
		o.logger = serviceName
	}
}

func WithVersion(version string) Option {
	return func(o *options) {
		if o.tags == nil {
			o.tags = Tags{
				{"version", version},
			}
		} else {
			o.tags = append(o.tags, Tag{"version", version})
		}
	}
}

func WithRelease(release string) Option {
	return func(o *options) {
		o.release = release
	}
}
