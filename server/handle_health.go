/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type healthStatus string

const (
	healthOK  healthStatus = "OK"
	healthNOK healthStatus = "NOK"
)

type systemStatus struct {
	HealthStatus healthStatus `json:"healthStatus"`
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, systemStatus{HealthStatus: healthOK})
}
