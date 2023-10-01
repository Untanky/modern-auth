package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/Untanky/modern-auth/internal/core"
	"github.com/Untanky/modern-auth/registry"
	"google.golang.org/grpc"
	"net"
)

var (
	port     = flag.Int("port", 5500, "the port to run the registry on")
	useTLS   = flag.Bool("useTLS", false, "use useTLS")
	certFile = flag.String("certFile", "", "path to the cert file")
	keyFile  = flag.String("keyFile", "", "path to the key file")
)

func main() {
	flag.Parse()

	var listener net.Listener
	var err error

	cfg := ListenerConfig{Addr: fmt.Sprintf(":%d", *port)}
	if *useTLS {
		tlsCfg := TLSConfig{CertificateFilepath: *certFile, KeyFilepath: *keyFile, ListenerConfig: cfg}
		listener, err = tlsCfg.NewListener()
	} else {
		listener, err = cfg.NewListener()
	}
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption

	store := core.NewInMemoryKeyValueStore[*registry.RegistrationInfo]()
	index := core.NewInMemoryKeyValueStore[core.List[string]]()
	registerChan := make(chan *registry.RegistrationInfo)
	unregisterChan := make(chan *registry.RegistrationInfo)

	grpcServer := grpc.NewServer(opts...)
	registry.RegisterRegistryServer(grpcServer, registry.NewRegistryServer(store, index, registerChan, unregisterChan))
	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}
}

type ListenerConfig struct {
	Addr string
}

func (cfg *ListenerConfig) NewListener() (net.Listener, error) {
	return net.Listen("tcp", cfg.Addr)
}

type TLSConfig struct {
	CertificateFilepath string
	KeyFilepath         string
	ListenerConfig
}

func (cfg *TLSConfig) NewListener() (net.Listener, error) {
	listener, err := cfg.ListenerConfig.NewListener()
	if err != nil {
		return nil, err
	}

	if *useTLS {
		cert, err := tls.LoadX509KeyPair(cfg.CertificateFilepath, cfg.KeyFilepath)
		if err != nil {
			return nil, err
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
		listener = tls.NewListener(listener, config)
	}
	return listener, nil
}
