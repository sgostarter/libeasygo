package grpctoolset

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClientConfig struct {
	Target    string         `yaml:"Target" json:"target"`
	TLSConfig *GRPCTlsConfig `yaml:"TLSConfig" json:"tls_config"`
}

func DialGRPC(cfg *GRPCClientConfig) (*grpc.ClientConn, error) {
	var dialOptions []grpc.DialOption

	if cfg.TLSConfig != nil {
		tlsConfig, err := GenClientTLSConfig(cfg.TLSConfig)
		if err != nil {
			return nil, err
		}

		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	return grpc.Dial(cfg.Target, dialOptions...)
}
