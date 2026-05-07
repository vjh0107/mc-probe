package ping

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

type Response struct {
	Version     Version         `json:"version"`
	Players     Players         `json:"players"`
	Description json.RawMessage `json:"description"`
	Favicon     string          `json:"favicon,omitempty"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type Players struct {
	Max    int `json:"max"`
	Online int `json:"online"`
}

func Query(addr string, timeout time.Duration) (*Response, error) {
	body, err := query(addr, timeout)
	if err != nil {
		return nil, err
	}
	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing status response: %w", err)
	}
	return &resp, nil
}

func query(addr string, timeout time.Duration) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(timeout))

	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("parsing address: %w", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 65535 {
		return nil, fmt.Errorf("invalid port %q", portStr)
	}

	if err := sendHandshake(conn, host, uint16(port)); err != nil {
		return nil, fmt.Errorf("sending handshake: %w", err)
	}

	if err := sendStatusRequest(conn); err != nil {
		return nil, fmt.Errorf("sending status request: %w", err)
	}

	body, err := readStatusResponse(conn)
	if err != nil {
		return nil, fmt.Errorf("reading status response: %w", err)
	}

	return body, nil
}

func sendHandshake(conn net.Conn, host string, port uint16) error {
	var payload bytes.Buffer

	writeVarInt(&payload, 0x00)
	writeVarInt(&payload, -1)
	writeString(&payload, host)
	binary.Write(&payload, binary.BigEndian, port)
	writeVarInt(&payload, 1)

	return writePacket(conn, payload.Bytes())
}

func sendStatusRequest(conn net.Conn) error {
	var payload bytes.Buffer
	writeVarInt(&payload, 0x00)
	return writePacket(conn, payload.Bytes())
}

func readStatusResponse(conn net.Conn) ([]byte, error) {
	if _, err := readVarInt(conn); err != nil {
		return nil, fmt.Errorf("reading packet length: %w", err)
	}

	packetID, err := readVarInt(conn)
	if err != nil {
		return nil, fmt.Errorf("reading packet id: %w", err)
	}

	if packetID != 0x00 {
		return nil, fmt.Errorf("unexpected packet id: %d", packetID)
	}

	bodyLen, err := readVarInt(conn)
	if err != nil {
		return nil, fmt.Errorf("reading body length: %w", err)
	}
	if bodyLen < 0 {
		return nil, fmt.Errorf("negative body length: %d", bodyLen)
	}

	body := make([]byte, bodyLen)
	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	return body, nil
}

func writePacket(conn net.Conn, data []byte) error {
	var buf bytes.Buffer
	writeVarInt(&buf, int32(len(data)))
	buf.Write(data)
	_, err := conn.Write(buf.Bytes())
	return err
}

func writeVarInt(buf *bytes.Buffer, value int32) {
	uval := uint32(value)
	for {
		if uval&^0x7F == 0 {
			buf.WriteByte(byte(uval))
			return
		}
		buf.WriteByte(byte(uval&0x7F) | 0x80)
		uval >>= 7
	}
}

func writeString(buf *bytes.Buffer, s string) {
	writeVarInt(buf, int32(len(s)))
	buf.WriteString(s)
}

func readVarInt(conn net.Conn) (int32, error) {
	var result int32
	var shift uint
	buf := make([]byte, 1)

	for {
		if _, err := conn.Read(buf); err != nil {
			return 0, err
		}
		result |= int32(buf[0]&0x7F) << shift
		if buf[0]&0x80 == 0 {
			return result, nil
		}
		shift += 7
		if shift >= 35 {
			return 0, fmt.Errorf("VarInt too big")
		}
	}
}
