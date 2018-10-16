package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var username, password string

//GetPage .
func GetPage(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("i/o timeout")
	}

	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()
	html, _ := ioutil.ReadAll(resp.Body)
	return string(html)
}

//GetKey .
func GetKey(hCode string) (string, string, string) {
	//提取出POST需要的数据
	//这里不考虑重用性，只针对当前版本的签到系统
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(hCode))
	if err != nil {
		log.Fatal(err)
	}

	keyVIEWSTATE, _ := dom.Find("#__VIEWSTATE").Attr("value")
	keyEVENTVALIDATION, _ := dom.Find("#__EVENTVALIDATION").Attr("value")
	keyVIEWSTATEGENERATOR, _ := dom.Find("#__VIEWSTATEGENERATOR").Attr("value")

	return keyEVENTVALIDATION, keyVIEWSTATE, keyVIEWSTATEGENERATOR
}

//LoginPostData .
func LoginPostData(EVENTVALIDATION, VIEWSTATE, VIEWSTATEGENERATOR, PostURL string) bool {
	postRaw := "__VIEWSTATE=" + url.QueryEscape(VIEWSTATE) + "&__VIEWSTATEGENERATOR=" + VIEWSTATEGENERATOR + "&__EVENTVALIDATION=" + url.QueryEscape(EVENTVALIDATION) + "&TextBoxTeacherName=" + url.QueryEscape(username) + "&TextBoxPassword=" + password + "&ButtonLogin=%E7%99%BB%E5%BD%95"
	//构建POST数据

	//fmtfmt.Println(postRaw)
	pReq, err := http.NewRequest("POST", PostURL, strings.NewReader(postRaw))

	pReq.Header.Add("Connection", "Keep-Alive")
	pReq.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
	//必须加入这个header，否则post无效
	//下面的header不加入会被WAF拦截
	pReq.Header.Add("X-originating-IP", "127.0.0.1")
	pReq.Header.Add("Origin", "http://cc.szpt.edu.cn")
	pReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	pReq.Header.Add("Referer", "http://cc.szpt.edu.cn/Login.aspx")

	jar, _ := cookiejar.New(nil)

	client := http.Client{}
	client.Jar = jar
	repReq, err := client.Do(pReq)

	//fmtfmt.Println(repReq.Header)
	//fmt.Println(repReq)
	//提交请求

	if err != nil {
		log.Fatal("Data POST Error")
	}

	defer pReq.Body.Close()

	ConvertCheckLoginByte, _ := ioutil.ReadAll(repReq.Body)
	var checkLogin = string(ConvertCheckLoginByte)

	//fmtfmt.Println(checkLogin)
	//rpCheckLogin := regexp.Compile("/退出登录/")
	//if checkLogin, _ := rpCheckLogin.FindAllString(ioutil.ReadAll(pReq.Body)); checkLogin != "" {
	if !strings.Contains(checkLogin, "退出登录") {
		log.Fatal("请检查你的账号或密码")
	}
	return true

}

//PunchCardAuthKey .
func PunchCardAuthKey(hCode string) string {
	//VIEWSTATE,_ := resp.Find("__VIEWSTATE")
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(hCode))
	resp, _ := dom.Find("#__VIEWSTATE").Attr("value")
	return resp
}

//PunchCard .
func PunchCard(pCardAuthKey, authurl string) bool {
	postRaw := "__VIEWSTATE=" + pCardAuthKey + "&TextBoxStudentNo=" + url.QueryEscape(username) + "&ButtonSign=%E6%8F%90%E4%BA%A4"

	punchc, err := http.NewRequest("POST", authurl, strings.NewReader(postRaw))

	if err != nil {
		log.Fatal("签到页面错误:" + err.Error())
	}

	punchc.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
	client := http.Client{}
	_, err = client.Do(punchc)

	if err != nil {
		fmt.Println("签到未开通")
		return false
	}

	return true
}

//Init Login
func Init() {
	loginURL := "http://cc.szpt.edu.cn/Login.aspx"

	getUser()

	//loginURL := "http://localhost:8080/"
	//本地数据测试

	htmlcode := GetPage(loginURL)
	key1, key2, key3 := GetKey(htmlcode)
	if LoginPostData(key1, key2, key3, loginURL) {
		fmt.Println("登录成功")
	}
}

func getUser() {
	file, err := os.Open("user.txt")
	if err != nil {
		log.Fatal("请检查user.txt是否存在")
	}

	defer file.Close()

	buffer := bufio.NewReader(file)
	for i := 1; i < 3; i++ {
		info, _, _ := buffer.ReadLine()
		if i == 1 {
			username = string(info)
		} else {
			password = string(info)
		}
	}
}

func main() {
	pCardAuthKey := "http://cc.szpt.edu.cn/sSign.aspx"
	//pCardAuthKey := "http://localhost:8080/"

	Init()
	aKey := PunchCardAuthKey(GetPage(pCardAuthKey))

	fmt.Println("开始抢沙发")
	for stop := false; !stop; {
		if PunchCard(aKey, pCardAuthKey) {
			fmt.Println("签到完成")
			stop = true
		}
	}
}
