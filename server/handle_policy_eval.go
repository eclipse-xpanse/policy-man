/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eclipse-xpanse/policy-man/config"
	"github.com/eclipse-xpanse/policy-man/log"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/open-policy-agent/opa/rego"
	"net/http"
	"runtime/debug"
)

type EvalRego struct {
	Policy  string `json:"policy" binding:"required"`
	IsAllow bool   `json:"isAllow" binding:"required"`
}

type EvalCmd struct {
	Rego  EvalRego `json:"rego" binding:"required"`
	Input string   `json:"input" binding:"required"`
}

type EvalCmdList struct {
	Input    string     `json:"input" binding:"required"`
	RegoList []EvalRego `json:"rego_list" binding:"required"`
}

func abortWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"healthStatus": "OK",
	})
	return
}

func policiesEvaluateHandler(_ *config.Conf) gin.HandlerFunc {
	return func(c *gin.Context) {

		var cmdList EvalCmdList

		if err := c.ShouldBindWith(&cmdList, binding.JSON); err != nil {
			log.Debug(err)
			abortWithError(c, http.StatusBadRequest, "Missing eval cmd field.")
			return
		}

		for _, rego := range cmdList.RegoList {
			decision, err := policyQuery(rego.Policy, rego.IsAllow, cmdList.Input)
			if err != nil {
				abortWithError(c, 500, err.Error())
				return
			}
			if rego.IsAllow && !decision || !rego.IsAllow && decision {
				c.JSON(200, gin.H{
					"isSuccessful": false,
				})
			}
		}

		c.JSON(200, gin.H{
			"isSuccessful": true,
		})
		return
	}
}

func policyEvaluateHandler(_ *config.Conf) gin.HandlerFunc {
	return func(c *gin.Context) {

		var cmd EvalCmd

		if err := c.ShouldBindWith(&cmd, binding.JSON); err != nil {
			log.Debug(err)
			abortWithError(c, http.StatusBadRequest, "Missing eval cmd field.")
			return
		}

		decision, err := policyQuery(cmd.Rego.Policy, cmd.Rego.IsAllow, cmd.Input)
		if err != nil {
			abortWithError(c, 500, err.Error())
			return
		}

		c.JSON(200, gin.H{
			"isAllow":      cmd.Rego.IsAllow,
			"isSuccessful": decision,
		})
		return
	}
}

func policyQuery(policyRego string, isAllow bool, input interface{}) (decision bool, err error) {

	policyRegoEx := fmt.Sprintf("package example.auth\n\n%v", policyRego)
	var policyQuery string
	if isAllow {
		policyQuery = fmt.Sprintf("data.example.auth.allow")
	} else {
		policyQuery = fmt.Sprintf("data.example.auth.deny")
	}

	return policyEval(policyRegoEx, policyQuery, input)
}

func policyEval(policyRego string, policyQuery string, input interface{}) (decision bool, err error) {

	defer func() {
		if r := recover(); r != nil {
			stackInfo := string(debug.Stack())
			err = errors.New(stackInfo)
		}
	}()

	ctx := context.TODO()

	query, err := rego.New(
		rego.Query(fmt.Sprintf("result = %s", policyQuery)),
		rego.Module("policy-man.rego", policyRego),
	).PrepareForEval(ctx)
	if err != nil {
		return false, errors.New("prepare for eval failed")
	}

	var inputMap map[string]interface{}

	if m, ok := input.(map[string]any); ok {
		inputMap = m
	} else {
		err = json.Unmarshal([]byte(input.(string)), &inputMap)
		if err != nil {
			return false, err
		}
	}

	results, err := query.Eval(ctx, rego.EvalInput(inputMap))

	if err != nil {
		return false, err
	} else if len(results) == 0 {
		return false, nil
	} else {
		decision, ok := results[0].Bindings["result"].(bool)
		if !ok || !decision {
			return false, nil
		}
		return decision, nil
	}
}
