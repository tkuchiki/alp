package parsers

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

const (
	conjoinPcapKeyHeader   = "Internal-ALP-Pcap-Conjoin-Key"
	timestampPcapKeyHeader = "Internal-ALP-Pcap-Timestamp-Key"
)

type PcapParser struct {
	queryString    bool
	qsIgnoreValues bool

	resCh chan *http.Response // conjoined *http.Response
}

func NewPcapParser(r io.Reader, rawServerIPs []string, serverPort uint16, query, qsIgnoreValues bool) (Parser, error) {
	h, err := pcapgo.NewReader(r)
	if err != nil {
		return nil, err
	}

	serverIPs := make([]net.IP, len(rawServerIPs))
	for i, rawServerIP := range rawServerIPs {
		serverIP := net.ParseIP(rawServerIP)
		if serverIP == nil {
			return nil, fmt.Errorf("failed to parse ip: %s", rawServerIP)
		}

		serverIPs[i] = serverIP
	}

	reqCh := make(chan *http.Request)
	resCh := make(chan *http.Response)
	go func() {
		ps := gopacket.NewPacketSource(h, h.LinkType())
		sf := newPcapHttpStreamFactory(reqCh, resCh, serverIPs, serverPort)
		sp := tcpassembly.NewStreamPool(sf)
		asmblr := tcpassembly.NewAssembler(sp)
		readAndAssembleAllPackets(ps, asmblr)
		sf.gracefulShutdown()
	}()

	conjoinedResCh := make(chan *http.Response)
	go conjoinRequestAndResponse(reqCh, resCh, conjoinedResCh)

	return &PcapParser{
		queryString:    query,
		qsIgnoreValues: qsIgnoreValues,

		resCh: conjoinedResCh,
	}, nil
}

func (j *PcapParser) Parse() (*ParsedHTTPStat, error) {
	res := <-j.resCh
	if res == nil {
		return nil, io.EOF
	}
	req := res.Request

	reqTimestamp, err := unixNanoStrToTime(req.Header.Get(timestampPcapKeyHeader))
	if err != nil {
		return nil, err
	}

	resTimestamp, err := unixNanoStrToTime(res.Header.Get(timestampPcapKeyHeader))
	if err != nil {
		return nil, err
	}
	resTime := resTimestamp.Sub(reqTimestamp)

	var uri string
	if j.queryString {
		if j.qsIgnoreValues {
			values := req.URL.Query()
			for q := range values {
				values.Set(q, "xxx")
			}
			req.URL.RawQuery = values.Encode()
		}
		uri = req.URL.String()
	} else {
		req.URL.RawQuery = ""
		uri = req.URL.String()
	}

	resBodyBytes := res.ContentLength
	stat := NewParsedHTTPStat(uri, req.Method, reqTimestamp.Format(time.RFC3339), math.Abs(resTime.Seconds()), float64(resBodyBytes), res.StatusCode)
	return stat, nil
}

func (j *PcapParser) ReadBytes() int {
	return 0
}

func (j *PcapParser) SetReadBytes(n int) {
	// not supported
}

func (j *PcapParser) Seek(n int) error {
	return errors.New("not supported")
}

type pcapHttpStreamStat struct {
	waiting atomic.Value
	cond    *sync.Cond

	doingReqCount int64
	doingResCount int64
}

func (s *pcapHttpStreamStat) startReq() {
	atomic.AddInt64(&s.doingReqCount, 1)
}

func (s *pcapHttpStreamStat) completeReq() {
	n := atomic.AddInt64(&s.doingReqCount, -1)
	if n == 0 {
		waiting := s.waiting.Load().(bool)
		if waiting {
			s.cond.Broadcast()
		}
	}
}

func (s *pcapHttpStreamStat) startRes() {
	atomic.AddInt64(&s.doingResCount, 1)
}

func (s *pcapHttpStreamStat) completeRes() {
	n := atomic.AddInt64(&s.doingResCount, -1)
	if n == 0 {
		waiting := s.waiting.Load().(bool)
		if waiting {
			s.cond.Broadcast()
		}
	}
}

func (s *pcapHttpStreamStat) waitForCompleteAll() {
	s.waiting.Store(true)
	s.cond.L.Lock()
	for atomic.LoadInt64(&s.doingReqCount) > 0 || atomic.LoadInt64(&s.doingResCount) > 0 {
		s.cond.Wait()
	}
	s.cond.L.Unlock()
}

type pcapHttpStreamFactory struct {
	reqCh      chan *http.Request
	resCh      chan *http.Response
	serverIPs  []net.IP
	serverPort uint16

	stat pcapHttpStreamStat
}

func newPcapHttpStreamFactory(reqCh chan *http.Request, resCh chan *http.Response, serverIPs []net.IP, serverPort uint16) *pcapHttpStreamFactory {
	f := &pcapHttpStreamFactory{
		reqCh:      reqCh,
		resCh:      resCh,
		serverIPs:  serverIPs,
		serverPort: serverPort,
	}
	f.stat.waiting.Store(false)
	f.stat.cond = sync.NewCond(&sync.Mutex{})
	return f
}

func (h *pcapHttpStreamFactory) New(nf, tf gopacket.Flow) tcpassembly.Stream {
	rs := newTCPReaderStream()

	clientAddr, isReq, unknown := h.detectTrafficDirection(nf, tf)
	if unknown {
		go tcpreader.DiscardBytesToEOF(&rs)
		return &rs
	}

	if isReq {
		go parseHTTPRequest(&rs, clientAddr, h.reqCh, &h.stat)
	} else {
		go parseHTTPResponse(&rs, clientAddr, h.resCh, &h.stat)
	}
	return &rs
}

func (h *pcapHttpStreamFactory) detectTrafficDirection(nf, tf gopacket.Flow) (clientAddr *net.TCPAddr, isReq bool, unknown bool) {
	if nf.EndpointType() != layers.EndpointIPv4 && nf.EndpointType() != layers.EndpointIPv6 {
		unknown = true
		return
	}
	if tf.EndpointType() != layers.EndpointTCPPort {
		unknown = true
		return
	}

	srcIP := net.IP(nf.Src().Raw())
	srcPort := binary.BigEndian.Uint16(tf.Src().Raw())
	dstIP := net.IP(nf.Dst().Raw())
	dstPort := binary.BigEndian.Uint16(tf.Dst().Raw())
	for _, serverIP := range h.serverIPs {
		if srcIP.Equal(serverIP) && srcPort == h.serverPort {
			clientAddr = &net.TCPAddr{
				IP:   dstIP,
				Port: int(dstPort),
			}
			isReq = false
			return
		} else if dstIP.Equal(serverIP) && dstPort == h.serverPort {
			clientAddr = &net.TCPAddr{
				IP:   srcIP,
				Port: int(srcPort),
			}
			isReq = true
			return
		}
	}

	unknown = true
	return
}

func (h *pcapHttpStreamFactory) gracefulShutdown() {
	h.stat.waitForCompleteAll()
	close(h.reqCh)
	close(h.resCh)
}

func readAndAssembleAllPackets(packetSource *gopacket.PacketSource, assembler *tcpassembly.Assembler) {
	defer assembler.FlushAll()
	for {
		p, err := packetSource.NextPacket()
		if err == io.EOF {
			return
		} else if err != nil {
			log.Printf("Failed to read packet: %v", err)
			return
		}

		if p.NetworkLayer() == nil || p.TransportLayer() == nil || p.TransportLayer().LayerType() != layers.LayerTypeTCP {
			log.Printf("Unusable packet: %v", p)
			continue
		}

		tcp := p.TransportLayer().(*layers.TCP)
		assembler.AssembleWithTimestamp(p.NetworkLayer().NetworkFlow(), tcp, p.Metadata().Timestamp)
	}
}

func conjoinRequestAndResponse(reqCh chan *http.Request, resCh chan *http.Response, conjoinedResCh chan *http.Response) {
	var copyBuf [4096]byte
	var bb bytes.Buffer

	var onReq func(req *http.Request)
	var onRes func(res *http.Response)

	reqBufMap := map[string]*http.Request{}
	onReq = func(req *http.Request) {
		if req == nil {
			onReq = nil
			return
		}

		key := req.Header.Get(conjoinPcapKeyHeader)
		req.Header.Del(conjoinPcapKeyHeader)
		reqBufMap[key] = req
	}
	onRes = func(res *http.Response) {
		if res == nil {
			onRes = nil
			return
		}

		key := res.Header.Get(conjoinPcapKeyHeader)
		if req, ok := reqBufMap[key]; ok {
			res.Request = req
			res.Header.Del(conjoinPcapKeyHeader)
			delete(reqBufMap, key)

			// re-parse with request for chunked response
			if res.ContentLength == -1 {
				_ = res.Write(&bb)
				defer bb.Reset()
				br := bufio.NewReader(io.MultiReader(&bb, res.Body))

				var err error
				res, err = http.ReadResponse(br, req)
				if err != nil {
					clientAddr := key
					timestamp, _ := unixNanoStrToTime(res.Header.Get(timestampPcapKeyHeader))
					log.Printf("Failed to read HTTP request from the client %v at %s: %v", clientAddr, timestamp.Format(time.RFC3339Nano), err)
					return
				}

				// discard HTTP body
				n, err := io.CopyBuffer(io.Discard, res.Body, copyBuf[:])
				if err != nil {
					clientAddr := key
					timestamp, _ := unixNanoStrToTime(res.Header.Get(timestampPcapKeyHeader))
					log.Printf("Failed to read HTTP body to the client %v at %s: %v", clientAddr, timestamp.Format(time.RFC3339Nano), err)
					return
				}
				res.ContentLength = n // real response length
				_ = res.Body.Close()
			}

			conjoinedResCh <- res
		} else {
			// ignore it
		}
	}

	for {
		if onReq != nil && onRes != nil {
			select {
			case req := <-reqCh:
				onReq(req)
			case res := <-resCh:
				onRes(res)
			}
		} else if onReq != nil {
			req := <-reqCh
			onReq(req)
		} else if onRes != nil {
			res := <-resCh
			onRes(res)
		} else {
			close(conjoinedResCh)
			return
		}
	}
}

type tcpReaderStream struct {
	tcpreader.ReaderStream
	timestamps []time.Time
}

func newTCPReaderStream() tcpReaderStream {
	return tcpReaderStream{
		ReaderStream: tcpreader.NewReaderStream(),
	}
}

func (s *tcpReaderStream) Reassembled(rs []tcpassembly.Reassembly) {
	for _, r := range rs {
		if len(r.Bytes) == 0 {
			continue
		}

		// HTTP messages aren't started in the middle of a packet generally.
		if bytes.HasPrefix(r.Bytes, []byte("HTTP/1.")) {
			s.timestamps = append(s.timestamps, r.Seen)
		} else if len(r.Bytes) > len("GET / HTTP/1.0\r\n") {
			nr := bytes.Index(r.Bytes, []byte{'\r', '\n'})
			if nr == -1 {
				continue // not http request
			}
			if !bytes.Contains(r.Bytes[:nr], []byte("HTTP/1.")) {
				continue // not http request
			}

			s.timestamps = append(s.timestamps, r.Seen)
		}
	}
	s.ReaderStream.Reassembled(rs)
}

func (s *tcpReaderStream) consumeTimestamp() (timestamp time.Time) {
	if len(s.timestamps) == 0 {
		return // zero time.Time
	}

	timestamp = s.timestamps[0]
	s.timestamps = s.timestamps[1:]
	return
}

func parseHTTPRequest(rs *tcpReaderStream, clientAddr *net.TCPAddr, reqCh chan *http.Request, stat *pcapHttpStreamStat) {
	stat.startReq()
	defer stat.completeReq()

	var copyBuf [4096]byte
	bufr := bufio.NewReader(rs)
	for {
		// check EOF
		if _, err := bufr.ReadByte(); err == io.EOF {
			return
		} else if err != nil {
			timestamp := rs.consumeTimestamp()
			log.Printf("Failed to read next HTTP request first byte to the client %v at %s: %v", clientAddr, timestamp.Format(time.RFC3339Nano), err)
			return
		} else {
			if err := bufr.UnreadByte(); err != nil {
				timestamp := rs.consumeTimestamp()
				log.Printf("Failed to unread next HTTP request first byte to the client %v at %s: %v", clientAddr, timestamp.Format(time.RFC3339Nano), err)
				return
			}
		}

		// parse request
		req, err := http.ReadRequest(bufr)
		if err != nil {
			timestamp := rs.consumeTimestamp()
			log.Printf("Failed to read HTTP request from the client %v at %s: %v", clientAddr, timestamp.Format(time.RFC3339Nano), err)
			return
		}

		// set internal headers
		timestamp := rs.consumeTimestamp()
		req.RemoteAddr = clientAddr.String()
		req.Header.Set(conjoinPcapKeyHeader, req.RemoteAddr)
		req.Header.Set(timestampPcapKeyHeader, timeToUnixNanoStr(timestamp))

		// discard body
		if req.ContentLength > 0 && req.Body != http.NoBody {
			_, err = io.CopyBuffer(io.Discard, req.Body, copyBuf[:])
			_ = req.Body.Close()
			if err != nil {
				log.Printf("Failed to read HTTP body from the client %v at %s: %v", clientAddr, timestamp.Format(time.RFC3339Nano), err)
				return
			}
		}

		// send parsed request
		reqCh <- req
	}
}

func parseHTTPResponse(rs *tcpReaderStream, clientAddr *net.TCPAddr, resCh chan *http.Response, stat *pcapHttpStreamStat) {
	stat.startRes()
	defer stat.completeRes()

	var copyBuf [4096]byte
	bufr := bufio.NewReader(rs)
	for {
		// check EOF
		if _, err := bufr.ReadByte(); err == io.EOF {
			return
		} else if err != nil {
			timestamp := rs.consumeTimestamp()
			log.Printf("Failed to read next HTTP response first byte to the client %v at %s: %v", clientAddr, timestamp.Format(time.RFC3339Nano), err)
			return
		} else {
			if err := bufr.UnreadByte(); err != nil {
				timestamp := rs.consumeTimestamp()
				log.Printf("Failed to unread next HTTP response first byte to the client %v at %s: %v", clientAddr, timestamp.Format(time.RFC3339Nano), err)
				return
			}
		}

		// parse response
		res, err := http.ReadResponse(bufr, nil)
		if err != nil {
			timestamp := rs.consumeTimestamp()
			log.Printf("Failed to read HTTP response to the client %v at %s: %v", clientAddr, timestamp.Format(time.RFC3339Nano), err)
			return
		}

		// set internal headers
		timestamp := rs.consumeTimestamp()
		res.Header.Set(conjoinPcapKeyHeader, clientAddr.String())
		res.Header.Set(timestampPcapKeyHeader, timeToUnixNanoStr(timestamp))

		// copy body to buffer
		if res.ContentLength != 0 && res.Body != http.NoBody {
			var bb bytes.Buffer
			_, err = io.CopyBuffer(&bb, res.Body, copyBuf[:])
			if res.Close {
				_ = res.Body.Close()
			}
			if err != nil {
				log.Printf("Failed to read HTTP body to the client %v at %s: %v", clientAddr, timestamp.Format(time.RFC3339Nano), err)
				return
			}
			res.Body = io.NopCloser(&bb)
		}

		// send parsed response
		resCh <- res
	}
}

func timeToUnixNanoStr(t time.Time) string {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(t.UnixNano()))

	var e [16]byte
	base64.RawStdEncoding.Encode(e[:], b[:])
	return string(e[:base64.RawStdEncoding.EncodedLen(len(b))])
}

func unixNanoStrToTime(s string) (time.Time, error) {
	var b [8]byte
	_, err := base64.RawStdEncoding.Decode(b[:], []byte(s))
	if err != nil {
		return time.Time{}, err
	}

	unixNano := binary.LittleEndian.Uint64(b[:])
	return time.Unix(0, int64(unixNano)), nil
}
