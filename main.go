package main

import (
	"log"
	"os"

	"github.com/things-go/go-socks5"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <username> <password>")
	}

	username := os.Args[1]
	password := os.Args[2]

	cator := socks5.UserPassAuthenticator{
		socks5.StaticCredentials{
			username: password,
		},
	}

	server := socks5.NewServer(
		socks5.WithAuthMethods([]socks5.Authenticator{cator}),
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
	)

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", ":8080"); err != nil {
		panic(err)
	}

}
