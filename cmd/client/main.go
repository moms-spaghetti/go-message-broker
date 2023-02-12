package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"task2/internal/message"
	"time"
)

const (
	tcpaddr               = ":8181"
	tcpnetwork            = "tcp"
	udpnetwork            = "udp"
	udpaddr               = ":9001"
	apicreateSubscriber   = "createSubscriber"
	apisubscribeToTopic   = "subscribeToTopic"
	apicreateTopic        = "createTopic"
	apipublishMessage     = "publishMessage"
	apigetMessages        = "getMessages"
	apiudpGetMessageCount = "udpGetMessageCount"
	apiudpGetNextMessage  = "udpGetNextMessage"
	apiudpCompleteMessage = "udpCompleteMessage"
)

var (
	laddr = &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 9002,
		Zone: "",
	}
	raddr = &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 9001,
		Zone: "",
	}
)

type rawRequest struct {
	API     string                 `json:"api"`
	Query   map[string]interface{} `json:"query"`
	Payload map[string]interface{} `json:"Payload"`
}

type jsonResponse struct {
	Err    string      `json:"err"`
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func main() {

	menuSelect()
}

func menuSelect() {
	// option loop
	var protocol string
	for {
		fmt.Println(`
		**Protocol**
		1: TCP
		2: UDP
		q: QUIT`)

		protocol = readEntry("choose and option 1, 2, q")

		switch strings.ToLower(protocol) {
		case "1":
			protocol = tcpnetwork
		case "2":
			protocol = udpnetwork
		case "q":
			fmt.Printf("quitting\n")
			os.Exit(0)
		default:
			fmt.Printf("\n\ninvalid option\n\n")
			continue
		}

		break
	}

	var (
		request rawRequest
		conn    net.Conn
		err     error
	)

	for {
		fmt.Println("***" + protocol + " mode***")
		fmt.Println(`
		**Publisher**
		1: create topic
		2: publish message
		
		**Subscriber**
		3: create subscriber
		4: subscribe to topic
		5: get messages
	
		q: QUIT`)
		reqtype := readEntry("choose an option (1, 2, 3, q, etc...)")

		switch reqtype {
		case "1":
			topic := readEntry("enter topic name")
			request = rawRequest{
				Query:   map[string]interface{}{"name": topic},
				API:     apicreateTopic,
				Payload: nil,
			}
		case "2":
			topic := readEntry("enter topic name")
			body := readEntry("enter body content")

			request = rawRequest{
				Query:   map[string]interface{}{"topic": topic, "body": body},
				API:     apipublishMessage,
				Payload: nil,
			}
		case "3":
			request = rawRequest{
				Query:   nil,
				API:     apicreateSubscriber,
				Payload: nil,
			}
		case "4":
			subscriber := readEntry("enter subscriber id")
			topic := readEntry("enter topic name")

			request = rawRequest{
				Query:   map[string]interface{}{"subid": subscriber, "topic": topic},
				API:     apisubscribeToTopic,
				Payload: nil,
			}
		case "5":
			subid := readEntry("enter subscriber id")
			topic := readEntry("enter topic name")

			request = rawRequest{
				API:   apigetMessages,
				Query: map[string]interface{}{"topic": topic, "subid": subid},
			}
		case "q":
			conn.Close()
			fmt.Printf("quitting\n")
			os.Exit(0)
		default:
			fmt.Printf("\n\ninvalid option\n\n")
			continue
		}

		switch protocol {
		case tcpnetwork:
			conn, err = net.Dial(tcpnetwork, tcpaddr)
		case udpnetwork:

			conn, err = net.DialUDP(udpnetwork, laddr, raddr)
		}
		if err != nil {
			panic(err)
		}

		if reqtype == "5" {
			udpconn, ok := conn.(*net.UDPConn)
			if !ok {
				tcpSubscriberGetMessages(conn, request)
			} else {
				count, err := udpGetMessageCount(udpconn, request)
				if err != nil {
					panic(err)
				}

				for i := 0; i < count; i++ {
					message, err := udpGetMessage(udpconn, request)
					if err != nil {
						panic(err)
					}

					log.Print("downloaded message: ", message.ID)
					log.Print("working on message: ", message.ID)
					// simulate work
					time.Sleep(3 * time.Second)

					log.Print("marking message as done")
					response, err := udpDoneMessage(udpconn, request, message)
					if err != nil {
						panic(err)
					}
					fmt.Printf("\n***response***: \n%+v\n\n", response)
				}

				conn.Close()
			}
		} else {
			serverConn(request, conn)
		}
	}
}

func udpDoneMessage(
	conn *net.UDPConn,
	rawRequest rawRequest,
	m message.Message,
) (jsonResponse, error) {
	var (
		out []byte
		err error
		buf []byte
		n   int
		in  jsonResponse
	)
	rawRequest.API = apiudpCompleteMessage
	m.Done = true

	rawRequest.Payload = map[string]interface{}{
		"id":    m.ID,
		"topic": m.Topic,
		"body":  m.Body,
		"done":  m.Done,
	}

	out, err = json.Marshal(rawRequest)
	if err == nil {
		_, err = conn.Write(out)
	}
	if err == nil {
		buf = make([]byte, 1024)
		n, err = conn.Read(buf)
	}
	if err == nil {
		err = json.Unmarshal(buf[:n], &in)
	}

	return in, err
}

func udpGetMessageCount(conn *net.UDPConn, rawRequest rawRequest) (int, error) {
	var (
		out []byte
		err error
		buf []byte
		n   int
		in  struct {
			Err    string `json:"err"`
			Status int    `json:"status"`
			Data   struct {
				Count int `json:"count"`
			}
		}
	)
	rawRequest.API = apiudpGetMessageCount
	out, err = json.Marshal(rawRequest)
	if err == nil {
		_, err = conn.Write(out)
	}
	if err == nil {
		buf = make([]byte, 1024)
		n, err = conn.Read(buf)
	}
	if err == nil {
		err = json.Unmarshal(buf[:n], &in)
	}

	return in.Data.Count, err
}

func udpGetMessage(conn *net.UDPConn, rawRequest rawRequest) (message.Message, error) {
	var (
		out []byte
		err error
		buf []byte
		n   int
		in  struct {
			Err    string `json:"err"`
			Status int    `json:"status"`
			Data   struct {
				Message message.Message `json:"message"`
			}
		}
	)
	rawRequest.API = apiudpGetNextMessage
	out, err = json.Marshal(rawRequest)
	if err == nil {
		_, err = conn.Write(out)
	}
	if err == nil {
		buf = make([]byte, 1024)
		n, err = conn.Read(buf)
	}
	if err == nil {
		err = json.Unmarshal(buf[:n], &in)
	}

	return in.Data.Message, err
}

func tcpSubscriberGetMessages(conn net.Conn, request rawRequest) {
	var (
		response jsonResponse
		err      error
		in       struct {
			Err    string
			Status int
			Data   struct {
				Message message.Message `json:"message"`
				Count   int             `json:"count"`
			}
		}
		inm    message.Message
		c      int
		out    jsonResponse
		n      int
		buf    []byte
		reqout []byte
		resout []byte
	)

	// encode request
	reqout, err = json.Marshal(request)

	// send jsonrequest with method,api,subid,topicid
	if err == nil {
		_, err = conn.Write(reqout)
	}

	for {
		// accept reply, save to buffer
		if err == nil {
			buf = make([]byte, 1024)
			n, err = conn.Read(buf)
		}
		// unmarshal to jsonResponse, check status and finish if server validation fails
		if err == nil {
			var jr jsonResponse
			err = json.Unmarshal(buf[:n], &jr)
			if jr.Status != http.StatusOK {
				response = jr
				break
			}
		}
		// if ok unmarshal saved buffer to correct type
		if err == nil {
			err = json.Unmarshal(buf[:n], &in)
			c = in.Data.Count
			inm = in.Data.Message
			log.Println("downloading message " + in.Data.Message.ID)
		}
		// simulate work, set msg done flag, reduce count by 1
		if err == nil {
			log.Println("simulate working on message " + in.Data.Message.ID)
			time.Sleep(3 * time.Second)
			inm.Done = true
			c = c - 1
		}
		// create response
		if err == nil {
			out.Status = http.StatusOK
			out.Data = struct {
				Message message.Message `json:"message"`
			}{
				Message: inm,
			}
			out.Err = ""
		}
		// encode request
		if err == nil {
			resout, err = json.Marshal(out)
		}
		// send response
		if err == nil {
			_, err = conn.Write(resout)
		}
		// if count zero exit loop
		if c == 0 {
			response = jsonResponse{
				Err:    "",
				Status: http.StatusOK,
				Data:   nil,
			}
			break
		}
		// if count not zero await next message
	}

	conn.Close()
	fmt.Printf("\n***response***: \n%+v\n\n", response)
	fmt.Print("Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func serverConn(request rawRequest, conn net.Conn) {
	var (
		reqout   []byte
		err      error
		n        int
		buf      []byte
		udpconn  *net.UDPConn
		ok       bool
		response jsonResponse
	)
	// encode request
	reqout, err = json.Marshal(request)

	// send jsonrequest with method,api,subid,topicid

	if err == nil {
		udpconn, ok = conn.(*net.UDPConn)
		if !ok {
			err = json.NewEncoder(conn).Encode(request)
		} else {
			_, err = udpconn.Write(reqout)
		}
	}

	buf = make([]byte, 1024)
	if err == nil {
		udpconn, ok = conn.(*net.UDPConn)
		if !ok {
			n, err = conn.Read(buf)
		} else {
			n, err = udpconn.Read(buf)
		}
	}

	if err != nil {
		response = jsonResponse{
			Err:    err.Error(),
			Status: http.StatusInternalServerError,
			Data:   nil,
		}
	}

	json.Unmarshal(buf[:n], &response)

	conn.Close()
	fmt.Printf("\n***response***: \n%+v\n\n", response)
	fmt.Print("Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func readEntry(instruction string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s:\n", instruction)
	fmt.Print("-> ")
	entry, _ := reader.ReadString('\n')
	return strings.Replace(entry, "\n", "", -1)
}
