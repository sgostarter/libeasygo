package grpctoolset

import (
	"context"
	"net"
	"sync"

	"github.com/sgostarter/i/l"
	"github.com/sgostarter/libeasygo/commerr"
	"github.com/sgostarter/libeasygo/routineman"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCServerConfig struct {
	Address   string         `yaml:"Address" json:"address"`
	TLSConfig *GRPCTlsConfig `yaml:"TLSConfig" json:"tls_config"`
}

type GRPCServer interface {
	Start(init func(s *grpc.Server)) (err error)
	Wait()
	Stop()
	StopAndWait()
}

func NewGRPCServer(routineMan routineman.RoutineMan, cfg *GRPCServerConfig, logger l.Wrapper) (GRPCServer, error) {
	if routineMan == nil {
		routineMan = routineman.NewRoutineMan(context.Background(), logger)
	}

	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	serverOptions := make([]grpc.ServerOption, 0)

	if cfg.TLSConfig != nil {
		tlsConfig, err := GenServerTLSConfig(cfg.TLSConfig)
		if err != nil {
			return nil, err
		}

		serverOptions = append(serverOptions, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}

	return &gRPCServerImpl{
		routineMan:    routineMan,
		address:       cfg.Address,
		logger:        logger.WithFields(l.StringField(l.ClsKey, "gRPCServerImpl")),
		serverOptions: serverOptions,
	}, nil
}

type gRPCServerImpl struct {
	lock sync.Mutex

	routineMan    routineman.RoutineMan
	address       string
	logger        l.Wrapper
	serverOptions []grpc.ServerOption

	listen net.Listener
	s      *grpc.Server
}

func (impl *gRPCServerImpl) Start(init func(s *grpc.Server)) (err error) {
	impl.lock.Lock()
	defer impl.lock.Unlock()

	if impl.listen != nil || impl.s != nil {
		err = commerr.ErrAlreadyExists

		return
	}

	impl.listen, err = net.Listen("tcp", impl.address)
	if err != nil {
		impl.logger.WithFields(l.StringField("listen", impl.address), l.ErrorField(err)).Error("listen")

		return
	}

	impl.s = grpc.NewServer(impl.serverOptions...)
	init(impl.s)

	impl.routineMan.StartRoutine(impl.mainRoutine, "mainRoutine")

	return
}

func (impl *gRPCServerImpl) mainRoutine(ctx context.Context, exiting func() bool) {
	err := impl.s.Serve(impl.listen)
	if err != nil {
		impl.logger.WithFields(l.ErrorField(err)).Error("GRPCServe")
	}
}

func (impl *gRPCServerImpl) Wait() {
	impl.routineMan.Wait()
}

func (impl *gRPCServerImpl) Stop() {
	impl.s.GracefulStop()
}

func (impl *gRPCServerImpl) StopAndWait() {
	impl.s.GracefulStop()
	impl.routineMan.StopAndWait()
}
