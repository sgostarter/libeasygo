package netutils

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/sgostarter/i/commerr"
)

type iCMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

/*
if errors.Is(err, os.ErrPermission) {
			fmt.Println("xx")
		}
*/

var (
	_noPermission4CheckLatency atomic.Bool
)

/*
# mac os $ ping www.baidu.com -c 1
PING www.a.shifen.com (220.181.38.150): 56 data bytes
64 bytes from 220.181.38.150: icmp_seq=0 ttl=53 time=25.318 ms

--- www.a.shifen.com ping statistics ---
1 packets transmitted, 1 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 25.318/25.318/25.318/0.000 ms

// linux ping www.baidu.com -c 1
PING www.a.shifen.com (220.181.38.149) 56(84) bytes of data.
64 bytes from 220.181.38.149 (220.181.38.149): icmp_seq=1 ttl=251 time=3.26 ms

--- www.a.shifen.com ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 3.264/3.264/3.264/0.000 ms
*/

func OSPing(host string) (latency time.Duration, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	countKey := "-c"
	if runtime.GOOS == "windows" {
		countKey = "-n"
	}

	cmd := exec.CommandContext(ctx, "ping", countKey, "1", host)

	d, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	s := string(d)

	ps := strings.Split(s, "\n")
	if len(ps) < 2 {
		err = commerr.ErrUnknown

		return
	}

	idx := strings.Index(ps[1], "time=")
	if idx < 0 {
		err = commerr.ErrUnknown

		return
	}

	ps = strings.Split(ps[1][idx+5:], " ")
	if len(ps) != 2 {
		err = commerr.ErrUnknown

		return
	}

	if ps[1] != "ms" {
		err = commerr.ErrUnknown

		return
	}

	f, err := strconv.ParseFloat(ps[0], 64)
	if err != nil {
		return
	}

	latency = time.Millisecond * time.Duration(f)

	return
}

func Latency(host string) (latency time.Duration, err error) {
	if _noPermission4CheckLatency.Load() {
		return OSPing(host)
	}

	latency, err = LatencyEx(host)
	if err == nil {
		return
	}

	if errors.Is(err, os.ErrPermission) {
		_noPermission4CheckLatency.Store(true)
	}

	return OSPing(host)
}

func LatencyEx(host string) (latency time.Duration, err error) {
	timeout := time.Second * 5

	conn, err := net.DialTimeout("ip:icmp", host, timeout)
	if err != nil {
		return
	}

	defer func() {
		_ = conn.Close()
	}()

	icmp := iCMP{8, 0, 0, 0, 0}

	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.BigEndian, icmp)

	data := make([]byte, 1024)
	buffer.Write(data)
	data = buffer.Bytes()

	icmp.SequenceNum = uint16(1) // 检验和设为0
	data[2], data[3] = byte(0), byte(0)

	data[6], data[7] = byte(icmp.SequenceNum>>8), byte(icmp.SequenceNum)

	icmp.Checksum = checkSum(data)
	data[2], data[3] = byte(icmp.Checksum>>8), byte(icmp.Checksum)

	timeStart := time.Now()
	_ = conn.SetDeadline(timeStart.Add(timeout))

	_, err = conn.Write(data)
	if err != nil {
		return
	}

	buf := make([]byte, 65535)

	_, err = conn.Read(buf)
	if err != nil {
		return
	}

	latency = time.Since(timeStart)

	return
}

func checkSum(data []byte) (rt uint16) {
	var (
		sum    uint32
		length = len(data)
		index  int
	)

	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}

	if length > 0 {
		sum += uint32(data[index]) << 8
	}

	rt = uint16(sum) + uint16(sum>>16)

	return ^rt
}
