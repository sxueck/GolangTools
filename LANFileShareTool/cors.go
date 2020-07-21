package main

import (
	"fmt"
	"github.com/semicircle/gocors"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
)

var (
	flags   []cli.Flag
	port    int
	address string
	cors    bool
)

func init() {
	flags = []cli.Flag{
		cli.IntFlag{
			Name:        "p, port",
			Value:       8000,
			Destination: &port,
		},
		cli.StringFlag{
			Name:        "t, host, ip-address",
			Value:       "0.0.0.0",
			Destination: &address,
		},
		cli.BoolFlag{
			Name:        "c, cors",
			Usage:       "default close",
			Destination: &cors,
		},
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "SimpleHTTPServer"
	app.Usage = "Application Usage"
	app.HideVersion = true
	app.Flags = flags
	app.Action = Action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func CORSEnable(enable bool) http.Handler {
	cs := gocors.New()
	cs.SetAllowOrigin("*")
	fHttp := http.FileServer(http.Dir("."))
	if enable {
		return cs.Handler(fHttp)
	}
	return fHttp
}

func Action(c *cli.Context) error {
	ipAddress := fmt.Sprintf("%s:%d", address, port)
	fmt.Println("listening", ipAddress, "CORS", cors)
	log.Fatal(
		http.ListenAndServe(ipAddress, CORSEnable(cors)))
	return nil
}
