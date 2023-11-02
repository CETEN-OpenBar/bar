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

		// post the same request to API_URL
		r, err := http.Post(conf.ApiURL, "application/json", rdr)
		if err != nil {
			return c.JSON(400, gin.H{"error": err.Error()})
		}

		if r.StatusCode != http.StatusOK {
			return c.JSON(400, gin.H{"error": "invalid card"})
		}

		r2, err := http.Post("http://localhost:8080/account/admin", "application/json", rdr)
		if err != nil {
			return c.JSON(400, gin.H{"error": err.Error()})
		}

		if r2.StatusCode != http.StatusOK {
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
