// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// Args provides plugin execution arguments.
type Args struct {

	// Level defines the plugin log level.
	Level                   string `envconfig:"PLUGIN_LOG_LEVEL"`
	Username                string `envconfig:"PLUGIN_USERNAME"`
	Password                string `envconfig:"PLUGIN_PASSWORD"`
	APIKey                  string `envconfig:"PLUGIN_API_KEY"`
	AccessToken             string `envconfig:"PLUGIN_ACCESS_TOKEN"`
	URL                     string `envconfig:"PLUGIN_URL"`
	IncludeDirs             string `envconfig:"PLUGIN_INCLUDE_DIRS"`
	SourcePath              string `envconfig:"PLUGIN_SOURCE_PATH"`
	TargetPath              string `envconfig:"PLUGIN_TARGET_PATH"`
	Insecure                string `envconfig:"PLUGIN_INSECURE"`
	Retries                 int    `envconfig:"PLUGIN_RETRIES"`
	Spec                    string `envconfig:"PLUGIN_SPEC"`
	Threads                 int    `envconfig:"PLUGIN_THREADS"`
	SpecVars                string `envconfig:"PLUGIN_SPEC_VARS"`
	Flat                    string `envconfig:"PLUGIN_FLAT"`
	ServerId                string `envconfig:"PLUGIN_SERVER_ID"`
	BuildName               string `envconfig:"PLUGIN_BUILD_NAME"`
	BuildNumber             string `envconfig:"PLUGIN_BUILD_NUMBER"`
	Project                 string `envconfig:"PLUGIN_PROJECT"`
	Module                  string `envconfig:"PLUGIN_MODULE"`
	Props                   string `envconfig:"PLUGIN_PROPS"`
	ExcludeProps            string `envconfig:"PLUGIN_EXCLUDE_PROPS"`
	Build                   string `envconfig:"PLUGIN_BUILD"`
	Bundle                  string `envconfig:"PLUGIN_BUNDLE"`
	RetryWaitTime           string `envconfig:"PLUGIN_RETRY_WAIT_TIME"`
	SyncDeletes             string `envconfig:"PLUGIN_SYNC_DELETES"`
	SortBy                  string `envconfig:"PLUGIN_SORT_BY"`
	SortOrder               string `envconfig:"PLUGIN_SORT_ORDER"`
	GPGKey                  string `envconfig:"PLUGIN_GPG_KEY"`
	Exclusions              string `envconfig:"PLUGIN_EXCLUSIONS"`
	DetailedSummary         string `envconfig:"PLUGIN_DETAILED_SUMMARY"`
	Recursive               string `envconfig:"PLUGIN_RECURSIVE"`
	DryRun                  string `envconfig:"PLUGIN_DRY_RUN"`
	Explode                 string `envconfig:"PLUGIN_EXPLODE"`
	ByPassArchiveInspection string `envconfig:"PLUGIN_BY_PASS_ARCHIVE_INSPECTION"`
	ValidationSymLinks      string `envconfig:"PLUGIN_VALIDATION_SYS_LINKS"`
	Quiet                   string `envconfig:"PLUGIN_QUIET"`
	FailNoOp                string `envconfig:"PLUGIN_FAIL_NO_OP"`
	SplitCount              int    `envconfig:"PLUGIN_SPLIT_COUNT"`
	MinSplit                int    `envconfig:"PLUGIN_MIN_SPLIT"`
	Limit                   int    `envconfig:"PLUGIN_LIMIT"`
	Offset                  int    `envconfig:"PLUGIN_OFFSET"`
}

// Exec executes the plugin.
func Exec(ctx context.Context, args Args) error {
	// write code here
	if args.URL == "" {
		return fmt.Errorf("url needs to be set")
	}

	if args.SourcePath == "" {
		return fmt.Errorf("download path needs to be set")
	}

	cmdArgs := []string{getJfrogBin(), "rt", "dl", fmt.Sprintf("--url %s", args.URL)}

	if args.Retries != 0 {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--retries=%d", args.Retries))
	}

	// Set authentication params
	cmdArgs, error := setAuthParams(cmdArgs, args)
	if error != nil {
		return error
	}

	optionalStringParams := []struct {
		flag  string
		value string
	}{
		{"--server-id", args.ServerId},
		{"--build-name", args.BuildName},
		{"--build-number", args.BuildNumber},
		{"--project", args.Project},
		{"--module", args.Module},
		{"--props", args.Props},
		{"--exclude-props", args.ExcludeProps},
		{"--build", args.Build},
		{"--bundle", args.Bundle},
		{"--retry-wait-time", args.RetryWaitTime},
		{"--sync-deletes", args.SyncDeletes},
		{"--sort-by", args.SortBy},
		{"--sort-order", args.SortOrder},
		{"--gpg-key", args.GPGKey},
		{"--exclusions", args.Exclusions},
	}

	for _, param := range optionalStringParams {
		if param.value != "" {
			cmdArgs = append(cmdArgs, fmt.Sprintf("%s=%s", param.flag, param.value))
		}
	}

	optionalBooleanParams := []struct {
		flag         string
		defaultValue bool
		value        string
	}{
		{"--insecure-tls", false, args.Insecure},
		{"--flat", false, args.Flat},
		{"--detailed-summary", false, args.DetailedSummary},
		{"--recursive", true, args.Recursive},
		{"--dry-run", false, args.DryRun},
		{"--explode", false, args.Explode},
		{"--bypass-archive-inspection", false, args.ByPassArchiveInspection},
		{"--validate-symlinks", false, args.ValidationSymLinks},
		{"--include-dirs", false, args.IncludeDirs},
		{"--quiet", false, args.Quiet},
		{"--fail-no-op", false, args.FailNoOp},
	}

	for _, param := range optionalBooleanParams {
		val := parseBoolOrDefault(param.defaultValue, param.value)
		cmdArgs = append(cmdArgs, fmt.Sprintf("%s=%s", param.flag, strconv.FormatBool(val)))
	}

	optionalIntegerParams := []struct {
		flag  string
		value int
	}{
		{"--retries", args.Retries},
		{"--threads", args.Threads},
		{"--split-count", args.SplitCount},
		{"--min-split", args.MinSplit},
		{"--limit", args.Limit},
		{"--offset", args.Offset},
	}

	for _, param := range optionalIntegerParams {
		if param.value != 0 {
			cmdArgs = append(cmdArgs, fmt.Sprintf("%s=%d", param.flag, param.value))
		}
	}

	// Take in spec file or use source/target arguments
	if args.Spec != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--spec=%s", args.Spec))
		if args.SpecVars != "" {
			cmdArgs = append(cmdArgs, fmt.Sprintf("--spec-vars='%s'", args.SpecVars))
		}
	}

	//source path
	cmdArgs = append(cmdArgs, args.SourcePath)

	//target path
	if args.TargetPath != "" {
		cmdArgs = append(cmdArgs, args.TargetPath)
	}

	cmdStr := strings.Join(cmdArgs[:], " ")

	shell, shArg := getShell()

	cmd := exec.Command(shell, shArg, cmdStr)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "JFROG_CLI_OFFER_CONFIG=false", "JFROG_CLI_TRANSITIVE_DOWNLOAD_EXPERIMENTAL=true")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)

	err := cmd.Run()

	destination := parseBoolOrDefault(false, args.TargetPath)
	if destination {
		cmdArgs = append(cmdArgs, "--include-dirs")
	}

	return err
}

// setAuthParams appends authentication parameters to cmdArgs based on the provided credentials.
func setAuthParams(cmdArgs []string, args Args) ([]string, error) {
	// Set authentication params
	envPrefix := getEnvPrefix()
	if args.Username != "" && args.Password != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--user %sPLUGIN_USERNAME", envPrefix))
		cmdArgs = append(cmdArgs, fmt.Sprintf("--password %sPLUGIN_PASSWORD", envPrefix))
	} else if args.APIKey != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--apikey %sPLUGIN_API_KEY", envPrefix))
	} else if args.AccessToken != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--access-token %sPLUGIN_ACCESS_TOKEN", envPrefix))
	} else {
		return nil, fmt.Errorf("either username/password, api key or access token needs to be set")
	}
	return cmdArgs, nil
}

func getShell() (string, string) {
	if runtime.GOOS == "windows" {
		return "powershell", "-Command"
	}
	return "sh", "-c"
}

func getJfrogBin() string {
	if runtime.GOOS == "windows" {
		return "C:/bin/jfrog.exe"
	}
	return "jfrog"
}

func getEnvPrefix() string {
	if runtime.GOOS == "windows" {
		return "$Env:"
	}
	return "$"
}

func parseBoolOrDefault(defaultValue bool, s string) (result bool) {
	var err error
	result, err = strconv.ParseBool(s)
	if err != nil {
		result = defaultValue
	}

	return
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}
