package main

import(
	"fmt"
	"net"
	"strings"
	"bufio"
)

const (
	rtspPort=8554
)

func HandleConnection(conn net.Conn){
	defer conn.Close()

	// reader := bufio.NewReader(conn)
	for {

		//OPTIONS rtsp://example.com/media.sdp RTSP/1.0
		//CSeq: 1


		requestLine ,err :=reader.ReadString('\n')
		if err !=nil{
			fmt.Println("Error reading requset:",err)
			return
		}
		fmt.Println("requestLine:"+requestLine)
		// OPTIONS rtsp://localhost:8554/test RTSP/1.0
		requestLine=strings.TrimSpace(requestLine)

		//parse the RTSP

		parts := strings.Split(requestLine," ")
		if len(parts)<3{
			fmt.Println("Invalid request line:",requestLine)
			return
		}

		method := parts[0]//OPTIONS
		uri := parts[1]//rtsp://localhost:8554/test
		protocol := parts[2]//协议和版本RTSP/1.0

		//parse headers
		headers := make(map[string]string)

		for{
			line ,err := reader.ReadString('\n')
			if err !=nil{
				fmt.Println("Error reading header",err)
				return
			}

			line =strings.TrimSpace(line)
			if line == "" {
				break
			}

			headerParts := strings.SplitN(line,": ",2)
			if len(headerParts) !=2{
				continue
			}
			headers[headerParts[0]]=headerParts[1]
		}

		switch method {
		case "OPTIONS":
			fmt.Println(" Options:",requestLine)
			handleOptions(conn,protocol)
		case "ANNOUNCE":
			fmt.Println(" ANNOUNCE:",requestLine)
			handleAnnounce(conn)
		case "DESCRIBE":
			fmt.Println(" DESCRIBE:",requestLine)
			handleDescribe(conn,protocol,uri)
		case "SETUP":
			fmt.Println(" SETUP:",requestLine)
			handleSetup(conn)
		case "RECORD":
			fmt.Println(" RECORD:",requestLine)
            handleRecord(conn)
        case "TEARDOWN":
            handleTeardown(conn)
		default:
			fmt.Println("Unhandled method:",method)
		}
	}
}


func handleOptions(conn net.Conn,protocol string){
		respones := fmt.Sprintf(" %s 200 OK\r\n",protocol)
		respones +="CSeq:1\r\n"
		respones +="Public: OPTIONS, ANNOUNCE, SETUP, RECORD, TEARDOWN\r\n"
		respones +="\r\n"
		fmt.Println("handleOptions:",respones)
		conn.Write([]byte(respones))
}

func  handleAnnounce(conn net.Conn){
	response := "RTSP/1.0 200 OK\r\n" +
        "CSeq: 2\r\n" +
        "\r\n"
	conn.Write([]byte(response))
}

func handleDescribe(conn net.Conn,protocol string,uri string){
		sdp :="v=0\r\n"+
		"o=- 0 0 IN IP4 127.0.0.1\r\n"+
		"s=No Name\r\n"+
		"c=IN IP4 127.0.0.1\r\n"+
		"t=0 0\r\n"+
		"a=tool:libavformat 58.29.100\r\n"+
		"m=video 0 RTP/AVP 96\r\n"+
		"a=rtspmap:96 H264/90000\r\n"

		response := fmt.Sprintf("%s 200 OK\r\n",protocol)
		response +="CSeq:1\r\n"
		response +="Content-Base: "+uri+"/\r\n"
		response +="Content-Type: appliation/\r\n"
		response +=fmt.Sprintf("Content-length: %d\r\n",len(sdp))
		response += "\r\n"
		response +=sdp
		// response :=sdp
		conn.Write([]byte(response))
}

func handleSetup(conn net.Conn){
	response :="RTSP/1.0 200 OK\r\n"+
	"CSeq: 3\r\n"+
	"Transport: RTP/AVP;unicast;client_port=8000-8001;server_port=9000-9001\r\n"+
	"Session: 12345678\r\n"
	conn.Write([]byte(response))
}

func handleRecord(conn net.Conn) {
    response := "RTSP/1.0 200 OK\r\n" +
        "CSeq: 4\r\n" +
        "Session: 12345678\r\n" +
        "\r\n"
    conn.Write([]byte(response))
}

func handleTeardown(conn net.Conn) {
    response := "RTSP/1.0 200 OK\r\n" +
        "CSeq: 5\r\n" +
        "Session: 12345678\r\n" +
        "\r\n"
    conn.Write([]byte(response))
    conn.Close()
}

func listenRTSP(){
    linster ,err :=net.Listen("tcp",fmt.Sprintf(":%d",8554))
    if err !=nil{
        fmt.Println("Error starting RTSP server:",err)
        return
    }
    defer linster.Close()
    fmt.Println("RTSP server is running at rtsp:localhost:%d/\n",8554)

    for{
        conn,err :=  linster.Accept()
        if err !=nil{
            fmt.Println("Error accepting connection:",err)
            continue
        }
        go HandleConnection(conn)
    }
}

func main(){
	listenRTSP()
}