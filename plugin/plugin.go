// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	ApiKey string `envconfig:"PLUGIN_NUGET_APIKEY"`
	NugetUri string `envconfig:"PLUGIN_NUGET_URI"`
	PackageLocation string `envconfig:"PLUGIN_PACKAGE_LOCATION"`
}

const globalNugetUri = "https://api.nuget.org/v3/index.json"

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ValidateAndSetArgs(args Args) (validatedArgs Args, err error) {
	if args.ApiKey == ""{
		err = errors.New("nuget api key must be set in settings")
	}
	if args.NugetUri == "" {
		args.NugetUri = globalNugetUri
	}
	if args.PackageLocation != "" && !FileExists(args.PackageLocation){
		err = errors.New("the package location: " + args.PackageLocation + " does not exist")
	}
	validatedArgs = args
	return args, err
}

func WalkPath(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Fatal(err)
		}
		if filepath.Ext(path) == ".nupkg" {
			*files = append(*files, path)
		}
		return nil
	}
}

func PushToNuget(file string, args Args) *exec.Cmd {
	cmd := exec.Command("dotnet", "nuget", "push", file, "--api-key", args.ApiKey, "--source", args.NugetUri, "--skip-duplicate")
	return cmd
}

func Exec(ctx context.Context, args Args) error {
	logrus.Debugln("Starting ...")

	args, err := ValidateAndSetArgs(args)
	if err != nil {
		logrus.Errorln("Issues with the parameters passed: ")
		return err
	}

	var files []string
	// checks if single package location was provided, if not push all.
	if args.PackageLocation == ""{
		root := "/drone/src"
		err = filepath.Walk(root, WalkPath(&files))
		if err != nil {
			logrus.Errorln(err)
		}
	} else {
		files = append(files, args.PackageLocation)
	}

	for _, file := range files {
		if file != "" {
			logrus.Debugln("Pushing package: " +  file)
			cmd := PushToNuget(file, args)
			output, err := cmd.Output()
			if err != nil {
				logrus.Errorln(string(output))
				return err
			}
		} else {
			logrus.Debugln("No packages to publish ...")
		}
	}
	logrus.Debugln("Finished ...")
	return nil
}
