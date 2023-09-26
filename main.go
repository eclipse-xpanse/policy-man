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
)

func main() {

	cfg, err := config.LoadConf()
	if err != nil {
		fmt.Print("load config failed")
	}
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
