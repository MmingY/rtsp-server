## FFmpeg想服务推流

ffmpeg -re -i input.mp4 -c:v libx264 -f rtsp rtsp://localhost:8554/stream

## RTSP客户端推流流程

**客户端：**

```
OPTIONS rtsp://192.168.1.108:554/live.sdp RTSP/1.0
CSeq: 1
User-Agent: Lavf57.71.100
```

**服务端：**

```
RTSP/1.0 200 OK
Server: EasyDarwin/7.3 (Build/17.0325; Platform/Win32; Release/EasyDarwin; State/Development; )
Cseq: 1
Public: DESCRIBE, SETUP, TEARDOWN, PLAY, PAUSE, OPTIONS, ANNOUNCE, RECORD
 ```

**客户端：**

```
ANNOUNCE rtsp://192.168.1.108:554/live.sdp RTSP/1.0
Content-Type: application/sdp
CSeq: 2
User-Agent: Lavf57.71.100
Content-Length: 325

v=0
o=- 0 0 IN IP4 127.0.0.1
s=Media Server
c=IN IP4 192.168.1.108
t=0 0
a=tool:libavformat 57.71.100
m=video 0 RTP/AVP 96
a=rtpmap:96 H264/90000
a=fmtp:96 packetization-mode=1;
sprop-parameter-sets=Z2QAHqw0ygsBJ/wFuCgoKgAAB9AAAYah0MALFAALE9d5caGAFigAFieu8uFA,aO48MA==; profile-level-id=64001E
a=control:streamid=0
```

**服务端：**
```
RTSP/1.0 200 OK
Server: EasyDarwin/7.3 (Build/17.0325; Platform/Win32; Release/EasyDarwin; State/Development; )
Cseq: 2
```
**客户端：**

```
SETUP rtsp://192.168.1.108:554/live.sdp/streamid=0 RTSP/1.0
Transport: RTP/AVP/TCP;unicast;interleaved=0-1;mode=record
CSeq: 3
User-Agent: Lavf57.71.100
```

**服务端：**

```
RTSP/1.0 200 OK
Server: EasyDarwin/7.3 (Build/17.0325; Platform/Win32; Release/EasyDarwin; State/Development; )
Cseq: 3
Cache-Control: no-cache
Session: 132169028622239
Date: Tue, 13 Nov 2018 02:49:48 GMT
Expires: Tue, 13 Nov 2018 02:49:48 GMT
Transport: RTP/AVP/TCP;unicast;mode=record;interleaved=0-1

```

**客户端：**

```
RECORD rtsp://192.168.1.108:554/live.sdp RTSP/1.0
Range: npt=0.000-
CSeq: 4
User-Agent: Lavf57.71.100
Session: 132169028622239
```

**服务端：**

```
RTSP/1.0 200 OK
Server: EasyDarwin/7.3 (Build/17.0325; Platform/Win32; Release/EasyDarwin; State/Development; )
Cseq: 4
Session: 132169028622239
RTP-Info: url=rtsp://192.168.1.108:554/live.sdp/live.sdp
```