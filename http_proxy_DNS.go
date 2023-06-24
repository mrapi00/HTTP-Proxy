/*****************************************************************************
 * http_proxy_DNS.go
 * Names: Mahmudul Rapi, Pranav Thatte
 * NetIds: mrapi, pthatte
 *****************************************************************************/
// TODO: implement an HTTP proxy with DNS Prefetching

// Note: it is highly recommended to complete http_proxy.go first, then copy it
// with the name http_proxy_DNS.go, thus overwriting this file, then edit it
// to add DNS prefetching (don't forget to change the filename in the header
// to http_proxy_DNS.go in the copy of http_proxy.go)
package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/html"
	"net/http"
	"net/url"
	"os"
	"strings"
)

/* server()
 * Open socket and wait for client to connect
 * Print received message to stdout
 */
func server(server_port string) {
	socket, err := net.Listen("tcp", fmt.Sprint("127.0.0.1:", server_port))
	if err != nil {
		log.Fatalf("Could not listen on port %v", server_port)
	}

	for {
		conn, err := socket.Accept()

		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// need to eventuall close connection
	defer conn.Close()
	// read in connection from telnet
	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)

	serverError := []byte("HTTP/1.1 500 Internal Server Error\r\n")

	if err != nil || req.Method != "GET" {
		conn.Write(serverError)
		return
	}

	req.Header.Add("Connection", "close")
	// using ReadRequest writes Request into URI instead of URL field, need to swap fields
	// https://cs.opensource.google/go/go/+/refs/tags/go1.20.2:src/net/http/request.go;l=291
	urlLink, err := url.Parse(req.RequestURI)
	if err != nil {
		conn.Write(serverError)
		return
	}
	req.URL = urlLink
	req.RequestURI = ""

	// open client socket to online server

	// https://stackoverflow.com/questions/23297520/how-can-i-make-the-go-http-client-not-follow-redirects-automatically
	var ErrUseLastResponse = errors.New("net/http: use last response")
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return ErrUseLastResponse
		},
	}

	// proxy is server to telnet client (specified host name i.e. google.com), client to web server
	resp, err := client.Do(req)
	if err != nil {
		// https://stackoverflow.com/questions/22170942/how-can-i-get-an-error-message-in-a-string-in-go
		if !strings.Contains(fmt.Sprint(err), "net/http: use last response") {
			conn.Write(serverError)
			return
		}
	}

	// duplicate input stream for DNS
	rawBody, err := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
	resp.Write(conn)

	z := html.NewTokenizer(bytes.NewReader(rawBody))
	go func() {
		for {
			tt := z.Next()
			switch tt {
			case html.ErrorToken:
				return
			case html.StartTagToken:
				tn, _ := z.TagName()
				if len(tn) == 1 && tn[0] == 'a' { // the starting <a> tag

					k, v, _ := z.TagAttr() // <a href=<URL>>, then k = href, v=<URL>
					if string(k) == "href" && strings.HasPrefix(string(v), "http") {
						go net.LookupHost(string(v))
					}
				}
			}
		}
	}()

}

// Main parses command-line arguments and calls server function
func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./server [server port]")
	}
	server_port := os.Args[1]
	server(server_port)
}
