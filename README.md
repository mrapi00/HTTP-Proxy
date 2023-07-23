# HTTP Proxy

Golang web proxy that passes requests and data from local web clients to remote web servers for improved web performance by leveraging DNS prefetching, data caching, and content filtering specified HTML elements

![alt text]https://mahmudulrapi.netlify.app/proxy.79a14fab.png

### Starting the Proxy

``` sh
go build
./http_proxy <port> &

# testing proxy via terminal

telnet localhost <port>
GET http://www.<website>.com/ HTTP/1.1
```
