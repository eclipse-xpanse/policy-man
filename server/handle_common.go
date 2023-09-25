/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package server

import "github.com/gin-gonic/gin"

type ErrorResult struct {
	ErrCode int    `json:"err_code,omitempty"`
	ErrMsg  string `json:"err_msg,omitempty"`
}

func abortWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, ErrorResult{
		ErrCode: code,
		ErrMsg:  message,
	})
}
