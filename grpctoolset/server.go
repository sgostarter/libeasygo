package grpctoolset

import (
	"context"
	"net"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
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

func NewGRPCServer(routineMan routineman.RoutineMan, cfg *GRPCServerConfig, logger l.Wrapper, extraInterceptors ...interface{}) (GRPCServer, error) {
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
		routineMan:        routineMan,
		address:           cfg.Address,
		extraInterceptors: extraInterceptors,
		logger:            logger.WithFields(l.StringField(l.ClsKey, "gRPCServerImpl")),
		serverOptions:     serverOptions,
	}, nil
}

type gRPCServerImpl struct {
	lock sync.Mutex

	routineMan        routineman.RoutineMan
	address           string
	extraInterceptors []interface{}
	logger            l.Wrapper
	serverOptions     []grpc.ServerOption

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

	impl.s = grpc.NewServer(impl.getServerOptions()...)
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
	// impl.s.GracefulStop()
	impl.s.Stop()
}

func (impl *gRPCServerImpl) StopAndWait() {
	impl.routineMan.TriggerStop()
	// impl.s.GracefulStop()
	impl.s.Stop()
	impl.routineMan.StopAndWait()
}

func (impl *gRPCServerImpl) getInterceptors() []grpc.ServerOption {
	var interceptors []grpc.UnaryServerInterceptor

	var streamInterceptors []grpc.StreamServerInterceptor

	interceptors = append(interceptors, grpc_recovery.UnaryServerInterceptor())
	streamInterceptors = append(streamInterceptors, grpc_recovery.StreamServerInterceptor())

	for _, v := range impl.extraInterceptors {
		// 不是很明白type出来的和直接写func有什么区别，但这俩type在switch的时候确实不一样
		// 而且case用逗号也不行，也很疑惑
		switch interceptor := v.(type) {
		case func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error):
			interceptors = append(interceptors, interceptor)
		case func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error:
			streamInterceptors = append(streamInterceptors, interceptor)
		case grpc.UnaryServerInterceptor:
			interceptors = append(interceptors, interceptor)
		case grpc.StreamServerInterceptor:
			streamInterceptors = append(streamInterceptors, interceptor)
		default:
			impl.logger.Warn("interceptor not valid")
		}
	}

	return []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)),
	}
}

func (impl *gRPCServerImpl) getServerOptions() (options []grpc.ServerOption) {
	options = append(options, impl.serverOptions...)
	options = append(options, impl.getInterceptors()...)

	return
}
