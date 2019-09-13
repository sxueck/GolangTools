package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const TESTLINK = "https://mirrors.aliyun.com/archlinux/iso/archboot/latest/archlinux-2018.06-1-archboot-network.iso"

type FILEINFO struct {
	link       string
	contentLen int64
	f          *os.File
}

type DOWNLOAD struct {
	start      int64
	end        int64
	threadsNum int
}

var (
	ifile FILEINFO
	fdown DOWNLOAD

	//下载线程数量
	dthreads = 4

	wg    sync.WaitGroup
	mutex sync.Mutex
)

func main() {
	download(Wireshark())
	other()
}

func download(link, fpath string) error {
	//清空结构体，避免之前的信息影响后面
	ifile = FILEINFO{link: link}

	if _, err := os.Stat(fpath); err != nil {
		if os.IsNotExist(err) {
			if createFile, err := os.Create(fpath); err == nil {
				ifile.f = createFile
				defer ifile.f.Close()
			}
		}
	} else {
		log.Println("File exists,covering...")
	}

	if ifile.checkHead() {
		dispSliceDownload()
	} else {
		log.Println(link, "Direct Download")
		directDownload()
	}
	return nil
}

//判断是否支持多线程下载
func (*FILEINFO) checkHead() bool {
	var client = new(http.Client)
	if request, err := http.NewRequest("HEAD", ifile.link, nil); err == nil {
		if response, err := client.Do(request); err == nil {
			defer response.Body.Close()
			ifile.contentLen = response.ContentLength
			//包含accept-ranges: bytes即可多线程下载
			if strings.Compare(response.Header.Get("accept-ranges"), "bytes") == 0 {
				return true
			}
			return false
		}
	}
	return false
}

func directDownload() {
	log.Println("Download Starting...")
	res,err := http.Get(ifile.link)
	if err != nil {
		log.Printf("download error : ",err)
	}
	io.Copy(ifile.f,res.Body)
	log.Println("Download Success")
}

func dispSliceDownload() {
	log.Println("Download Starting...")
	defer ifile.f.Close()
	//计算每个线程的下载区块大小
	dispi := ifile.contentLen / int64(dthreads)
	if ifile.contentLen-int64(dthreads)*dispi != 0 {
		dthreads += 1
	}
	wg.Add(dthreads)
	for i := 0; i < dthreads; i++ {
		//计算每个线程的下载区块位置
		start := int64(int64(i) * dispi)
		end := start + dispi
		//判断结尾是否到达
		if end > ifile.contentLen {
			end -= end - ifile.contentLen
		}

		//开始构建请求
		if req, err := http.NewRequest("GET", ifile.link, nil); err == nil {
			req.Header.Set(
				"range",
				"bytes="+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10))
			//强迫症
			fdown.start = start
			fdown.end = end
			fdown.threadsNum = i + 1

			go sliceDownload(req, fdown, ifile)
		}
	}
	wg.Wait()
	log.Println("Download Success")
}

func sliceDownload(req *http.Request, fdown DOWNLOAD, f FILEINFO) {
	client := new(http.Client)
	if res, err := client.Do(req); err == nil && strings.Contains(res.Status, "206") {
		defer res.Body.Close()
		b,err := ioutil.ReadAll(res.Body)
		if err != io.EOF && err != nil {
			log.Printf("threads %d download error : %s",fdown.threadsNum,err)
		}
		//防止同时对同一文件进行读写
		mutex.Lock()
		_,err = ifile.f.WriteAt(b,fdown.start)
		mutex.Unlock()
		if err != nil  {
			log.Printf("threads %d write error: %s", fdown.threadsNum,err)
		}
	}
	wg.Done()
}