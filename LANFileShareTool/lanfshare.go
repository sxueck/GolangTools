package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
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

func getIPByMyself(lisPort string) bool {
	//addrs, err := net.InterfaceAddrs()
	ifaceAddr, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	var (
		ip net.IP
		//priorityJu = make([]string, 1)
		//judge ip net priority
	)

	validCount := 0

	for _, i := range ifaceAddr {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP.To4()
			case *net.IPAddr:
				ip = v.IP.To4()
			}
			//fmt.Println("check " + ip.String())
			if ip != nil && !ip.IsLoopback() {
				//delete ipv6 address and localhost address
				if judgeNetStatus(createICMP(uint16(1)), ip.String()) {
					fmt.Println(ip.String() + lisPort)
					validCount++
				}
			}
		}
	}

	if validCount == 0 {
		return false
	}

	return true
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
	icmp.Checknum = checkSum(buffer.Bytes())
	buffer.Reset()
	return icmp
}

func checkSum(data []byte) uint16 {
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
		sum += uint32(data[index])
	}
	sum += (sum >> 16)
	return uint16(^sum)
}

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

	fmt.Printf("选择一个大于1024的端口[默认8000]:")
	fmt.Scanln(&port)

	if strings.Compare(port, "") == 0 {
		return ":8000"
	}

	return (":" + port)
}

//TerminalInput .
func TerminalInput() string {
	if len(os.Args) != 1 {
		//终端输入
		return os.Args[1]
	}
	return ""
}

func fileUploadHander(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		uploadPage := `
		<html>
		<head>
			<title>Upload File Page</title>
		</head>
		<body>
			<form enctype="multipart/form-data" action="{{.PostURL}}" method="POST">
				<input type="file" name="uFile" />
				<input type="submit" value="上传" />
			</form>
		</body>
		</html>
		`
		t := template.New("upload.html")
		t, _ = t.Parse(uploadPage)

		uploadURL := struct {
			PostURL string
		}{
			PostURL: "http://" + r.Host + "/upload",
		}

		t.Execute(w, uploadURL)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, hander, err := r.FormFile("uFile")
		if err != nil {
			fmt.Println(err)
		}

		defer file.Close()

		if file == nil {
			fmt.Fprintln(w, "你没有选择文件")
			return
		}

		os.Mkdir("uploadDir/", 0666)

		f, err := os.OpenFile("uploadDir/"+hander.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintln(w, "文件写入错误")
		} else {
			fmt.Fprintln(w, "文件上传成功")
		}

		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {

	var lisPort string

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/upload", fileUploadHander)

Loop:

	lisPort = ""
	if lisPort = TerminalInput(); lisPort == "" {
		lisPort = choicePort()
	}

	//fmt.Println(lisPort)
	if !getIPByMyself(lisPort) {
		fmt.Println("0.0.0.0" + lisPort)
	}

	fmt.Printf(getCurrentDir() + " 正在被共享\n")
	fmt.Println("在地址后加入/upload即可上传文件")

	if err := http.ListenAndServe(lisPort, nil); err != nil {
		fmt.Println("你选择的端口已被占用,请重新选择")
		goto Loop
	}
}
