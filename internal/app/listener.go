package app

import (
	"crypto/tls"
	"net"
)

type ListenerConfig struct {
	Addr string
}

func (cfg *ListenerConfig) NewListener() (net.Listener, error) {
	addr, err := net.ResolveTCPAddr("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	return net.ListenTCP("tcp", addr)
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

	cert, err := tls.LoadX509KeyPair(cfg.CertificateFilepath, cfg.KeyFilepath)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	listener = tls.NewListener(listener, config)
	return listener, nil
}
