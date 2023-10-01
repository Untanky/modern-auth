package main

import (
	"flag"
	"fmt"
	"github.com/Untanky/modern-auth/internal/app"
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

	cfg := app.ListenerConfig{Addr: fmt.Sprintf(":%d", *port)}
	if *useTLS {
		tlsCfg := app.TLSConfig{CertificateFilepath: *certFile, KeyFilepath: *keyFile, ListenerConfig: cfg}
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
