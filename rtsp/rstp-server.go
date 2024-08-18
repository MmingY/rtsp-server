package rtsp_server

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

const (
	rtspPort = 8554
)

type RtspSession struct {
	sessionId    string
	RtpVideoPort int
	RtpAudioPort int
	StreamID     string        // 流ID
	Timeout      time.Duration // 会话超时时间
}

// RTPInfo 代表RTP信息
type RTPInfo struct {
	SSRC      uint32 // RTP同步源标识符
	PT        uint8  // RTP有效负载类型
	Payload   []byte // RTP负载数据
	Timestamp uint32 // RTP时间戳
	Sequence  uint16 // RTP序列号
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	id, err := GenerateSessionID()
	if err != nil {
		fmt.Println("generate session ID error", err)
		return
	}
	id = "f2gP-LrSR"
	rtspSession := &RtspSession{
		sessionId: id,
	}
	reader := bufio.NewReader(conn)
	buf := make([]byte, 1)
	rtpLen := make([]byte, 2)
	for {

		if _, err = io.ReadFull(conn, buf); err != nil {
			fmt.Println("RTP error:", err)
			return
		}

		if buf[0] == 0x24 { //RTP
			if _, err = io.ReadFull(conn, buf); err != nil {
				fmt.Println("RTP error:", err)
				return
			}
			channelId := int(buf[0])
			handleRTP(conn, &rtpLen, channelId)
		} else { //RTSP
			// OPTIONS rtsp://localhost:8554/test RTSP/1.0
			var method string   //OPTIONS
			var uri string      //rtsp://localhost:8554/test
			var protocol string //协议和版本RTSP/1.0

			//parse the RTSP
			reqBuf := bytes.NewBuffer(nil)
			reqBuf.Write(buf)

			var request *Request
			for {
				//parse headers
				line, isPrefix, err := reader.ReadLine()
				if err != nil {
					fmt.Println("Error reading header", err)
					return
				}
				reqBuf.Write(line)
				if !isPrefix {
					reqBuf.WriteString("\r\n")
				}

				if len(line) == 0 {
					request = NewRequest(reqBuf.String())
					fmt.Println("request byte:", reqBuf.Bytes())
					contentLen := request.GetContentLength()
					//处理请求body中的数据
					if contentLen > 0 {
						bodyBuf := make([]byte, contentLen)
						if bodyLen, err := io.ReadFull(conn, bodyBuf); err != nil {
							fmt.Println("Error reading body:", err)
						} else if bodyLen != contentLen {
							fmt.Println("read rtsp request body failed, expect size[%d], got size[%d]", contentLen, bodyLen)
						}
						request.Body = string(bodyBuf)
					}
					//headers[request] = headerParts[1]
					protocol = request.Version
					uri = request.URL
					break
				}

			}

			switch request.Method {
			case "OPTIONS":
				fmt.Println(" Options:", request)
				handleOptions(conn, protocol)
			case "ANNOUNCE":
				fmt.Println(" ANNOUNCE:", request)
				handleAnnounce(conn, rtspSession)
			case "DESCRIBE":
				fmt.Println(" DESCRIBE:", request)
				handleDescribe(conn, protocol, uri)
			case "SETUP":
				fmt.Println(" SETUP:", request)
				handleSetup(conn, rtspSession, request)
			case "RECORD":
				fmt.Println(" RECORD:", request)
				handleRecord(conn, rtspSession)
			case "TEARDOWN":
				handleTeardown(conn, rtspSession)
			default:
				fmt.Println("Unhandled method:", method)
			}
		}
	}
}

/*
*
处理RTP数据
*/
func handleRTP(conn net.Conn, rtspLenB *[]byte, channelId int) {
	if _, err := io.ReadFull(conn, *rtspLenB); err != nil {
		fmt.Println("rtp read len error:", err)
		return
	}
	rtspLen := int(binary.BigEndian.Uint16(*rtspLenB))
	rtpBuff := make([]byte, rtspLen)
	//channelId 判断通道类型：音频通道，视频通道，音频控制通道，视频控制通道
	if _, err := io.ReadFull(conn, rtpBuff); err != nil {
		fmt.Println("Read RTCP frame error:", err)
		return
	}
	fmt.Println("channel id:%s ,rtp len:%s \r\n", channelId, rtspLen)
	fmt.Println(rtpBuff[:20])
}

func handleOptions(conn net.Conn, protocol string) {
	respones := fmt.Sprintf(" %s 200 OK\r\n", protocol)
	respones += "CSeq:1\r\n"
	respones += "Public: OPTIONS, ANNOUNCE, SETUP, RECORD, TEARDOWN\r\n"
	respones += "\r\n"
	fmt.Println("handleOptions:", respones)
	conn.Write([]byte(respones))
}

func handleAnnounce(conn net.Conn, session *RtspSession) {
	video, err := StartVideo()
	if err != nil {
		fmt.Println("open tfp server error:", err)
		return
	}
	session.RtpVideoPort = video
	response := "RTSP/1.0 200 OK\r\n" +
		"CSeq: 2\r\n" +
		"\r\n"
	conn.Write([]byte(response))
}

func handleDescribe(conn net.Conn, protocol string, uri string) {
	sdp := "v=0\r\n" +
		"o=- 0 0 IN IP4 127.0.0.1\r\n" +
		"s=No Name\r\n" +
		"c=IN IP4 127.0.0.1\r\n" +
		"t=0 0\r\n" +
		"a=tool:libavformat 58.29.100\r\n" +
		"m=video 0 RTP/AVP 96\r\n" +
		"a=rtspmap:96 H264/90000\r\n"

	response := fmt.Sprintf("%s 200 OK\r\n", protocol)
	response += "CSeq:1\r\n"
	response += "Content-Base: " + uri + "/\r\n"
	response += "Content-Type: appliation/\r\n"
	response += fmt.Sprintf("Content-length: %d\r\n", len(sdp))
	response += "\r\n"
	response += sdp
	// response :=sdp
	conn.Write([]byte(response))
}

func handleSetup(conn net.Conn, session *RtspSession, request *Request) {
	if strings.Contains(request.Protocol, "TCP") { //是否使用TCP传输RTP数据
		/*response := "RTSP/1.0 200 OK\r\n" +
			"CSeq: 3\r\n" +
			"Transport: RTP/AVP;unicast;interleaved=" + request.Interleaved + ";mode=record\r\n" +
			"Session: " + session.sessionId + "\r\n" +
			"\r\n"
		write, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Errorf("setup error:", err)
			return
		} else {
			fmt.Println("write len:", write)
		}*/
		response := "RTSP/1.0 200 OK\nServer: EasyDarwin/7.3 (Build/17.0325; Platform/Win32; Release/EasyDarwin; State/Development; )\nCseq: 3\nCache-Control: no-cache\nSession: 132169028622239\nDate: Tue, 13 Nov 2018 02:49:48 GMT\nExpires: Tue, 13 Nov 2018 02:49:48 GMT\nTransport: RTP/AVP/TCP;unicast;mode=record;interleaved=0-1" +
			"\r\n" +
			"\r\n"
		conn.Write([]byte(response))
	} else { //uYT4wxdT:mhpHLE
		response := "RTSP/1.0 200 OK\r\n" +
			"CSeq: 3\r\n" +
			//"Transport: RTP/AVP;unicast;client_port=" + request.ClientPort + ";server_port=9000-9001" + strconv.Itoa(session.RtpVideoPort) + "\r\n" +
			"Transport: RTP/AVP/UDP;unicast;client_port=" + request.ClientPort + ";server_port=9000-9001\r\n" +
			"Session: " + session.sessionId + "\r\n" +
			"\r\n"
		conn.Write([]byte(response))
	}

}

func handleRecord(conn net.Conn, session *RtspSession) {
	response := "RTSP/1.0 200 OK\r\n" +
		"CSeq: 4\r\n" +
		"Session: " + session.sessionId + "\r\n" +
		"\r\n"
	conn.Write([]byte(response))
}

func handleTeardown(conn net.Conn, session *RtspSession) {
	response := "RTSP/1.0 200 OK\r\n" +
		"CSeq: 5\r\n" +
		"Session: " + session.sessionId + "\r\n" +
		"\r\n"
	conn.Write([]byte(response))
	conn.Close()
}

func ListenRTSP() {
	linster, err := net.Listen("tcp", fmt.Sprintf(":%d", 8554))
	if err != nil {
		fmt.Println("Error starting RTSP server:", err)
		return
	}
	defer linster.Close()
	fmt.Println("RTSP server is running at rtsp:localhost:%d/\n", 8554)

	for {
		conn, err := linster.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go HandleConnection(conn)
	}
}

// GenerateSessionID 生成一个随机的 Session ID
func GenerateSessionID() (string, error) {
	// 创建一个 16 字节的数组
	b := make([]byte, 16)
	// 使用 crypto/rand 包生成随机字节
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	// 将字节数组编码为十六进制字符串
	return hex.EncodeToString(b), nil
}
