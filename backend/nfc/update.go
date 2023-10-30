package main

import (
	"crypto/sha256"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	// Get checksum of current process file
	h := sha256.New()
	filename := os.Args[0]
	if _, err := os.Stat(filename); err != nil {
		logrus.WithError(err).Fatal("failed to stat process file")
	}

	f, err := os.Open(filename)
	if err != nil {
		logrus.WithError(err).Fatal("failed to open process file")
	}

	if _, err := io.Copy(h, f); err != nil {
		logrus.WithError(err).Fatal("failed to copy process file")
	}

	// Get checksum from https://api.bar.telecomnancy.net/nfc/checksum
	// and compare
	r, err := http.Get("https://api.bar.telecomnancy.net/nfc/checksum")
	if err != nil {
		logrus.WithError(err).Fatal("failed to get checksum")
	}

	if r.StatusCode != http.StatusOK {
		logrus.WithField("status", r.StatusCode).Fatal("failed to get checksum")
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Fatal("failed to read checksum")
	}

	if string(d) != string(h.Sum(nil)) {
		performUpdate()
		return
	}

	logrus.Info("checksum OK")
}

func performUpdate() {
	logrus.Info("performing update")
	r, err := http.Get("https://api.bar.telecomnancy.net/nfc")
	if err != nil {
		logrus.WithError(err).Fatal("failed to get update")
	}

	if r.StatusCode != http.StatusOK {
		logrus.WithField("status", r.StatusCode).Fatal("failed to get update")
	}

	f, err := os.OpenFile(os.Args[0], os.O_RDWR, 0755)
	if err != nil {
		logrus.WithError(err).Fatal("failed to open process file")
	}

	if _, err := io.Copy(f, r.Body); err != nil {
		logrus.WithError(err).Fatal("failed to copy process file")
	}

	logrus.Info("update OK")
	logrus.Info("restarting")
	os.Exit(0)
}
