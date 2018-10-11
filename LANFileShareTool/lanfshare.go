package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

//ICMP Proccotol Datapack
type ICMP struct {
	Type        uint8
	Code        uint8
	Checknum    uint16
	Identifer   uint16
	SequenceNum uint16
}

func getCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func getIPByMyself() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, addr := range addrs {
		if ipnet, check := addr.(*net.IPNet); check && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if judgeNetStatus(createICMP(uint16(0)), ipnet.IP.String()) {
					return ipnet.IP.String()
				}
			}
		}
	}

	return "0.0.0.0"
}

func createICMP(seq uint16) ICMP {
	icmp := ICMP{
		Type:        8,
		Code:        0,
		Checknum:    0,
		Identifer:   0,
		SequenceNum: seq,
	}
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	//icmp.Checknum = checkSum(buffer.Bytes())
	buffer.Reset()
	return icmp
}

/*func checkSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index]<<8 + uint32(data[index]))
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)
	return uint16(^sum)
}
*/

func judgeNetStatus(icmp ICMP, host string) bool {

	raddr, err := net.ResolveIPAddr("ip", host)

	conn, err := net.DialIP("ip4:icmp", nil, raddr)
	if err != nil {
		log.Fatal("ip address error")
	}
	defer conn.Close()

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	if _, err := conn.Write(buffer.Bytes()); err != nil {
		log.Fatal(err)
	}

	conn.SetDeadline(time.Now().Add(time.Microsecond * 5))
	recv := make([]byte, 1024)
	_, err = conn.Read(recv)

	if err != nil {
		//log.Fatal("timeout")
		return false
	}

	return true
}

func choicePort() string {
	var port string

	fmt.Printf("选择一个端口[默认8000]:")
	fmt.Scanln(&port)

	if strings.Compare(port, "") == 0 {
		return ":8000"
	}

	return (":" + port)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))

	lisPort := choicePort()

	myAddr := getIPByMyself()

	if strings.Compare(myAddr, "0.0.0.0") == 0 {
		fmt.Println("您的目录已共享成功,但是访问IP需要您来手动查找")
	} else {
		fmt.Println("通过浏览器访问" + myAddr + lisPort)
	}

	fmt.Printf(getCurrentDir() + " 正在被共享")
	http.ListenAndServe(lisPort, nil)
}
