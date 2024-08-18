## 一.RTSP 推流流程

推流是指将媒体内容从客户端发送到服务器。通常由流媒体服务器来处理推流请求。以下是推流的基本步骤：

### 1.建立连接：
客户端（推流端）向服务器发送 OPTIONS 请求以确定服务器支持的请求方法。
```
OPTIONS rtsp://example.com/media.sdp RTSP/1.0
CSeq: 1
```
服务器返回支持的方法列表。
```
RTSP/1.0 200 OK
CSeq: 1
Public: DESCRIBE, SETUP, TEARDOWN, PLAY, PAUSE, OPTIONS, ANNOUNCE, RECORD

```
### 2.发送描述：
客户端发送 ANNOUNCE 请求，将媒体描述（如 SDP 描述）发送到服务器。
```
DESCRIBE rtsp://example.com/media.sdp RTSP/1.0
CSeq: 2
Accept: application/sdp
```
服务器确认接收 ANNOUNCE 请求。
```
RTSP/1.0 200 OK
CSeq: 2
Content-Base: rtsp://example.com/media.sdp
Content-Type: application/sdp
Content-Length: 460

v=0
o=- 2890844526 2890842807 IN IP4 127.0.0.1
s=RTSP Session
m=video 0 RTP/AVP 96
a=control:streamid=0
m=audio 0 RTP/AVP 97
a=control:streamid=1
```
### 3.设置传输参数：
客户端发送 SETUP 请求，指定传输的详细信息（如传输协议、端口等）。

服务器返回传输参数的确认信息。

### 4.开始推流：
客户端发送 RECORD 请求开始推流。
``` 
RECORD rtsp://example.com/media.sdp RTSP/1.0
CSeq: 5
Session: 12345678
```
服务器返回确认信息，表示可以开始传输媒体数据。
``` 
RTSP/1.0 200 OK
CSeq: 5
Session: 12345678
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

## 四。RTP
静态负载类型
定义: 静态负载类型由IANA（Internet Assigned Numbers Authority）预定义，且在所有RTP会话中具有固定的含义。
范围: 通常在0到95之间。
使用: 不需要在SDP中显式声明，因为其含义是标准化的、众所周知的。
动态负载类型
定义: 动态负载类型在会话中根据需要分配，可以在96到127之间。其含义由会话中的参与方协商确定。
范围: 96到127之间。
使用: 需要在SDP或其他信令协议中显式声明，以告知接收方如何解释这些负载类型。

PT	编码	说明
0	PCMU	音频，8 kHz, 单声道
1	Reserved	保留
2	Reserved	保留
3	GSM	音频，8 kHz
4	G723	音频，8 kHz
5	DVI4	音频，8 kHz
6	DVI4	音频，16 kHz
7	LPC	音频，8 kHz
8	PCMA	音频，8 kHz
9	G722	音频，8 kHz
10	L16	音频，44.1 kHz, 立体声
11	L16	音频，44.1 kHz, 单声道
12	QCELP	音频
13	CN	舒适噪声（Comfort Noise）
14	MPA	音频，90 kHz
15	G728	音频
16	DVI4	音频，11.025 kHz
17	DVI4	音频，22.05 kHz
18	G729	音频
19	Reserved	保留
20	Unassigned	未分配
21	Unassigned	未分配
22	Unassigned	未分配
23	Unassigned	未分配
24	Unassigned	未分配
25	CelB	视频
26	JPEG	视频
27	Unassigned	未分配
28	nv	视频
29	Unassigned	未分配
30	Unassigned	未分配
31	H261	视频
32	MPV	视频
33	MP2T	视频
34	H263	视频
35	Unassigned	未分配
36	Unassigned	未分配
37	Unassigned	未分配
38	Unassigned	未分配
39	Unassigned	未分配
40	Unassigned	未分配
41	Unassigned	未分配
42	Unassigned	未分配
43	Unassigned	未分配
44	Unassigned	未分配
45	Unassigned	未分配
46	Unassigned	未分配
47	Unassigned	未分配
48	Unassigned	未分配
49	Unassigned	未分配
50	Unassigned	未分配
51	Unassigned	未分配
52	Unassigned	未分配
53	Unassigned	未分配
54	Unassigned	未分配
55	Unassigned	未分配
56	Unassigned	未分配
57	Unassigned	未分配
58	Unassigned	未分配
59	Unassigned	未分配
60	Unassigned	未分配
61	Unassigned	未分配
62	Unassigned	未分配
63	Unassigned	未分配
64	Unassigned	未分配
65	Unassigned	未分配
66	Unassigned	未分配
67	Unassigned	未分配
68	Unassigned	未分配
69	Unassigned	未分配
70	Unassigned	未分配
71	Unassigned	未分配
72	Unassigned	未分配
73	Unassigned	未分配
74	Unassigned	未分配
75	Unassigned	未分配
76	Unassigned	未分配
77	Unassigned	未分配
78	Unassigned	未分配
79	Unassigned	未分配
80	Unassigned	未分配
81	Unassigned	未分配
82	Unassigned	未分配
83	Unassigned	未分配
84	Unassigned	未分配
85	Unassigned	未分配
86	Unassigned	未分配
87	Unassigned	未分配
88	Unassigned	未分配
89	Unassigned	未分配
90	Unassigned	未分配
91	Unassigned	未分配
92	Unassigned	未分配
93	Unassigned	未分配
94	Unassigned	未分配
95	Unassigned	未分配