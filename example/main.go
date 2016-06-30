package main

import (
	"crypto/tls"
	"log"
	"net"
	"sync/atomic"

	"websocket"
)

var i uint32 = 0

func main() {
	tlsEnable := true

	rawl, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Panic(err)
	}

	var ln net.Listener
	if tlsEnable {
		config := tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			Certificates:             make([]tls.Certificate, 1),
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			},
		}
		config.Certificates[0], err = tls.LoadX509KeyPair("./ssl/server.crt", "./ssl/server.key")
		if err != nil {
			log.Fatalf("load certFile or keyFile failed!err:=%v", err)
			return
		}
		config.NextProtos = append(config.NextProtos, "http/1.1")
		ln = tls.NewListener(rawl.(*net.TCPListener), &config)
	} else {
		ln = rawl
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept err:", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	gi := atomic.AddUint32(&i, 1)
	log.Printf("this is no.%d", gi)
	if tlsConn, ok := conn.(*tls.Conn); ok {
		if err := tlsConn.Handshake(); err != nil {
			log.Printf("http: TLS handshake error from %s: %v", tlsConn.RemoteAddr(), err)
			return
		}
	}

	wssocket := websocket.NewWsSocket(conn)
	err := wssocket.HandShake()
	if err != nil {
		return
	}
	for {
		_, err = wssocket.ReadIframe()
		if err != nil {
			log.Println("readIframe err:", err)
			return
		}
	}
}
