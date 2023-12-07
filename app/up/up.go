package up

import (
	"context"

	"github.com/fatih/color"
	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
	"github.com/spf13/viper"
	"go.uber.org/multierr"

	"github.com/iyear/tdl/app/internal/tgc"
	"github.com/iyear/tdl/pkg/consts"
	"github.com/iyear/tdl/pkg/dcpool"
	"github.com/iyear/tdl/pkg/prog"
	"github.com/iyear/tdl/pkg/storage"
	"github.com/iyear/tdl/pkg/uploader"
	"github.com/iyear/tdl/pkg/utils"
)

type Options struct {
	Chat     string
	Paths    []string
	Excludes []string
	Remove   bool
	Photo    bool
}

func Run(ctx context.Context, opts *Options) error {
	files, err := walk(opts.Paths, opts.Excludes)
	if err != nil {
		return err
	}

	color.Blue("Files count: %d", len(files))

	c, kvd, err := tgc.NoLogin(ctx)
	if err != nil {
		return err
	}

	return tgc.RunWithAuth(ctx, c, func(ctx context.Context) (rerr error) {
		middlewares, err := tgc.NewDefaultMiddlewares(ctx)
		if err != nil {
			return errors.Wrap(err, "create middlewares")
		}

		pool := dcpool.NewPool(c, int64(viper.GetInt(consts.FlagPoolSize)), middlewares...)
		defer multierr.AppendInvoke(&rerr, multierr.Close(pool))

		manager := peers.Options{Storage: storage.NewPeers(kvd)}.Build(pool.Default(ctx))

		to, err := resolveDestPeer(ctx, manager, opts.Chat)
		if err != nil {
			return errors.Wrap(err, "get target peer")
		}

		upProgress := prog.New(utils.Byte.FormatBinaryBytes)
		upProgress.SetNumTrackersExpected(len(files))
		prog.EnablePS(ctx, upProgress)

		options := uploader.Options{
			Client:   pool.Default(ctx),
			PartSize: viper.GetInt(consts.FlagPartSize),
			Threads:  viper.GetInt(consts.FlagThreads),
			Iter:     newIter(files, to, opts.Photo, opts.Remove),
			Progress: newProgress(upProgress),
		}

		up := uploader.New(options)

		go upProgress.Render()
		defer prog.Wait(ctx, upProgress)

		return up.Upload(ctx, viper.GetInt(consts.FlagLimit))
	})
}

func resolveDestPeer(ctx context.Context, manager *peers.Manager, chat string) (peers.Peer, error) {
	if chat == "" {
		return manager.FromInputPeer(ctx, &tg.InputPeerSelf{})
	}

	return utils.Telegram.GetInputPeer(ctx, manager, chat)
}
