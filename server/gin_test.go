/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package server_test

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
)

func GinGetTestCtx(method string, body string) (c *gin.Context, r *httptest.ResponseRecorder, err error) {

	if method == "POST" {
		return GinGetTestCtxPost(body)
	}

	return nil, nil, errors.New("not supported")
}

func GinGetTestCtxPost(body string) (c *gin.Context, r *httptest.ResponseRecorder, err error) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	ctx.Request.Method = "POST"

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(body)))

	return ctx, w, nil
}
