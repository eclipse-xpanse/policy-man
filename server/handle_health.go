/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type healthStatus string // @name HealthStatus

const (
	healthOK  healthStatus = "OK"
	healthNOK healthStatus = "NOK"
)

type systemStatus struct {
	HealthStatus healthStatus `json:"healthStatus" binding:"required"`
} // @name SystemStatus

// @Tags			Admin
// @Summary		Check health
// @description	Check health status of service
// @Accept			json
// @Produce		json
// @Router			/health [GET]
// @Success		200	{object}	systemStatus	"OK"
// @Failure		400	{object}	ErrorResult		"Bad Request"
// @Failure		500	{object}	ErrorResult		"Internal Server Error"
// @Failure		502	{object}	ErrorResult		"Bad Gateway"
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, systemStatus{HealthStatus: healthOK})
}
