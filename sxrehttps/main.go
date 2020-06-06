package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	log.Fatal(StartServe())
}

func RedirectToHttps(c echo.Context) error {
	src := c.Request().Host + c.Request().RequestURI
	fmt.Println("RedirectToHttps", src)
	return c.Redirect(302, fmt.Sprintf("https://%s",src))
}

func StartServe() error {
	e := echo.New()
	e.HideBanner = true
	e.GET("/*", RedirectToHttps)
	e.POST("/*", RedirectToHttps)
	e.PUT("/*", RedirectToHttps)
	e.DELETE("/*", RedirectToHttps)
	return e.Start(":80")
}
