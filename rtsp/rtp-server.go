package rtsp_server

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

/**
0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|V=2|P|X|  CC   |M|     PT      |       sequence number         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           timestamp                           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           synchronization source (SSRC) identifier            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|            contributing source (CSRC) identifiers             |
|                             ....                              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                                                               |
|                         RTP Payload                           |
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/

const UDP_BUF_SIZE = 1048576
const rtpFile = "video_output.rtp"

type RTPPack struct {
	*RTPHeader
	RTPPayload []byte
}

type RTPHeader struct {
	Version   uint8 //版本
	Padding   bool  //填充(p):1比特，如果设置，表示数据末尾有填充字节
	Extension bool  //扩展（X）：1比特，如果设置，表示有扩展头部
	CSRCCount uint8 //贡献源计数（CC）：4比特，表示贡献源标识符的数量
	Marker    bool  //标志位（M）：1比特，用于标记重要事件（例如视频帧的边界）
	/**
	静态负载类型
	0: PCMU (音频，8 kHz, 单声道)
	3: GSM (音频, 8 kHz)
	4: G723 (音频, 8 kHz)
	8: PCMA (音频, 8 kHz, 单声道)
	9: G722 (音频, 8 kHz)
	10: L16 (音频, 44.1 kHz, 立体声)
	11: L16 (音频, 44.1 kHz, 单声道)
	14: MPA (音频, 90 kHz)
	26: JPEG (视频)
	31: H.261 (视频)
	32: MPV (视频)
	33: MP2T (视频)
	34: H.263 (视频)

	动态负载类型
	96: 通常用于H.264视频
	97: 通常用于AAC音频
	98: 通常用于 VP8 视频
	99: 通常用于 MP2 音频
	100: 通常用于 MP4 音频
	101: 通常用于 H.264 音频/视频的前向错误修复（FEC）
	102: 通常用于 H.264 SVC（可伸缩视频编码）
	103 被映射为 H.265 编码
	104 被映射为 H.266 编码
	*/
	PayloadType    uint8    //负载类型（PT）：7比特，表示负载的类型（例如音频、视频）
	SequenceNumber uint16   //序列号：16比特，每发送一个RTP包，序列号加1。
	Timestamp      uint32   //时间戳：32比特，表示这个RTP包的采样时刻
	SSRC           uint32   //同步源标识符（SSRC）：32比特，唯一标识RTP流
	CSRC           []uint32 //贡献源标识符（CSRC）：0到15个32比特，用于混音器识别参与流的源
}

func CreatePacket(packet []byte) *RTPPack {
	var rtpPack = new(RTPPack)
	var rtpHeader = new(RTPHeader)
	rtpPack.RTPHeader = rtpHeader

	rtpPack.Version = packet[0] >> 6
	rtpPack.Padding = (packet[0]>>5)&1 == 1
	rtpPack.Extension = (packet[0]>>4)&1 == 1
	rtpPack.CSRCCount = packet[0] & 0x0F
	rtpPack.Marker = (packet[1]>>7)&1 == 1
	rtpPack.PayloadType = packet[1] & 0x7F

	rtpPack.SequenceNumber = binary.BigEndian.Uint16(packet[2:4])
	rtpPack.Timestamp = binary.BigEndian.Uint32(packet[4:8])
	rtpPack.SSRC = binary.BigEndian.Uint32(packet[8:12])
	rtpPack.CSRC = make([]uint32, rtpPack.CSRCCount)

	for i := 0; i < int(rtpPack.CSRCCount); i++ {
		rtpPack.CSRC[i] = binary.BigEndian.Uint32(packet[12+4*i : 16+4*i])
	}

	fmt.Printf("RTP Packet - Seq: %d, Timestamp: %d, SSRC: %d, CSRC: %v\n",
		rtpPack.SequenceNumber, rtpPack.Timestamp, rtpPack.SSRC, rtpPack.CSRC)

	payloadOffset := 12 + 4*rtpPack.CSRCCount
	rtpPack.RTPPayload = packet[payloadOffset:]
	return rtpPack
}

/*
  *
 *开音频通信端口
*/
func StartAudio() (err error) {

	return nil
}

func StartVideo() (port int, err error) {
	addr := net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 0}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening for UDP packets:", err)
		return 0, err
	}
	defer conn.Close()

	localAdd := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println("Listening on UDP port", localAdd.Port)

	//go readRtByUDP(conn)
	return addr.Port, nil
}

func readRtByUDP(conn *net.UDPConn) {
	bufUDP := make([]byte, UDP_BUF_SIZE)

	var n = -1
	var err error
	//TODO 如果rtsp关闭这个也需要关闭
	for true {
		if n, _, err = conn.ReadFromUDP(bufUDP); err != nil {
			fmt.Println("Error reading RTP packet:", err)
			continue
		}
		rtpBytes := make([]byte, n)
		copy(rtpBytes, bufUDP)

		rtpPack := CreatePacket(rtpBytes)
		HandleRTP(rtpPack)
	}
}

func HandleRTP(rtpPack *RTPPack) {
	file, err := os.Create(rtpFile)
	if err != nil {
		fmt.Println("file open error", err)
		return
	}
	for {
		file.Write(rtpPack.RTPPayload)
	}
}
