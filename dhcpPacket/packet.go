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
	Flags   Flags
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

	flagsBytes := d.Flags.ToBytes()
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

func (f *Flags) ToBytes() (bytePacket [2]byte) {

	if f.broadcast {
		fmt.Println("Test")
		bytePacket = [2]byte{0x80, 0x00}
	}

	return [2]byte{0x00, 0xff}
}

func (f *Flags) SetBroadcast(value bool) {
	f.broadcast = value
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
