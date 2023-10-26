/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package main

import (
	"context"
	"fmt"
	"github.com/eclipse-xpanse/policy-man/config"
	"github.com/eclipse-xpanse/policy-man/log"
	"github.com/eclipse-xpanse/policy-man/server"
	"github.com/spf13/pflag"
)

// Add comments to describe server openAPI information
//
//	@title			OpenAPI of policy-man
//	@version		1.0
//	@description	OpenAPI of policy-man server
func main() {

	cfg, err := config.LoadConf()
	if err != nil {
		fmt.Print("loading config failed\n")
		return
	}

	if isFlagPassed("version") {
		fmt.Println(version)
		return
	}

	if isFlagPassed("help") {
		return
	}

	fmt.Print(config.Logo)

	if err = log.InitLog(cfg.Log.Level, cfg.Log.Path); err != nil {
		return
	}

	ctx := context.Background()

	go func() {
		err := server.RunHTTPServer(ctx, cfg)
		if err != nil {
			log.Errorf("run http server failed: %v", err)
			return
		}
	}()

	<-ctx.Done()
}

func isFlagPassed(name string) bool {
	found := false
	config.RootCmd.Flags().Visit(func(flag *pflag.Flag) {
		if flag.Name == name {
			found = true
			return
		}
	})
	return found
}
