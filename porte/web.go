package main

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type OpenRequest struct {
	CardID  string `json:"card_id"`
	CardPIN string `json:"card_pin"`
}

func checkRequest(c echo.Context) error {
	conf := GetConfig()
	var req OpenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}

	data, err := json.Marshal(req)
	if err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}

	logrus.Debug("POST /auth/card: ", string(data))
	rdr := bytes.NewReader(data)

	// Http client with X-Local-Token
	client := &http.Client{}
	request, err := http.NewRequest("POST", conf.ApiURL+"/auth/card", rdr)
	if err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}

	request.Header.Add("X-Local-Token", conf.LocalToken)
	request.Header.Add("Content-Type", "application/json")

	// post the same request to API_URL
	r, err := client.Do(request)
	if err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}

	if r.StatusCode != http.StatusOK {
		logrus.Debug("/auth/card returned ", r.StatusCode)
		return c.JSON(400, gin.H{"error": "invalid card"})
	}

	request, err = http.NewRequest("GET", conf.ApiURL+"/account/admin", nil)
	if err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}

	request.Header.Add("X-Local-Token", conf.LocalToken)
	for _, cookie := range r.Cookies() {
		request.AddCookie(cookie)
	}

	r2, err := client.Do(request)
	if err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}

	if r2.StatusCode != http.StatusOK {
		logrus.Debug("/account/admin returned ", r2.StatusCode)
		return c.JSON(400, gin.H{"error": "invalid card"})
	}

	type Resp struct {
		IsAllowed bool `json:"is_allowed"`
	}

	var resp Resp
	if err := json.NewDecoder(r2.Body).Decode(&resp); err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}

	if !resp.IsAllowed {
		return c.JSON(400, gin.H{"error": "invalid card"})
	}
	return nil
}

func sendToACM(c echo.Context, data []byte) error {
	// Write 1 to ttyACM*
	files := os.DirFS("/dev")
	err := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(path, "ttyACM") {
			f, err := os.OpenFile("/dev/"+path, os.O_RDWR, 0755)
			if err != nil {
				return err
			}

			f.Write([]byte(data))
			f.Close()
		}

		return nil
	})
	if err != nil {
		return c.JSON(400, gin.H{"error": err.Error()})
	}
	return nil
}

func routes(e *echo.Echo) {
	e.POST("/porte", func(c echo.Context) error {
		if err := checkRequest(c); err != nil {
			return err
		}

		return sendToACM(c, []byte("1"))
	})

	e.POST("/caisse", func(c echo.Context) error {
		if err := checkRequest(c); err != nil {
			return err
		}

		return sendToACM(c, []byte("2"))
	})

	e.POST("/ventilo", func(c echo.Context) error {
		if err := checkRequest(c); err != nil {
			return err
		}

		return sendToACM(c, []byte("3"))
	})
}
