package dhcpPacket

import (
	"bytes"
	"encoding/binary"
	"errors"
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
		bytePacket = [2]byte{0x80, 0x00}
		return
	}

	return [2]byte{0x00, 0x00}
}

func (f *Flags) SetBroadcast(value bool) {
	f.broadcast = value
}

func ParseDhcpPacket(b []byte, o *dhcpPacket) (err error) {

	if len(b) < 237 {
		return errors.New("Mauvais packet")
	}

	magic := [...]byte{0x63, 0x82, 0x53, 0x63}
	parsedMagic := [4]byte{}
	copy(parsedMagic[:], b[236:240])

	if parsedMagic != magic {
		return errors.New("Mauvais Magic Cookie")
	}

	o.op = b[0]
	o.htype = b[1]
	o.hlen = b[2]
	o.hops = b[3]
	copy(o.xid[:], b[4:8])
	copy(o.secs[:], b[8:10])
	o.flags = parseFlags(b[10:12])
	copy(o.ciaddr[:], b[12:16])
	copy(o.yiaddr[:], b[16:20])
	copy(o.siaddr[:], b[20:24])
	copy(o.giaddr[:], b[24:28])
	copy(o.chaddr[:], b[28:44])
	copy(o.sname[:], b[44:108])
	copy(o.file[:], b[108:236])

	// On saute le magic Cookie
	o.Options = parseOptions(b[240:])

	return

}

func (d *dhcpPacket) GetMessageType() (msgType int) {
	value, err := d.Options.GetOption(53)
	if err != nil {
		fmt.Println(err)
	}

	msgType = byteToInt(value)

	return
}

func (d *dhcpPacket) GetRequestedIp() (ip net.IP) {
	value, err := d.Options.GetOption(50)
	if err != nil {
		fmt.Println(err)
	}

	ip = net.IP(value)
	return
}

func (d *dhcpPacket) GetHostName() (hostname string) {
	value, err := d.Options.GetOption(12)
	if err != nil {
		fmt.Println(err)
	}

	hostname = string(value)
	return
}

func (d *dhcpPacket) GetMaximumPacketSize() (size int) {
	value, err := d.Options.GetOption(57)
	if err != nil {
		fmt.Println(err)
	}

	size = int(binary.BigEndian.Uint16(value))
	return
}

func (d *dhcpPacket) GetParameterRequestList() (list []int) {
	value, err := d.Options.GetOption(55)
	if err != nil {
		fmt.Println(err)
	}

	buf := bytes.NewBuffer(value)

	for i := 0; i < len(value); i++ {
		codeRaw := make([]byte, 1)

		n, err := buf.Read(codeRaw)
		if err != nil {
			fmt.Println(err)
		}
		if n == 0 {
			fmt.Println("Unexpected End")
		}
		list = append(list, byteToInt(codeRaw))
	}

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

func byteToInt(value []byte) (output int) {
	value = append(value, byte(0x00))
	return int(binary.LittleEndian.Uint16(value))
}
