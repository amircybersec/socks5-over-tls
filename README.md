# socks5-over-tls

In this repo, I explain the steps to setup a server to act as proxy that supports socks5-over-tls transport.

socks5-over-tls transport is a wrapped transport where the outer transport is TLS and the inner transport is SOCKS5. The transport essentially looks like a normal TLS. It starts with a TLS handshake and payload is encrypted. The payload contains a SOCKS5 connection, which starts with username and password authentication, and relays the TCP payload just like a normal SOCKS5. 

For the client side, I use Outline SDK configurable dialers to build a custom config like this:




