package main // import "hello-filepath"

import (
	_ "embed"
	"fmt"
	"hello-filepath/dir"
	"log"

	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/practice-golang/lorca"
)

//go:embed assets/index.html
var index []byte

//go:embed assets/logo.png
var logo []byte

func getFiles(c echo.Context) error {
	pathRequest := new(dir.PathRequest)

	if err := c.Bind(pathRequest); err != nil {
		r := map[string]string{
			"result":  "error",
			"message": "Wrong request",
		}

		return c.JSON(http.StatusBadRequest, r)
	}

	// flist, err := GetDirectoryList("..", NAME, DESC)
	// flist, err := GetDirectoryList("..", TYPE, ASC)
	// flist, err := GetDirectoryList("..", SIZE, ASC)
	flist, err := dir.GetFileList(pathRequest.Path, dir.NAME, dir.ASC)
	if err != nil {
		fmt.Println("GetFileList:", err)

		resultError := map[string]string{
			"result":  "error",
			"message": err.Error(),
		}

		return c.JSON(http.StatusOK, resultError)
	}

	return c.JSON(http.StatusOK, flist)
}

func getDirectories(c echo.Context) error {
	pathRequest := new(dir.PathRequest)

	if err := c.Bind(pathRequest); err != nil {
		r := map[string]string{
			"result":  "error",
			"message": "Wrong request",
		}

		return c.JSON(http.StatusBadRequest, r)
	}

	dlist, err := dir.GetVisibleDirectories(pathRequest.Path, dir.NAME, dir.ASC)
	if err != nil {
		fmt.Println("GetVisibleDirectories:", err)

		resultError := map[string]string{
			"result":  "error",
			"message": err.Error(),
		}

		return c.JSON(http.StatusOK, resultError)
	}

	return c.JSON(http.StatusOK, dlist)
}

func initEcho() {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, string(index))
	})
	e.GET("/logo.png", func(c echo.Context) error {
		return c.HTML(http.StatusOK, string(logo))
	})

	e.POST("/files", getFiles)
	e.POST("/dir-upstream", getDirectories)

	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}

func initLorca() {
	cwd, _ := os.Getwd()
	profilePath := cwd + `\profile`

	lorca.DefaultChromeArgs = []string{
		"--disable-background-networking",
		"--disable-background-timer-throttling",
		"--disable-backgrounding-occluded-windows",
		"--disable-breakpad",
		"--disable-client-side-phishing-detection",
		"--disable-default-apps",
		"--disable-dev-shm-usage",
		"--disable-infobars",
		"--disable-extensions",
		"--disable-features=site-per-process",
		"--disable-hang-monitor",
		"--disable-ipc-flooding-protection",
		"--disable-popup-blocking",
		"--disable-prompt-on-repost",
		"--disable-renderer-backgrounding",
		"--disable-sync",
		"--disable-translate",
		// "--disable-windows10-custom-titlebar",
		"--metrics-recording-only",
		// "--no-first-run",
		"--no-default-browser-check",
		"--safebrowsing-disable-auto-update",
		// "--enable-automation",
		"--password-store=basic",
		"--use-mock-keychain",
	}

	// args := []string{"--ash-force-desktop", "--force-app-mode"}
	args := []string{}

	ui, err := lorca.New("http://localhost:1323", profilePath, 1024, 768, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	<-ui.Done()

	os.RemoveAll(cwd + `/profile`)
}

func main() {
	go func() { initEcho() }()
	initLorca()
}
