package dhcpUDP

import (
	"fmt"
	"net"
)

var listener *Listener

type Listener struct {
	C chan []byte
}

func InitListener(handler func([]byte)) {
	listener = NewListener()
	go listener.Run()

	for {
		tmp := <-listener.C
		handler(tmp)
	}
}

func NewListener() *Listener {
	out := new(Listener)
	out.C = make(chan []byte)
	return out
}

func (l *Listener) Run() {

	addr, err := net.ResolveUDPAddr("udp", ":67")

	if err != nil {
		fmt.Println(err)
	}
	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for {
			var buf [1024]byte
			sock.ReadFromUDP(buf[:])
			l.C <- buf[:]
		}
	}()
}
