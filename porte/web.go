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

func routes(e *echo.Echo) {
	conf := GetConfig()

	e.POST("/open", func(c echo.Context) error {
		var req OpenRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, gin.H{"error": err.Error()})
		}

		data, err := json.Marshal(req)
		if err != nil {
			return c.JSON(400, gin.H{"error": err.Error()})
		}

		rdr := bytes.NewReader(data)

		// Http client with X-Local-Token
		client := &http.Client{}
		request, err := http.NewRequest("POST", conf.ApiURL+"/auth/card", rdr)
		if err != nil {
			return c.JSON(400, gin.H{"error": err.Error()})
		}

		request.Header.Add("X-Local-Token", conf.LocalToken)

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

		r2, err := client.Do(request)
		if err != nil {
			logrus.Debug("/account/admin returned ", r.StatusCode)
			return c.JSON(400, gin.H{"error": err.Error()})
		}

		if r2.StatusCode != http.StatusOK {
			logrus.Debug("/account/admin returned ", r.StatusCode)
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

		// Write 1 to ttyACM*
		files := os.DirFS("/dev")
		err = fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if strings.HasPrefix(path, "ttyACM") {
				f, err := os.OpenFile("/dev/"+path, os.O_RDWR, 0755)
				if err != nil {
					return err
				}

				f.Write([]byte("1"))
				f.Close()
			}

			return nil
		})
		if err != nil {
			return c.JSON(400, gin.H{"error": err.Error()})
		}

		return nil
	})
}