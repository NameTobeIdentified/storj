// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package server

import (
	"crypto"
	"crypto/ecdsa"
	"log"
	"os"
	"path/filepath"

	"github.com/gtank/cryptopasta"
	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"gopkg.in/spacemonkeygo/monkit.v2"

	"storj.io/storj/pkg/peertls"
	"storj.io/storj/pkg/piecestore"
	"storj.io/storj/pkg/piecestore/rpc/server/psdb"
	"storj.io/storj/pkg/provider"
	pb "storj.io/storj/protos/piecestore"
)

var (
	mon = monkit.Package()

	// ServerError wraps errors returned from Server struct methods
	ServerError = errs.Class("PSServer error")
)

// Config contains everything necessary for a server
type Config struct {
	Path string `help:"path to store data in" default:"$CONFDIR"`
}

// Run implements provider.Responsibility
func (c Config) Run(ctx context.Context, server *provider.Provider) (err error) {
	defer mon.Task()(&ctx)(&err)

	s, err := Initialize(ctx, c, server.Identity().Key)
	if err != nil {
		return err
	}

	go func() {
		err := s.DB.DeleteExpiredLoop(ctx)
		zap.S().Fatal("Error in DeleteExpiredLoop: %v\n", err)
	}()

	pb.RegisterPieceStoreRoutesServer(server.GRPC(), s)

	defer func() {
		log.Fatal(s.Stop(ctx))
	}()

	return server.Run(ctx)
}

// Server -- GRPC server meta data used in route calls
type Server struct {
	DataDir string
	DB      *psdb.PSDB
	pkey    crypto.PrivateKey
}

// Initialize -- initializes a server struct
func Initialize(ctx context.Context, config Config, pkey crypto.PrivateKey) (*Server, error) {
	dbPath := filepath.Join(config.Path, "piecestore.db")
	dataDir := filepath.Join(config.Path, "piece-store-data")

	psDB, err := psdb.OpenPSDB(ctx, dataDir, dbPath)
	if err != nil {
		return nil, err
	}

	return &Server{DataDir: dataDir, DB: psDB, pkey: pkey}, nil
}

// Stop the piececstore node
func (s *Server) Stop(ctx context.Context) (err error) {
	return s.DB.Close()
}

// Piece -- Send meta data about a stored by by Id
func (s *Server) Piece(ctx context.Context, in *pb.PieceId) (*pb.PieceSummary, error) {
	log.Printf("Getting Meta for %s...", in.Id)

	path, err := pstore.PathByID(in.GetId(), s.DataDir)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Read database to calculate expiration
	ttl, err := s.DB.GetTTLByID(in.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully retrieved meta for %s.", in.Id)
	return &pb.PieceSummary{Id: in.GetId(), Size: fileInfo.Size(), ExpirationUnixSec: ttl}, nil
}

// Delete -- Delete data by Id from piecestore
func (s *Server) Delete(ctx context.Context, in *pb.PieceDelete) (*pb.PieceDeleteSummary, error) {
	log.Printf("Deleting %s...", in.Id)

	if err := s.deleteByID(in.GetId()); err != nil {
		return nil, err
	}

	log.Printf("Successfully deleted %s.", in.Id)
	return &pb.PieceDeleteSummary{Message: OK}, nil
}

func (s *Server) deleteByID(id string) error {
	if err := pstore.Delete(id, s.DataDir); err != nil {
		return err
	}

	if err := s.DB.DeleteTTLByID(id); err != nil {
		return err
	}

	log.Printf("Deleted data of id (%s) from piecestore\n", id)

	return nil
}

func (s *Server) verifySignature(ctx context.Context, ba *pb.RenterBandwidthAllocation) error {
	pi, err := provider.PeerIdentityFromContext(ctx)
	if err != nil {
		return err
	}

	k, ok := pi.Leaf.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return peertls.ErrUnsupportedKey.New("%T", pi.Leaf.PublicKey)
	}

	if ok := cryptopasta.Verify(ba.GetData(), ba.GetSignature(), k); !ok {
		return ServerError.New("Failed to verify Signature")
	}
	return nil
}
