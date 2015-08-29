package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "net"
    "io"
    "github.com/gorilla/websocket"
    "time"
    "strings"
    )
const (
  // Time allowed to write a message to the peer.
  writeWait = 10 * time.Second

  // Time allowed to read the next pong message from the peer.
  pongWait = 60 * time.Second

  // Send pings to peer with this period. Must be less than pongWait.
  pingPeriod = (pongWait * 9) / 10

  // Maximum message size allowed from peer.
  maxMessageSize = 512
)
 
func address_decode(address_bin []byte) (string,string) {
   
  var host string // = "127.0.0.1"
  var port string //= "22";
 
  return host,port
}
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}
 
 
       
func print_binary(s []byte) {
  fmt.Printf("print b:");
  for n := 0;n < len(s);n++ {
    fmt.Printf("%d,",s[n]);
  }
  fmt.Printf("\n");
}
 

 
func forwardtcp(wsconn *websocket.Conn,conn net.Conn,r *http.Request ) {
 
  for {
    // Receive and forward pending data from tcp socket to web socket
    tcpbuffer := make([]byte, 1024)
 
    //n,err := conn.Read(tcpbuffer)
    n,err := conn.Read(tcpbuffer)
    if err == io.EOF { fmt.Printf("TCP Read failed"); break; }
    if err == nil {
      //fmt.Printf("Forwarding from tcp to ws: %d bytes: %s\n",n,tcpbuffer)
      //print_binary(tcpbuffer)
      //r.Header.Add("Sec-WebSocket-Protocol", "chat")
      wsconn.WriteMessage(websocket.BinaryMessage,tcpbuffer[:n])
    } else {
        fmt.Println (err)
    }
  }
}
 
func forwardws (wsconn *websocket.Conn,conn net.Conn) {
 
 for {
    // Send pending data to tcp socket

    //n,buffer,err := wsconn.ReadMessage()
    _,buffer,err := wsconn.ReadMessage()
    if err == io.EOF { fmt.Printf("WS Read Failed"); break; }
    if err == nil {
      //s := string(buffer[:len(buffer)])
     // fmt.Printf("Received (from ws) forwarding to tcp: %d bytes: %s %d\n",len(buffer),s,n)
      //print_binary(buffer)
      conn.Write(buffer)
    } else{
        //fmt.Println(err)
    }
  }
}
 
func wsProxyHandler(w http.ResponseWriter, r *http.Request, ip string) {
  //r.ParseForm()
  //fmt.Println(r.FormValue("IP"))
  header:=http.Header{}
  header.Add("Sec-WebSocket-Protocol", "binary")
  wsconn, err := upgrader.Upgrade(w, r, header)
  wsconn.SetReadDeadline(time.Now().Add(pongWait))
  wsconn.SetPongHandler(func(string) error { wsconn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
  if err != nil {
    fmt.Println(err)
    return
  }

  fmt.Println( websocket.Subprotocols(r))
  //get connection address and port
  /*address := make([]byte, 16)
  n,address,err := wsconn.ReadMessage()
  if err != nil {
    fmt.Printf("address read error");
    fmt.Printf("read %d bytes",n);  
  }
 
  print_binary(address)
 
  host, port := address_decode(address)
 fmt.Println(host,port)
 */
  conn, err := net.Dial("tcp", ip)
  if err != nil {
    fmt.Println(err)
  // handle error
  }
 
  go forwardtcp(wsconn,conn,r)
  go forwardws(wsconn,conn)
 
  fmt.Printf("websocket closed");
}
type ConnectionInfo struct {
    ServerNumber string
    UUID  string
    IP string
}
func getIP(servernumber string, w http.ResponseWriter) string {
    connectionInfo:=ConnectionInfo{}
    connectionInfo.ServerNumber=servernumber
    UUIDAdress:="http://127.0.0.1/test.php?servernumber="+connectionInfo.ServerNumber
    response,err := http.Get(UUIDAdress)
    if err != nil {
        fmt.Printf("%s", err)
        return ""
    }
    defer response.Body.Close()
    contents, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Printf("%s", err)
        return ""
    }
    fmt.Println(string(contents))
    connectionInfo.UUID=string(contents)
    IPAdress:="http://127.0.0.1/test.php?uuid="+connectionInfo.UUID
    response,err = http.Get(IPAdress)
    if err != nil {
        fmt.Printf("%s", err)
        return ""
    }
    defer response.Body.Close()
    contents, err = ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Printf("%s", err)
        return ""
    }
    fmt.Println(string(contents))
    connectionInfo.IP=string(contents)
    ip_cookie := http.Cookie {
            Name: "IPAdress",
            Value: string(contents),
            MaxAge: 12345, //TODO: change this
            Secure: false,
            HttpOnly: false,
            Domain: "me.not",
    }
    http.SetCookie(w,&ip_cookie)
    w.Write(contents)


        
    return string(contents)
}
func main() {
    http.HandleFunc("/login",func(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()
        fmt.Println(r.Form)
        fmt.Printf("%v",r.FormValue("number"))
        getIP(r.FormValue("number"), w)
    })
    http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./"))))
    http.HandleFunc("/websockify", func(w http.ResponseWriter, r *http.Request) { 
    if r.Method != "GET" {
      http.Error(w, "Method not allowed", 405)
      return
    }
    IP,err:=r.Cookie("IPAdress")
    if err!=nil {
      fmt.Printf("%s",err)
      return
    }
    fmt.Println(IP.String())
    IPSplit:=strings.Split(IP.String(),"=")
    if(len(IPSplit)==2){
      //r.Header.Add("Sec-WebSocket-Protocol", "chat")
      //r.Write([]byte("Here is a string...."))
      wsProxyHandler(w,r,IPSplit[1]) 
    } else {
      fmt.Println("incorrectIP")
      return
    }
    })
    panic(http.ListenAndServe(":8192", nil))

}
