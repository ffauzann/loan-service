package logging

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/interceptor/grpc/unary/ctxtags"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/util"
	"github.com/ffauzann/grpc-postgres-auth-user-asymmetric/pkg/common/util/sanitize"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type Option func(o *options)

type options struct {
	parseError bool
	mapError   map[error]codes.Code
	ctxtags    bool
}

// WithErrorParser parse unknown server error with unspecified internal error.
func WithErrorParser(m map[error]codes.Code) Option {
	return func(o *options) {
		o.mapError = m
		o.parseError = true
	}
}

// WithCtxTag parse unknown server error with unspecified internal error.
func WithCtxTag(val bool) Option {
	return func(o *options) {
		o.ctxtags = val
	}
}

func evaluateOptions(opts []Option) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *options) getError(err error) (codes.Code, error, bool) {
	if code, ok := o.mapError[err]; ok {
		return code, err, ok
	}

	if ok := strings.Contains(err.Error(), "VALIDATION_ERR: "); ok {
		return codes.InvalidArgument, errors.New(strings.ReplaceAll(err.Error(), "VALIDATION_ERR: ", "")), ok
	}

	if s := strings.SplitN(err.Error(), ":", 3); s[0] == "DYNAMIC_ERR" && len(s) > 2 {
		code, _ := strconv.Atoi(s[1])
		if code == 0 {
			return codes.Unknown, errors.New(strings.Join(s[1:], ":")), true
		}
		return codes.Code(code), errors.New(s[2]), true
	}

	return codes.Internal, errors.New(codes.Internal.String()), false
}

func (o *options) prepareLog(ctx context.Context, l *zap.Logger, resp interface{}) *zap.Logger {
	if o.ctxtags {
		tags := ctxtags.Extract(ctx)
		values := tags.Values()
		cID, _ := values[ctxtags.CIDKey].(string)
		md := sanitize.Sanitize(values[ctxtags.MDKey])
		req := sanitize.Sanitize(values[ctxtags.ReqKey])
		resp := sanitize.Sanitize(resp)
		return l.With(
			zap.String(ctxtags.CIDKey, cID),
			zap.Any(ctxtags.MDKey, util.StructToMap(md)),
			zap.Any(ctxtags.ReqKey, util.StructToMap(req)),
			zap.Any(ctxtags.RespKey, util.StructToMap(resp)),
		)
	}

	return l
}
