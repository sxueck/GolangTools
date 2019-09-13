package main

import (
	"io/ioutil"
	"net/http"
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

func other() {
	//fidder4
	download(dealink("https://telerik-fiddler.s3.amazonaws.com/fiddler/FiddlerSetup.exe"))
	//teamviewer
	download(dealink("https://dl.teamviewer.com/download/TeamViewer_Setup.exe"))
	//chrome
	download(dealink("https://dl.google.com/tag/s/lang%3Dzh-CN/chrome/install/ChromeStandaloneSetup64.exe"))
}

func dealink(dlink string) (string,string) {
	fnameSplit := strings.Split(dlink,"/")
	return dlink,fnameSplit[len(fnameSplit) - 1]
}