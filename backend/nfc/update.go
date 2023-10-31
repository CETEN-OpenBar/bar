package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

var apiUrl = "https://api.bar.telecomnancy.net/nfc"

// var apiUrl = "http://localhost:8080/nfc"

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
	r, err := http.Get(fmt.Sprintf("%s/checksum", apiUrl))
	if err != nil {
		logrus.Error("failed to get checksum, cannot perform update")
		return
	}

	if r.StatusCode != http.StatusOK {
		logrus.WithField("status", r.StatusCode).Fatal("failed to get checksum")
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Fatal("failed to read checksum")
	}

	type response struct {
		Checksum string `json:"checksum"`
	}

	var resp response
	if err := json.Unmarshal(d, &resp); err != nil {
		logrus.WithError(err).Fatal("failed to unmarshal checksum")
	}

	checksum := fmt.Sprintf("%x", h.Sum(nil))

	if resp.Checksum != checksum {
		performUpdate()
	}

	logrus.Info("checksum OK")
}

func performUpdate() {
	logrus.Info("performing update")
	r, err := http.Get(apiUrl)
	if err != nil {
		logrus.WithError(err).Fatal("failed to get update")
	}

	if r.StatusCode != http.StatusOK {
		logrus.WithField("status", r.StatusCode).Fatal("failed to get update")
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Fatal("failed to read update")
	}

	// open a subprocess that will replace the current process file
	err = os.WriteFile("update", data, 0755)
	if err != nil {
		logrus.WithError(err).Fatal("failed to open process file")
	}

	// write sh file that will finish the update
	err = os.WriteFile("update.sh", []byte(fmt.Sprintf(`#!/bin/sh
sleep 1s
mv update %s
chmod +x %s
rm update.sh
rm update
`, os.Args[0], os.Args[0])), 0755)
	if err != nil {
		logrus.WithError(err).Fatal("failed to open process file")
	}

	// chmod the sh file
	err = os.Chmod("update.sh", 0755)

	// run the sh file in another process
	exec.Command("sh", "update.sh").Start()

	logrus.Info("update OK")
	logrus.Info("restarting")
	os.Exit(0)
}
