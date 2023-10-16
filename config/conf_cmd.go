/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Logo string = `
  ___  ___  _    ___ _____   __        __  __   _   _  _
 | _ \/ _ \| |  |_ _/ __\ \ / /  ___  |  \/  | /_\ | \| |
 |  _/ (_) | |__ | | (__ \ V /  |___| | |\/| |/ _ \| .' |
 |_|  \___/|____|___\___| |_|         |_|  |_/_/ \_\_|\_|


`

var RootCmd = &cobra.Command{
	Use:  "policy-man",
	Long: Logo,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(Logo)
	},
}

func init() {
	RootCmd.Flags().StringP("mode", "m", "release", "The mode of the HTTP server.[release/debug/test]")

	RootCmd.Flags().StringP("config", "c", "", "config file (default is ./config.yml)")

	RootCmd.Flags().StringP("host", "a", "", "The host of the HTTP server")
	RootCmd.Flags().StringP("port", "p", "", "The port of the HTTP server")

	RootCmd.Flags().String("log.level", "warn", "The level of the log")
	RootCmd.Flags().String("log.path", "stdout", "The path of the log")

	err := viper.BindPFlags(RootCmd.Flags())
	if err != nil {
		return
	}
}
