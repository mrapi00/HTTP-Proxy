/*****************************************************************************
 * http_proxy.go                                                                 
 * Names: Mahmudul Rapi, Pranav Thatte
 * NetIds: mrapi, pthatte
 *****************************************************************************/

 // TODO: implement an HTTP proxy
 
 package main

 import (
   "bufio"
   "fmt"
   "log"
   "os"
   "net"
   "net/url"
   "net/http"
   "errors"
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
	    go handleConnection(conn);
   }
 }

 func handleConnection(conn net.Conn) {
  // need to eventuall close connection
  defer conn.Close()
  // read in connection from telnet
  reader := bufio.NewReader(conn)
  req, err := http.ReadRequest(reader) 
  
  serverError := []byte("HTTP/1.1 500 Internal Server Error\r\n")
  
  if (err != nil || req.Method != "GET") {
    conn.Write(serverError)
    return
  }

  req.Header.Add("Connection", "close")
  // using ReadRequest writes Request into URI instead of URL field, need to swap fields
  // https://cs.opensource.google/go/go/+/refs/tags/go1.20.2:src/net/http/request.go;l=291
  urlLink, err := url.Parse(req.RequestURI)
  if (err != nil) {
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
  if (err != nil) {
    // https://stackoverflow.com/questions/22170942/how-can-i-get-an-error-message-in-a-string-in-go
    if (!strings.Contains(fmt.Sprint(err), "net/http: use last response")) {
      conn.Write(serverError)
      return
    }
  }
  resp.Write(conn)
}
 
 // Main parses command-line arguments and calls server function
 func main() {
   if len(os.Args) != 2 {
     log.Fatal("Usage: ./server [server port]")
   }
   server_port := os.Args[1]
   server(server_port)
 }
 