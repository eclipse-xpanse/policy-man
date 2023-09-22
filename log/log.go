/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package log

import (
	"errors"
	"github.com/mattn/go-isatty"
	"github.com/sirupsen/logrus"
	"os"
)

var isTerm bool

// nolint
func init() {
	isTerm = isatty.IsTerminal(os.Stdout.Fd())
}

var (
	// Basic is log server error log
	basic = logrus.New()
)

// InitLog use for initial log module
func InitLog(level, path string) error {
	var err error

	if !isTerm {
		basic.SetFormatter(&logrus.JSONFormatter{})
		basic.SetFormatter(&logrus.JSONFormatter{})
	} else {
		basic.Formatter = &logrus.TextFormatter{
			TimestampFormat: "2023/01/01 - 01:01:01",
			FullTimestamp:   true,
		}

		basic.Formatter = &logrus.TextFormatter{
			TimestampFormat: "2023/01/01 - 01:01:01",
			FullTimestamp:   true,
		}
	}

	if err = SetLogLevel(basic, level); err != nil {
		return errors.New("set log level error: " + err.Error())
	}

	if err = SetLogOut(basic, path); err != nil {
		return errors.New("set log path error: " + err.Error())
	}

	return nil
}

// SetLogOut provide log stdout and stderr output
func SetLogOut(log *logrus.Logger, outString string) error {
	switch outString {
	case "stdout":
		log.Out = os.Stdout
	case "stderr":
		log.Out = os.Stderr
	default:
		f, err := os.OpenFile(outString, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o600)
		if err != nil {
			return err
		}

		log.Out = f
	}

	return nil
}

// SetLogLevel set the log level, available levels: panic, fatal, error, warn, info and debug
func SetLogLevel(log *logrus.Logger, levelString string) error {
	level, err := logrus.ParseLevel(levelString)
	if err != nil {
		return err
	}

	log.Level = level

	return nil
}

func Error(args ...interface{}) {
	basic.Error(args...)
}

func Info(args ...interface{}) {
	basic.Info(args...)
}

func Debug(args ...interface{}) {
	basic.Debug(args...)
}

func Warn(args ...interface{}) {
	basic.Warn(args...)
}

func Fatal(args ...interface{}) {
	basic.Fatal(args...)
}

func Errorf(format string, args ...interface{}) {
	basic.Errorf(format, args...)
}

func Infof(format string, args ...interface{}) {
	basic.Infof(format, args...)
}
func Debugf(format string, args ...interface{}) {
	basic.Debugf(format, args...)
}
func Warnf(format string, args ...interface{}) {
	basic.Warnf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	basic.Fatalf(format, args...)
}
