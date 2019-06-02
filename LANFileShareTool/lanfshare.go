package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"
)

/*
#include <stdio.h>
#include <ShlObj.h>

char* GetSharePath() {
    char* szPathName = (char*)malloc(sizeof(char)*MAX_PATH);
    BROWSEINFO bInfo = { 0 };
    bInfo.hwndOwner = GetForegroundWindow();
    bInfo.lpszTitle = TEXT("Choice Your Folder");
    bInfo.ulFlags = BIF_RETURNONLYFSDIRS | BIF_USENEWUI,BIF_UAHINT| BIF_NONEWFOLDERBUTTON;
    LPITEMIDLIST lpDlist = SHBrowseForFolder(&bInfo);
    if(lpDlist != NULL) {
        SHGetPathFromIDList(lpDlist, szPathName);
    }
    return szPathName;
}
*/
import "C"

//清空屏幕
var clear map[string]func()

//ICMP Proccotol Datapack
type ICMP struct {
	Type        uint8
	Code        uint8
	Checknum    uint16
	Identifer   uint16
	SequenceNum uint16
}

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
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
		return false
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

//检查用户输入(从终端)
func cmdInputPort() string {
	if len(os.Args) != 1 {
		//终端输入
		inValue1 := os.Args[1]
		if port, err := strconv.Atoi(inValue1); err == nil {
			//用户第一参数端口号
			if port < 65535 && port > 0 {
				return (":" + inValue1)
			}
		}
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

func shareFileCheck(path chan string) {

	var sharePath string

	//如果没有参数2但是有参数1则将参数2当为共享路径
	//如果不包含参数，则设置为当前目录
	switch len(os.Args) {
	case 2:
		sharePath = os.Args[1]
	case 3:
		sharePath = os.Args[2]
	default:
		if runtime.GOOS == "windows" {
			sharePath = winUserpath()
		} else {
			sharePath = getCurrentDir()
		}
	}

	_, err := os.Stat(sharePath)
	if err != nil {
		path <- getCurrentDir()
	}
	path <- sharePath
}

func callClear() {
	//runtime.GOOS -> linux, windows, darwin etc.
	value, resOk := clear[runtime.GOOS]
	if resOk {
		value()
	}
}

func winUserpath() string {
	path := C.GetSharePath()
	cString := []byte(C.GoString(path))
	enc := mahonia.NewDecoder("GBK")
	_, cdata, _ := enc.Translate(cString, true)
	var upath string = string(cdata[:])
	return upath
}


func main() {

	var lisPort string
	sharePathThread := make(chan string)

	go shareFileCheck(sharePathThread)
	sharePATH := <-sharePathThread

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(sharePATH)))
	mux.HandleFunc("/upload", fileUploadHander)


Loop:

	lisPort = ""
	if lisPort = cmdInputPort(); lisPort == "" {
		lisPort = choicePort()
	}

	//fmt.Println(lisPort)
	if !getIPByMyself(lisPort) {
		fmt.Println("你的网卡信息获取失败或者无可用内网IP")
		fmt.Println("请手动查找你的IP")
	}

	fmt.Printf(sharePATH + " 正在被共享\n")
	fmt.Println("在地址后加入/upload即可上传文件")
	fmt.Println("")

	server := http.Server{
		Addr:    "0.0.0.0"+lisPort,
		Handler: mux,
	}

	if err := server.ListenAndServeTLS("key/lfs.crt", "key/lfs.key"); err != nil {
		//log.Fatal(err)
		//TODO:这里闪太快了，需要优化
		fmt.Println("你选择的端口已被占用,请重新选择")
		time.Sleep(time.Second)

		//清空终端界面
		callClear()
		goto Loop
	}
}
