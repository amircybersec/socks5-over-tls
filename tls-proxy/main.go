package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// Parse command-line flags
	//domain := flag.String("domain", "", "Domain for which to generate a certificate")
	targetAddr := flag.String("target", "localhost:8080", "Target server address")
	flag.Parse()

	/*
	   if *domain == "" {
	       log.Fatal("Please provide a domain using the -domain flag")
	   }
	*/

	// Enable verbose logging for the autocert package
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//autocert.DefaultACME.Log = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	/*
	   certManager := &autocert.Manager{
	       Prompt:     autocert.AcceptTOS,
	       HostPolicy: autocert.HostWhitelist(*domain),
	       Cache:      autocert.DirCache("/root/tls-proxy/certs"), // Directory to store certificates
	   }
	   log.Printf("certmanager is %v", certManager)

	   // TLS server configuration
	    tlsConfig := &tls.Config{
	       GetCertificate: certManager.GetCertificate,
	   }
	*/

	certFile := "fullchain.pem"
	keyFile := "privkey.pem"
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatalf("Failed to load certificate and key: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
	}

	//tlsConfig := certManager.TLSConfig()

	// Listen for incoming TLS connections
	listener, err := tls.Listen("tcp", ":443", tlsConfig)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	log.Println("TLS server listening on :443")

	for {
		// Accept incoming TLS connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// Handle the connection in a separate goroutine
		go handleConnection(conn, *targetAddr)
	}
}

func handleConnection(conn net.Conn, targetAddr string) {
	defer conn.Close()

	// Connect to the target server
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Failed to connect to target server: %v", err)
		return
	}
	defer targetConn.Close()

	// Perform bidirectional data transfer between the client and target server
	go io.Copy(targetConn, conn)
	io.Copy(conn, targetConn)
}
