package dhcpPacket

import (
	"encoding/binary"
	"fmt"
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

type Options struct {
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
	tmp := make([]byte, 4)
	binary.BigEndian.PutUint32(tmp, value)
	copy(d.xid[:], tmp)
}

func (f *Flags) ToBytes() (bytePacket [2]byte) {

	if f.broadcast {
		fmt.Println("Test")
		bytePacket = [2]byte{0x80, 0x00}
	}

	return [2]byte{0x80, 0x00}
}

func (f *Flags) SetBroadcast(value bool) {
	f.broadcast = value
}

func (o *Options) ToBytes() (bytePacket [10]byte) {
	return
}
