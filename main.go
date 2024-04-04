package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/things-go/go-socks5"
)

func main() {
	if len(os.Args) < 5 {
		log.Fatal("Usage: go run main.go <username> <password> <listening_address> <listening_port>")
	}

	username := os.Args[1]
	password := os.Args[2]
	listeningAddress := os.Args[3]
	listeningPort := os.Args[4]

	// Validate listening address
	if ip := net.ParseIP(listeningAddress); ip == nil {
		log.Fatalf("Invalid listening address: %s", listeningAddress)
	}

	// Validate listening port
	if port, err := net.LookupPort("tcp", listeningPort); err != nil {
		log.Fatalf("Invalid listening port: %s", listeningPort)
	} else if port < 1 || port > 65535 {
		log.Fatalf("Invalid listening port: %d", port)
	}

	cator := socks5.UserPassAuthenticator{
		socks5.StaticCredentials{
			username: password,
		},
	}

	server := socks5.NewServer(
		socks5.WithAuthMethods([]socks5.Authenticator{cator}),
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
	)

	// Create SOCKS5 proxy on the specified address and port
	listeningAddr := fmt.Sprintf("%s:%s", listeningAddress, listeningPort)
	if err := server.ListenAndServe("tcp", listeningAddr); err != nil {
		panic(err)
	}

}
