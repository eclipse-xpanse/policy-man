/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package server

import "github.com/gin-gonic/gin"

type HealthStatus string

const (
	healthOk HealthStatus = "OK"
)

type HealthStatusResponse struct {
	HealthStatus HealthStatus `json:"healthStatus"`
}

func healthHandler(c *gin.Context) {
	c.JSON(200, HealthStatusResponse{
		HealthStatus: healthOk,
	})
	return
}
