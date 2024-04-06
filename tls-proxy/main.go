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

	// TODO: automate certificate generation using the autocert package
	/*
	   if *domain == "" {
	       log.Fatal("Please provide a domain using the -domain flag")
	   }
	*/

	// Enable verbose logging for the autocert package
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	/*
		manager := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(*domain),
			Cache:      autocert.DirCache("certs"),
		}
		manager.Client = &acme.Client{
			DirectoryURL: autocert.DefaultACMEDirectory,
			UserAgent:    "Your-User-Agent",
			Logger:       log.New(os.Stderr, "autocert: ", log.LstdFlags),
		}

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
