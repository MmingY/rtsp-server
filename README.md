## 一.RTSP 推流流程

推流是指将媒体内容从客户端发送到服务器。通常由流媒体服务器来处理推流请求。以下是推流的基本步骤：

### 1.建立连接：
客户端（推流端）向服务器发送 OPTIONS 请求以确定服务器支持的请求方法。
```
OPTIONS rtsp://example.com/media.sdp RTSP/1.0
CSeq: 1
User-Agent: Lavf60.3.100
```
服务器返回支持的方法列表。
```
RTSP/1.0 200 OK
CSeq: 1
Session: f2gP-LrSR
Public: DESCRIBE, SETUP, TEARDOWN, PLAY, PAUSE, OPTIONS, ANNOUNCE, RECORD

```
### 2.发送描述：
客户端发送 ANNOUNCE 请求，将媒体描述（如 SDP 描述）发送到服务器。
```
ANNOUNCE rtsp://example.com/media.sdp RTSP/1.0
Content-Type: application/sdp
CSeq: 2
User-Agent: Lavf60.3.100
Session: f2gP-LrSR
Content-Length: 494

v=0
o=- 0 0 IN IP4 127.0.0.1
s=No Name
c=IN IP4 127.0.0.1
t=0 0
a=tool:libavformat 60.3.100
m=video 0 RTP/AVP 96
b=AS:3461
a=rtpmap:96 H264/90000
a=fmtp:96 packetization-mode=1; sprop-parameter-sets=Z01AKZZWA8ARPyykBAQFAAADA+kAAOpghA==,aMqNSA==; profile-level-id=4D4029
a=control:streamid=0
m=audio 0 RTP/AVP 97
b=AS:189
a=rtpmap:97 MPEG4-GENERIC/48000/2
a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3; config=1190
a=control:streamid=1
```
服务器确认接收 ANNOUNCE 请求。
```
RTSP/1.0 200 OK
CSeq: 2
Session: f2gP-LrSR
```
### 3.设置传输参数：
客户端发送 SETUP 请求，指定传输的详细信息（如传输协议、端口等）。
```
SETUP rtsp://127.0.0.1:554/stream/streamid=0 RTSP/1.0
Transport: RTP/AVP/TCP;unicast;interleaved=0-1;mode=record
CSeq: 3
User-Agent: Lavf60.3.100
Session: f2gP-LrSR
```
### 4.服务器返回传输参数的确认信息。
```
RTSP/1.0 200 OK
CSeq: 3
Session: f2gP-LrSR
Transport: RTP/AVP/TCP;unicast;interleaved=0-1;mode=record
```
### 5.设置传输参数：
客户端发送 SETUP 请求，指定传输的详细信息（如传输协议、端口等）
```
SETUP rtsp://127.0.0.1:554/stream/streamid=1 RTSP/1.0
Transport: RTP/AVP/TCP;unicast;interleaved=2-3;mode=record
CSeq: 4
User-Agent: Lavf60.3.100
Session: f2gP-LrSR
```
### 6.服务器返回传输参数的确认信息。
```
RTSP/1.0 200 OK
CSeq: 4
Session: f2gP-LrSR
Transport: RTP/AVP/TCP;unicast;interleaved=2-3;mode=record
```

### 7.开始推流：
客户端发送 RECORD 请求开始推流。
``` 
RECORD rtsp://127.0.0.1:554/stream RTSP/1.0
Range: npt=0.000-
CSeq: 5
User-Agent: Lavf60.3.100
Session: f2gP-LrSR
```
### 8.服务器返回确认信息，表示可以开始传输媒体数据。
``` 
RTSP/1.0 200 OK
CSeq: 5
Session: f2gP-LrSR
```
客户端通过 RTP/RTCP 传输媒体流到服务器。

### 5.结束推流：

客户端发送 TEARDOWN 请求终止会话。
``` 
TEARDOWN rtsp://example.com/media.sdp RTSP/1.0
CSeq: 6
Session: 12345678
```
服务器确认终止请求，关闭会话。
``` 
RTSP/1.0 200 OK
CSeq: 6
Session: 12345678
```

## 二。RTSP 拉流流程
拉流是指从服务器获取媒体流并播放。以下是拉流的基本步骤：

### 1.建立连接：

客户端（拉流端）向服务器发送 OPTIONS 请求以确定服务器支持的请求方法。
服务器返回支持的方法列表。

### 2.获取描述：
客户端发送 DESCRIBE 请求获取媒体资源的描述信息（如 SDP 描述）。
服务器返回媒体描述信息。

### 3.设置传输参数：
客户端发送 SETUP 请求，指定传输的详细信息（如传输协议、端口等）。
```
SETUP rtsp://example.com/media.sdp/streamid=0 RTSP/1.0
CSeq: 3
Transport: RTP/AVP;unicast;client_port=8000-8001
```
服务器返回传输参数的确认信息。
```
RTSP/1.0 200 OK
CSeq: 3
Transport: RTP/AVP;unicast;client_port=8000-8001;server_port=9000-9001
Session: 12345678
```
### 4.开始播放：
客户端发送 PLAY 请求开始接收和播放流媒体。
``` 
PLAY rtsp://example.com/media.sdp RTSP/1.0
CSeq: 4
Session: 12345678
Range: npt=0.000-
```
服务器返回确认信息，开始传输媒体数据。
``` 
RTSP/1.0 200 OK
CSeq: 4
Session: 12345678
RTP-Info: url=rtsp://example.com/media.sdp/streamid=0;seq=9810092;rtptime=3450012
```
服务器通过 RTP/RTCP 将媒体流传输到客户端。

### 5.暂停播放（可选）：
客户端发送 PAUSE 请求暂停播放。
服务器返回确认信息，暂停传输媒体数据。

### 6.停止播放：
客户端发送 TEARDOWN 请求终止会话。
服务器确认终止请求，关闭会话。

## 三.字段含义总结
CSeq：序列号，用于标识请求顺序。

Public：服务器支持的方法列表。

Content-Base：描述信息的基URI。

Content-Type：描述信息的格式。

Content-Length：描述信息的长度。

Transport：传输参数，包含协议、传输方式、端口等。

Session：会话标识符，用于关联后续请求。

Range：播放范围。

RTP-Info：RTP流的信息。
