// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package pointerdb

import (
	"context"

	"go.uber.org/zap"

	"storj.io/storj/pkg/provider"
	"storj.io/storj/pkg/utils"
	proto "storj.io/storj/protos/pointerdb"
	"storj.io/storj/storage/boltdb"
)

// Config is a configuration struct that is everything you need to start a
// PointerDB responsibility
type Config struct {
	DatabaseURL          string `help:"the database connection string to use" default:"bolt://$CONFDIR/pointerdb.db"`
	MinInlineSegmentSize int64  `default:"1240" help:"minimum inline segment size"`
	MaxInlineSegmentSize int    `default:"8000" help:"maximum inline segment size"`
}

// Run implements the provider.Responsibility interface
func (c Config) Run(ctx context.Context, server *provider.Provider) error {
	dburl, err := utils.ParseURL(c.DatabaseURL)
	if err != nil {
		return err
	}
	if dburl.Scheme != "bolt" {
		return Error.New("unsupported db scheme: %s", dburl.Scheme)
	}
	bdb, err := boltdb.NewClient(zap.L(), dburl.Path, boltdb.PointerBucket)
	if err != nil {
		return err
	}
	defer func() { _ = bdb.Close() }()

	proto.RegisterPointerDBServer(server.GRPC(), NewServer(bdb, zap.L(), c))

	return server.Run(ctx)
}
