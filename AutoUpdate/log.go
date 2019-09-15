package main

import (
	"io"
	"log"
	"os"
	"time"
)

var LOGPATH string

func wLog(info string) {
	appendToFile(func() string {
		timestamp := time.Unix(0, time.Now().Unix()*int64(time.Second))
		return timestamp.String()
	}() + " " + info)
}

func appendToFile(text string) {
	fp := fpoen(LOGPATH)
	//n,_ := fp.Seek(0,os.SEEK_END)
	//os.SEEK_END已经过时了
	n, _ := fp.Seek(0, io.SeekEnd)
	_, err := fp.WriteAt([]byte(text + "\n"), n)
	if err != nil {
		log.Println(err)
	}
	defer fp.Close()
}

func checkfile(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func fpoen(path string) *os.File {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	if checkfile(path) {
		//os.O_APPEND和seek冲突了
		if fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755); err == nil {
			return fp
		} else {
			log.Println(err)
			return nil
		}
	}
	if fp, err := os.Create(path); err == nil {
		return fp
	} else {
		log.Println(err)
		return nil
	}
}
