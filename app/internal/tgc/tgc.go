package tgc

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-faster/errors"
	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	tdclock "github.com/gotd/td/clock"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
	"golang.org/x/time/rate"

	"github.com/iyear/tdl/pkg/clock"
	"github.com/iyear/tdl/pkg/consts"
	"github.com/iyear/tdl/pkg/key"
	"github.com/iyear/tdl/pkg/kv"
	"github.com/iyear/tdl/pkg/logger"
	"github.com/iyear/tdl/pkg/recovery"
	"github.com/iyear/tdl/pkg/retry"
	"github.com/iyear/tdl/pkg/storage"
	"github.com/iyear/tdl/pkg/utils"
)

func NewDefaultMiddlewares(ctx context.Context) ([]telegram.Middleware, error) {
	_clock, err := Clock()
	if err != nil {
		return nil, errors.Wrap(err, "create clock")
	}

	return []telegram.Middleware{
		recovery.New(ctx, Backoff(_clock)),
		retry.New(5),
		floodwait.NewSimpleWaiter(),
	}, nil
}

func New(ctx context.Context, login bool, middlewares ...telegram.Middleware) (*telegram.Client, kv.KV, error) {
	kvd, err := kv.From(ctx).Open(viper.GetString(consts.FlagNamespace))
	if err != nil {
		return nil, nil, errors.Wrap(err, "open kv")
	}

	_clock, err := Clock()
	if err != nil {
		return nil, nil, errors.Wrap(err, "create clock")
	}

	mode, err := kvd.Get(key.App())
	if err != nil {
		mode = []byte(consts.AppBuiltin)
	}
	app, ok := consts.Apps[string(mode)]
	if !ok {
		return nil, nil, fmt.Errorf("can't find app: %s, please try re-login", mode)
	}
	appId, appHash := app.AppID, app.AppHash

	// process proxy
	var dialer dcs.DialFunc = proxy.Direct.DialContext
	if p := viper.GetString(consts.FlagProxy); p != "" {
		d, err := utils.Proxy.GetDial(p)
		if err != nil {
			return nil, nil, errors.Wrap(err, "get dialer")
		}
		dialer = d.DialContext
	}

	opts := telegram.Options{
		Resolver: dcs.Plain(dcs.PlainOptions{
			Dial: dialer,
		}),
		ReconnectionBackoff: func() backoff.BackOff {
			return Backoff(_clock)
		},
		Device:         consts.Device,
		SessionStorage: storage.NewSession(kvd, login),
		RetryInterval:  5 * time.Second,
		MaxRetries:     -1, // infinite retries
		DialTimeout:    10 * time.Second,
		Middlewares:    middlewares,
		Clock:          _clock,
		Logger:         logger.From(ctx).Named("td"),
	}

	// test mode, hook options
	if viper.GetString(consts.FlagTest) != "" {
		appId, appHash = telegram.TestAppID, telegram.TestAppHash
		opts.DC = 2
		opts.DCList = dcs.Test()
		// add rate limit to avoid frequent flood wait
		opts.Middlewares = append(opts.Middlewares, ratelimit.New(rate.Every(100*time.Millisecond), 5))
	}

	logger.From(ctx).Info("New telegram client",
		zap.Int("app", app.AppID),
		zap.String("mode", string(mode)),
		zap.Bool("is_login", login))

	return telegram.NewClient(appId, appHash, opts), kvd, nil
}

func NoLogin(ctx context.Context, middlewares ...telegram.Middleware) (*telegram.Client, kv.KV, error) {
	mid, err := NewDefaultMiddlewares(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "create default middlewares")
	}

	return New(ctx, false, append(middlewares, mid...)...)
}

func Login(ctx context.Context, middlewares ...telegram.Middleware) (*telegram.Client, kv.KV, error) {
	mid, err := NewDefaultMiddlewares(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "create default middlewares")
	}
	return New(ctx, true, append(middlewares, mid...)...)
}

func Clock() (tdclock.Clock, error) {
	_clock := tdclock.System
	if ntp := viper.GetString(consts.FlagNTP); ntp != "" {
		var err error
		_clock, err = clock.New()
		if err != nil {
			return nil, err
		}
	}

	return _clock, nil
}

func Backoff(_clock tdclock.Clock) backoff.BackOff {
	b := backoff.NewExponentialBackOff()

	b.Multiplier = 1.1
	b.MaxElapsedTime = viper.GetDuration(consts.FlagReconnectTimeout)
	b.MaxInterval = 10 * time.Second
	b.Clock = _clock
	return b
}
