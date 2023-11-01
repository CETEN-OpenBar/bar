package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/v56/github"
	"github.com/sirupsen/logrus"
)

func init() {
	// Get checksum of current version
	localVersion, err := os.ReadFile("current.md5")
	if err != nil {
		logrus.WithError(err).Error("failed to open current.md5")
	}

	client := github.NewClient(nil)

	// list all organizations for user "willnorris"
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "yyewolf", "bar")
	if err != nil {
		logrus.WithError(err).Error("failed to get latest release")
		return
	}

	var md5 string
	var tarball string

	for _, asset := range release.Assets {
		name := asset.GetName()
		if !strings.HasPrefix(name, "nfc") {
			continue
		}
		if strings.HasSuffix(name, "linux-amd64.tar.gz.md5") {
			md5 = asset.GetBrowserDownloadURL()
			logrus.WithField("url", md5).Info("found md5")
		}
		if strings.HasSuffix(name, "linux-amd64.tar.gz") {
			tarball = asset.GetBrowserDownloadURL()
		}
	}

	// Download md5
	r, err := http.Get(md5)
	if err != nil {
		logrus.WithError(err).Error("failed to get checksum")
		return
	}

	if r.StatusCode != http.StatusOK {
		logrus.WithField("status", r.StatusCode).Error("failed to get checksum")
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Error("failed to read checksum")
		return
	}

	localChecksum := string(localVersion)
	remoteChecksum := string(d)

	logrus.WithFields(logrus.Fields{
		"local":  localChecksum,
		"remote": remoteChecksum,
	}).Info("checksums")

	if localChecksum != remoteChecksum {
		performUpdate(remoteChecksum, tarball)
	}

	logrus.Info("checksum OK")
}

func performUpdate(checksum, tarball string) {
	logrus.Info("performing update")
	r, err := http.Get(tarball)
	if err != nil {
		logrus.WithError(err).Error("failed to get update")
		return
	}

	if r.StatusCode != http.StatusOK {
		logrus.WithField("status", r.StatusCode).Error("failed to get update")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Error("failed to read update")
		return
	}

	// verify checksum
	h := md5.New()
	_, err = h.Write(data)
	if err != nil {
		logrus.WithError(err).Error("failed to write update")
		return
	}

	if fmt.Sprintf("%x\n", h.Sum(nil)) != checksum {
		logrus.Error("checksum mismatch")
		return
	}

	// open a subprocess that will replace the current process file
	err = os.WriteFile("update.tar.gz", data, 0755)
	if err != nil {
		logrus.WithError(err).Error("failed to open process file")
		return
	}

	// write sh file that will finish the update
	err = os.WriteFile("update.sh", []byte(fmt.Sprintf(`#!/bin/sh
sleep 1s
rm nfc
tar -xzf update.tar.gz
chmod +x %s
rm update.tar.gz
rm LICENSE
rm update.sh
`, os.Args[0])), 0755)
	if err != nil {
		logrus.WithError(err).Error("failed to open process file")
		return
	}

	// chmod the sh file
	err = os.Chmod("update.sh", 0755)
	if err != nil {
		logrus.WithError(err).Error("failed to chmod sh file")
		return
	}

	// run the sh file in another process
	exec.Command("sh", "update.sh").Start()

	// save the checksum in current.md5
	err = os.WriteFile("current.md5", []byte(checksum), 0755)
	if err != nil {
		logrus.WithError(err).Error("failed to open process file")
		return
	}

	logrus.Info("update OK")
	logrus.Info("restarting")
	os.Exit(0)
}
