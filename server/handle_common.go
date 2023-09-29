/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package server

import "github.com/gin-gonic/gin"

type ErrorResult struct {
	ErrMsg string `json:"err_msg,omitempty" binding:"required"`
} // @name ErrorResult

func abortWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, ErrorResult{
		ErrMsg: message,
	})
}
