package authentication

import "github.com/golang-jwt/jwt/v5"

type Option func(o *options)

type options struct {
	excludedMethods []string
	claims          interface{ jwt.Claims }
	mdKey           string
}

func WithExcludedMethods(methods ...string) Option {
	return func(o *options) {
		o.excludedMethods = methods
	}
}

func WithCustomClaims(claims interface{ jwt.Claims }) Option {
	return func(o *options) {
		o.claims = claims
	}
}

func WithCustomMetadataKey(mdKey string) Option {
	return func(o *options) {
		o.mdKey = mdKey
	}
}

func evaluateOptions(opts []Option) *options {
	o := &options{
		mdKey: "authorization",
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}
