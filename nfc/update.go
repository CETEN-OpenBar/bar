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
	// Get checksum of current process file
	h := md5.New()
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

	client := github.NewClient(nil)

	// list all organizations for user "willnorris"
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "yyewolf", "bar")
	if err != nil {
		logrus.WithError(err).Fatal("failed to get latest release")
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
		logrus.WithError(err).Fatal("failed to get checksum")
	}

	if r.StatusCode != http.StatusOK {
		logrus.WithField("status", r.StatusCode).Fatal("failed to get checksum")
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.WithError(err).Fatal("failed to read checksum")
	}

	localChecksum := fmt.Sprintf("%x\n", h.Sum(nil))
	remoteChecksum := string(d)

	logrus.WithFields(logrus.Fields{
		"local":  localChecksum,
		"remote": remoteChecksum,
	}).Info("checksums")

	if localChecksum != remoteChecksum {
		performUpdate(tarball)
	}

	logrus.Info("checksum OK")
}

func performUpdate(tarball string) {
	logrus.Info("performing update")
	r, err := http.Get(tarball)
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
	err = os.WriteFile("update.tar.gz", data, 0755)
	if err != nil {
		logrus.WithError(err).Fatal("failed to open process file")
	}

	// write sh file that will finish the update
	err = os.WriteFile("update.sh", []byte(fmt.Sprintf(`#!/bin/sh
sleep 1s
rm nfc
tar -xzf update.tar.gz
chmod +x %s
rm update.tar.gz
`, os.Args[0])), 0755)
	if err != nil {
		logrus.WithError(err).Fatal("failed to open process file")
	}

	// chmod the sh file
	err = os.Chmod("update.sh", 0755)
	if err != nil {
		logrus.WithError(err).Fatal("failed to chmod sh file")
	}

	// run the sh file in another process
	exec.Command("sh", "update.sh").Start()

	logrus.Info("update OK")
	logrus.Info("restarting")
	os.Exit(0)
}
