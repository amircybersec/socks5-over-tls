# socks5-over-tls

In this repo, I explain the steps to setup a server to act as proxy that supports socks5-over-tls transport.

socks5-over-tls transport is a wrapped transport where the outer transport is TLS and the inner transport is SOCKS5. The transport essentially looks like a normal TLS. It starts with a TLS handshake and payload is encrypted. The payload contains a SOCKS5 connection, which starts with username and password authentication, and relays the TCP payload just like a normal SOCKS5. 

## Client

For the client side, I use Outline SDK configurable dialers to build a custom config like this:

```tls:sni=open.com|socks5://user:pass@serverdomain.com:443```

where `open.com` is the domain that can be used to bypass SNI based filtering. The `serverdomain.com` is the target server domain name. In the current implementation you need a valid CA-signed certitcate for your domain. `user:pass` are your socks5 server credentials. 

For the client, you can use fyne-proxy example, or Blazer proxy if you need a GUI based client. Both clients run on `windows`, `macOS`, and `linux`. I will update this doc as more clients become available. 

## Server

The server side basicallt has two elements: 
1. A TLS/SSL server
This server sits at the front and listens on the port `443` for incoming connections. It performs TLS handshake with the client and decypts the traffic. It then passes the traffic to the local SOCK5 server. 

2. SOCKS5 server
The SOCKS5 server sits behind the TLS server and handles the SOCKS5 connection.

In the current PoC, I use `socat` as TLS server to perform SSL handshake and passed the traffic to a local `socks5` server (running on port 8080 in this example) after decrypting it. I used `certbot` to create a CA-signed certificate.

```
socat -dd openssl-listen:443,reuseaddr,cert=fullchain.pem,key=privkey.pem,verify=0,fork tcp-connect:localhost:8080
```

You can run the socks5 sever on command line:

```
./main foo bar 127.0.0.1 8080
```

Notes:

1 - This is just a PoC at this point and I am planning to merge these two elements into one executable and run it as systemd service. Please also make sure you that you open port 443 on your server with `ufw` and set your server firewall properly.


2 - You can add packet splits and TLS fragmentation to the dialer make fingerprinting based on packet sizes more difficult:

```
split:1|tlsfrag:1|tls:sni=open.com|socks5://user:pass@serverdomain.com:443|plit:1
```


