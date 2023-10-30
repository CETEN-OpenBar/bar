package api

import (
	"bar/autogen"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var storedChecksum string

// (GET /nfc)
func (s *Server) GetNFC(c echo.Context) error {
	f, err := os.Open("/nfc")
	if err != nil {
		logrus.WithError(err).Fatal("failed to open process file")
	}

	return autogen.GetNFC200ApplicationoctetStreamResponse{Body: f}.VisitGetNFCResponse(c.Response())
}

// (GET /nfc/checksum)
func (s *Server) GetNFCChecksum(c echo.Context) error {
	if storedChecksum != "" {
		return autogen.GetNFCChecksum200JSONResponse{Checksum: storedChecksum}.VisitGetNFCChecksumResponse(c.Response())
	}

	h := sha256.New()

	f, err := os.Open("/nfc")
	if err != nil {
		logrus.WithError(err).Fatal("failed to open process file")
	}

	if _, err := io.Copy(h, f); err != nil {
		logrus.WithError(err).Fatal("failed to copy process file")
	}

	storedChecksum = fmt.Sprintf("%x", h.Sum(nil))

	return autogen.GetNFCChecksum200JSONResponse{Checksum: storedChecksum}.VisitGetNFCChecksumResponse(c.Response())
}
