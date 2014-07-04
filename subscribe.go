package wemo

/* 
  Need to set up a subscription service to Wemo events.

  1. Take a discovered device and
  2. Send a subscribe message to 
    deviceIP:devicePort/upnp/event/basicevent1
  3. If the responce is 200, the subscription is successful and ...
  4. ... thus it should be added to the subscribed device list
  5. When state is emitted record state changes
  6. Should then send a 200 status response?
  
*/

import(
  "net"
  "fmt"
  "io"
  "log"
  "net/http"
)


func Listener(listenerAddress string){
  
  fmt.Println("Listening... ", listenerAddress)
 
	listener, err := net.Listen("tcp", listenerAddress)
	if err != nil {
		log.Fatal("Listen Err: ", err)
    return
	}
  
	//defer listener.Close()
  
	for {
    
		// Wait for a connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
    
    fmt.Println(conn)
    
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			// Echo all incoming data.
			fmt.Println(io.Copy(c, c))
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}

// Subscribe to the device event emitter, return the Subscription ID (sid) and StatusCode
func (self *Device) Subscribe(listenerAddress string) (string, int) {
  
  host := self.Host 
  
  address := fmt.Sprintf("http://%s/upnp/event/basicevent1", host)
  
  timeout := 300
  
  client := &http.Client{}
  
  req, err := http.NewRequest("SUBSCRIBE", address, nil)
  if err != nil {
    log.Fatal("http NewRequest Err: ", err)
  }
  
  req.Header.Add("host", fmt.Sprintf("http://%s", host))
  req.Header.Add("path", "/upnp/event/basicevent1")
	req.Header.Add("callback", fmt.Sprintf("<http://%s/listener>", listenerAddress))
	req.Header.Add("nt"      , "upnp:event")
	req.Header.Add("timeout" , fmt.Sprintf("Second-%d", timeout))
  
  resp, err := client.Do(req)
  if err != nil {
    log.Fatal("Client Request Error: ", err)
  }
  defer resp.Body.Close()
  
  if resp.StatusCode == 200 {
    fmt.Println("Subscription Successful: ", resp.StatusCode)
    return resp.Header.Get("Sid"), resp.StatusCode
  } else if resp.StatusCode == 400 {
    fmt.Println("Subscription Unsuccessful, Incompatible header fields: ", resp.StatusCode)
  } else if resp.StatusCode == 412 {
    fmt.Println("Subscription Unsuccessful, Precondition Failed: ", resp.StatusCode)
  } else {
    fmt.Println("Subscription Unsuccessful, Unable to accept renewal: ", resp.StatusCode)
  }

  return "", resp.StatusCode  
  
}

// According to the spec all subscribers must unsubscribe when the publisher is no longer required to provide state updates. Return the StatusCode
func (self *Device) UnSubscribe(sid string) (int){
  
  host := self.Host 
  
  address := fmt.Sprintf("http://%s/upnp/event/basicevent1", host)
  
  client := &http.Client{}
  
  req, err := http.NewRequest("UNSUBSCRIBE", address, nil)
  if err != nil {
    log.Fatal("http NewRequest Err: ", err)
  }
  
  req.Header.Add("host", fmt.Sprintf("http://%s", host))
  req.Header.Add("SID", sid)
  
  resp, err := client.Do(req)
  if err != nil {
    log.Fatal("Client Request Error: ", err)
  }
  defer resp.Body.Close()
  
  if resp.StatusCode == 200 {
    fmt.Println("Unsubscription Successful: ", resp.StatusCode)
  } else if resp.StatusCode == 400 {
    fmt.Println("Unsubscription Unsuccessful, Incompatible header fields: ", resp.StatusCode)
  } else if resp.StatusCode == 412 {
    fmt.Println("Unsubscription Unsuccessful, Precondition Failed: ", resp.StatusCode)
  } else {
    fmt.Println("Unsubscription Unsuccessful, Unable to accept renewal: ", resp.StatusCode)
  }
  
  return resp.StatusCode
  
}

// The subscription to the device must be renewed before the timeout. Return the Subscription ID (sid) and StatusCode
func (self *Device) ReSubscribe(sid string) (string, int){
  
  host := self.Host 
  
  address := fmt.Sprintf("http://%s/upnp/event/basicevent1", host)
  
  timeout := 300
  
  client := &http.Client{}
  
  req, err := http.NewRequest("SUBSCRIBE", address, nil)
  if err != nil {
    log.Fatal("http NewRequest Err: ", err)
  }
  
  req.Header.Add("host", fmt.Sprintf("http://%s", host))
  req.Header.Add("SID", sid)
  req.Header.Add("timeout" , fmt.Sprintf("Second-%d", timeout))
  
  resp, err := client.Do(req)
  if err != nil {
    log.Fatal("Client Request Error: ", err)
  }
  defer resp.Body.Close()
  
  if resp.StatusCode == 200 {
    fmt.Println("Unsubscription Successful: ", resp.StatusCode)
    return resp.Header.Get("Sid"), resp.StatusCode
  } else if resp.StatusCode == 400 {
    fmt.Println("Unsubscription Unsuccessful, Incompatible header fields: ", resp.StatusCode)
  } else if resp.StatusCode == 412 {
    fmt.Println("Unsubscription Unsuccessful, Precondition Failed: ", resp.StatusCode)
  } else {
    fmt.Println("Unsubscription Unsuccessful, Unable to accept renewal: ", resp.StatusCode)
  }
  
  return "", resp.StatusCode
  
}