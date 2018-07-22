package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Static("/", "html")
	e.GET("/addme", addMe)
	e.Logger.Fatal(e.Start(":80"))
}

func addMe(c echo.Context) error {
	ip := c.RealIP()
	if strings.Index(ip, "::ffff:") == 0 { // fix nginx pseudo ipv6 forward
		ip = strings.Replace(ip, "::ffff:", "", 1)
	}
	info, err := getASNInfo(ip)
	if err != nil || !info.Announced {
		log.Println(err)
		log.Println(info)
		log.Println(ip)
		return c.JSON(http.StatusOK, map[string]string{"result": "failed", "msg": "Failed to get ASN"})
	}

	err = branchChange(info)
	if err != nil && err.Error() == "Branch exists" {
		log.Println(err)
		return c.JSON(http.StatusOK, map[string]string{"result": "existing", "msg": "ASN already reported to us"})
	}

	err = createPR(info)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusOK, map[string]string{"result": "failed", "msg": "Failed to creare PR"})
	}

	return c.JSON(http.StatusOK, map[string]string{"result": "ok", "msg": "ASN added to list"})
}
