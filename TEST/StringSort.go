package main

import (
	"fmt"
	"sort"
)

func main() {
	//var ip = []string{"10.1.34.15", "10.1.35.44", "192.168.1.2", "172.168.141", "1.1.1.1"}
	ip := []string{"10.1.34.15", "10.1.35.44", "192.168.1.2", "172.168.141", "1.1.1.1"}
	ip = append(ip, "255.255.255.254")
	sort.Sort(sort.StringSlice(ip))
	fmt.Println(ip)
}
