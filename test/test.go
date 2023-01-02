package test

import (
	"log"
	"os/exec"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func RegisterTest(e *echo.Echo, lazy bool) {

	if !lazy {
		log.Println("Building Tests...")
		abs, err := filepath.Abs("./test/build.sh")

		if err != nil {
			log.Fatal(err)
		}

		cmd, err := exec.Command("/bin/sh", abs).Output()

		log.Println(string(cmd))

		if err != nil {
			log.Fatal(err)
		}
	}

	e.Static("/test", "test/build")
}
