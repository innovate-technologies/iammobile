package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Static("/", "html")
	e.GET("/addme", addMe)
	e.Logger.Fatal(e.Start(":80"))
}

func addMe(c echo.Context) error {
	info, err := getASNInfo(c.RealIP())
	if err != nil || !info.Announced {
		log.Println(err)
		log.Println(info)
		log.Println(c.RealIP())
		return c.JSON(http.StatusOK, map[string]string{"result": "failed", "msg": "Failed to get ASN"})
	}

	err = branchChange(fmt.Sprintf("%d", info.AsNumber), info.AsDescription)
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
