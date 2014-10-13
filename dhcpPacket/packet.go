package dhcpPacket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type dhcpPacket struct {
	op      byte
	htype   byte
	hlen    byte
	hops    byte
	xid     [4]byte
	secs    [2]byte
	flags   Flags
	ciaddr  [4]byte
	yiaddr  [4]byte
	siaddr  [4]byte
	giaddr  [4]byte
	chaddr  [16]byte
	sname   [64]byte
	file    [128]byte
	Options Options
}

type Flags struct {
	broadcast bool
}

// Returns the byte array forming the packet
func (d dhcpPacket) ToBytes() (bytePacket []byte) {
	bytePacket = append(bytePacket, d.op, d.htype, d.hlen, d.hops)
	bytePacket = append(bytePacket, d.xid[:]...)
	bytePacket = append(bytePacket, d.secs[:]...)

	flagsBytes := d.flags.ToBytes()
	bytePacket = append(bytePacket, flagsBytes[:]...)

	bytePacket = append(bytePacket, d.ciaddr[:]...)
	bytePacket = append(bytePacket, d.yiaddr[:]...)
	bytePacket = append(bytePacket, d.siaddr[:]...)
	bytePacket = append(bytePacket, d.giaddr[:]...)
	bytePacket = append(bytePacket, d.chaddr[:]...)
	bytePacket = append(bytePacket, d.sname[:]...)
	bytePacket = append(bytePacket, d.file[:]...)

	magic := [...]byte{0x63, 0x82, 0x53, 0x63}
	bytePacket = append(bytePacket, magic[:]...)

	optionsBytes := d.Options.ToBytes()
	bytePacket = append(bytePacket, optionsBytes[:]...)

	return
}

func NewDhcpPacket() (d *dhcpPacket) {

	d = new(dhcpPacket)
	d.htype = 0x01
	d.hlen = 0x06
	d.hops = 0x00

	return d

}

func (d *dhcpPacket) SetOp(value int) {
	if value == 1 {
		d.op = 0x01
	} else {
		d.op = 0x02
	}
}

func (d *dhcpPacket) GetOp() (value int) {
	if d.op == 0x01 {
		return 1
	} else {
		return 2
	}
}

func (d *dhcpPacket) GetXid() (value uint32) {
	value = binary.LittleEndian.Uint32(d.xid[:])
	return
}

func (d *dhcpPacket) SetXid(value uint32) {
	copy(d.xid[:], intToBytes(value))
}

func (d *dhcpPacket) SetCiaddr(value string) {
	tmp := net.ParseIP(value)
	copy(d.ciaddr[:], tmp[12:])
}

func (d *dhcpPacket) SetYiaddr(value string) {
	tmp := net.ParseIP(value)
	copy(d.yiaddr[:], tmp[12:])
}
func (d *dhcpPacket) SetSiaddr(value string) {
	tmp := net.ParseIP(value)
	copy(d.siaddr[:], tmp[12:])
}
func (d *dhcpPacket) SetGiaddr(value string) {
	tmp := net.ParseIP(value)
	copy(d.giaddr[:], tmp[12:])
}

func (d *dhcpPacket) SetChaddr(value string) {
	tmp, err := net.ParseMAC(value)
	if err != nil {
		fmt.Println("Error whild parsing MAC")
	}
	copy(d.chaddr[:], tmp)
}

func (d *dhcpPacket) GetChaddr() (address net.HardwareAddr) {
	address = d.chaddr[:]
	return
}

func (d *dhcpPacket) SetBroadcast(value bool) {
	d.flags.SetBroadcast(value)
}

func (f *Flags) ToBytes() (bytePacket [2]byte) {

	if f.broadcast {
		fmt.Println("Test")
		bytePacket = [2]byte{0x80, 0x00}
		return
	}

	return [2]byte{0x00, 0x00}
}

func (f *Flags) SetBroadcast(value bool) {
	f.broadcast = value
}

func ParseDhcpPacket(b []byte, o *dhcpPacket) (err error) {

	o.op = b[0]
	o.htype = b[1]
	o.hlen = b[2]
	o.hops = b[3]
	copy(o.xid[:], b[4:7])
	copy(o.secs[:], b[8:9])
	o.flags = parseFlags(b[10:11])
	copy(o.ciaddr[:], b[12:15])
	copy(o.yiaddr[:], b[16:19])
	copy(o.siaddr[:], b[20:23])
	copy(o.giaddr[:], b[24:27])
	copy(o.chaddr[:], b[28:43])
	copy(o.sname[:], b[44:107])
	copy(o.file[:], b[108:237])
	o.Options = parseOptions(b[238:])

	return

}

func parseFlags(f []byte) (o Flags) {
	if f[0] == byte(0x80) {
		o.SetBroadcast(true)
	}

	return
}

// Helpers

func intToBytes(value uint32) (output []byte) {

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, value)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	return buf.Bytes()
}
