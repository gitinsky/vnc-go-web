package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "net"
    "io"
    "github.com/gorilla/websocket"
    )
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
 
func address_decode(address_bin []byte) (string,string) {
   
  var host string // = "127.0.0.1"
  var port string //= "22";
 
  return host,port
}
  
 
func forwardtcp(wsconn *websocket.Conn,conn net.Conn) {
 
  for {
    // Receive and forward pending data from tcp socket to web socket
    tcpbuffer := make([]byte, 1024)
 
    n,err := conn.Read(tcpbuffer)
    if err == io.EOF { fmt.Printf("TCP Read failed"); break; }
    if err == nil {
      fmt.Printf("Forwarding from tcp to ws: %d bytes: %s\n",n,tcpbuffer)
      print_binary(tcpbuffer)
      wsconn.WriteMessage(websocket.BinaryMessage,tcpbuffer[:n])
    } else {
        fmt.Println (err)
    }
  }
}
 
func forwardws (wsconn *websocket.Conn,conn net.Conn) {
 
 for {
    // Send pending data to tcp socket
    n,buffer,err := wsconn.ReadMessage()
    if err == io.EOF { fmt.Printf("WS Read Failed"); break; }
    if err == nil {
      s := string(buffer[:len(buffer)])
      fmt.Printf("Received (from ws) forwarding to tcp: %d bytes: %s %d\n",len(buffer),s,n)
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
  wsconn, err := upgrader.Upgrade(w, r, nil)
 
  if err != nil {
    fmt.Println(err)
    return
  }
 
  // get connection address and port
  //address := make([]byte, 16)
 
  //n,address,err := wsconn.ReadMessage()
  /*if err != nil {
    fmt.Printf("address read error");
    fmt.Printf("read %d bytes",n);  
  }*/
 
 // print_binary(address)
 
  //host, port := address_decode(address)

  conn, err := net.Dial("tcp", ip)
  if err != nil {
    fmt.Println(err)
  // handle error
  }
 
  go forwardtcp(wsconn,conn)
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
    w.Write(contents)
    return string(contents)
}
func main() {
    //var ip_string string
    http.HandleFunc("/login",func(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()
        fmt.Println(r.Form)
        fmt.Printf("%v",r.FormValue("number"))
        getIP(r.FormValue("number"), w)
    })
    http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./"))))
    http.HandleFunc("/websockify", func(w http.ResponseWriter, r *http.Request) { wsProxyHandler(w,r,"172.28.1.82:5900") })
    panic(http.ListenAndServe(":8192", nil))

}
