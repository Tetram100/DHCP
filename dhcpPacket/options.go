package dhcpPacket

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type Options struct {
	options []Option
}

type Option struct {
	code   uint32
	length uint32
	value  []byte
}

func (o *Options) ToBytes() (bytePacket []byte) {
	for _, option := range o.options {
		tmp := option.ToBytes()
		bytePacket = append(bytePacket, tmp[:]...)
	}

	return
}

func (o *Options) Add(code uint32, value []byte) {

	var option Option

	if value != nil {
		option = Option{code, uint32(len(value)), value}
	} else {
		option = Option{code, 0, nil}
	}

	o.options = append(o.options, option)
}

func (o *Option) ToBytes() (bytePacket []byte) {
	code := intToBytes(o.code)
	bytePacket = append(bytePacket, code[0])
	if o.length != 0 {
		length := intToBytes(o.length)
		bytePacket = append(bytePacket, length[0])
	}
	bytePacket = append(bytePacket, o.value[:]...)

	return
}

func parseOptions(i []byte) (o Options) {

	buf := bytes.NewBuffer(i)

	for {
		err := o.parseAdd(buf)
		if err != nil {
			if err.Error() != "End" {
				fmt.Println(err)
			}
			break
		}
	}

	return
}

func (o *Options) parseAdd(buf *bytes.Buffer) (err error) {
	var n int

	codeRaw := make([]byte, 1)
	lengthRaw := make([]byte, 1)

	n, err = buf.Read(codeRaw)
	if err != nil {
		fmt.Println(err)
	}
	if n == 0 || codeRaw[0] == byte(0xff) {
		return errors.New("End")
	}

	n, err = buf.Read(lengthRaw)
	if err != nil {
		fmt.Println(err)
	}
	if n == 0 {
		return errors.New("End")
	}

	// On rajoute des zéros pour parse en uint16
	codeRaw = append(codeRaw, byte(0x00))
	lengthRaw = append(lengthRaw, byte(0x00))

	code := binary.LittleEndian.Uint16(codeRaw)
	length := binary.LittleEndian.Uint16(lengthRaw)

	data := make([]byte, length)
	n, err = buf.Read(data)
	if err != nil {
		fmt.Println(err)
	}
	if n != int(length) {
		return errors.New("Truncated Data")
	}

	fmt.Println(len(data))

	o.Add(uint32(code), data[:])
	return

}
