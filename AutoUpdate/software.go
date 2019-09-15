package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func Wireshark() (string,string) {
	var versparam []string

	dlinkInfo := "https://www.wireshark.org/download.html"
	client := new(http.Client)
	res,_ := client.Get(dlinkInfo)
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		h,_ := ioutil.ReadAll(res.Body)
		verRegexp := regexp.MustCompile("Stable Release \\(.*?\\)")
		versparam = verRegexp.FindStringSubmatch(string(h))
		verRegexp = regexp.MustCompile("\\d.*\\d")
		versparam = verRegexp.FindStringSubmatch(string(versparam[0]))
	}
	dlink := "https://1.as.dl.wireshark.org/win64/Wireshark-win64-" + versparam[0] + ".exe"
	return dlink,"wireshark.exe"
}

func other(linkfileargs string) {
	link := readlink(linkfileargs)
	for _,v := range link {
		download(dealink(v))
	}
}

func dealink(dlink string) (string,string) {
	fnameSplit := strings.Split(dlink,"/")
	return dlink,fnameSplit[len(fnameSplit) - 1]
}

func readlink(linkfile string) []string {
	f,err := os.Open(linkfile)
	if err != nil {
		return nil
	}
	buf := bufio.NewReader(f)
	var link []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		//这这要放在判断结尾之前，不然最后一行会无法写入
		link = append(link,line)
		if err == io.EOF {
			return link
		}
	}
}