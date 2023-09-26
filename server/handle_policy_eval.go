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
	"strings"
)

type evalCmd struct {
	Policy string `json:"policy" binding:"required"`
	Input  string `json:"input" binding:"required"`
}

type evalCmdList struct {
	Input      string   `json:"input" binding:"required"`
	PolicyList []string `json:"policy_list" binding:"required"`
}

type evalResult struct {
	Input        string `json:"input,omitempty"`
	Policy       string `json:"policy,omitempty"`
	IsSuccessful bool   `json:"isSuccessful"`
}

func policiesEvaluateHandler(_ *config.Conf) gin.HandlerFunc {
	return func(c *gin.Context) {

		var cmdList evalCmdList

		if err := c.ShouldBindWith(&cmdList, binding.JSON); err != nil {
			log.Debug(err)
			abortWithError(c, http.StatusBadRequest, "Missing eval cmd field.")
			return
		}

		for _, policy := range cmdList.PolicyList {
			decision, err := policyQuery(policy, cmdList.Input)
			if err != nil {
				abortWithError(c, 500, err.Error())
				return
			}
			if !decision {
				c.JSON(200, evalResult{
					IsSuccessful: false,
					Policy:       policy,
					Input:        cmdList.Input,
				})
				return
			}
		}

		c.JSON(200, evalResult{
			IsSuccessful: true,
		})
	}
}

func policyEvaluateHandler(_ *config.Conf) gin.HandlerFunc {
	return func(c *gin.Context) {

		var cmd evalCmd

		if err := c.ShouldBindWith(&cmd, binding.JSON); err != nil {
			log.Debug(err)
			abortWithError(c, http.StatusBadRequest, "Missing eval cmd field.")
			return
		}

		decision, err := policyQuery(cmd.Policy, cmd.Input)
		if err != nil {
			abortWithError(c, 500, err.Error())
			return
		}

		c.JSON(200, evalResult{
			IsSuccessful: decision,
		})
	}
}

func policyQuery(policyRego string, input interface{}) (decision bool, err error) {

	policyRegoFixed := removePackageAtTheBeginning(policyRego)
	policyRegoEx := fmt.Sprintf("package policyman.auth\n\n%v", policyRegoFixed)
	policyQuery := "data.policyman.auth"
	return policyEval(policyRegoEx, policyQuery, input)
}

func removePackageAtTheBeginning(input string) string {
	lines := strings.Split(input, "\n")
	var outputLines []string

	for _, line := range lines {
		// Strip the spaces
		line = strings.TrimSpace(line)

		// Skip the blank lines and the lines starts with `#`
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Remove the line start with "package"
		if len(outputLines) == 0 && strings.HasPrefix(line, "package") {
			continue
		}

		outputLines = append(outputLines, line)
	}

	result := strings.Join(outputLines, "\n")
	return result
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
		result, ok := results[0].Bindings["result"].(map[string]any)
		if !ok {
			return false, nil
		}
		if allow, ok := result["allow"]; ok {
			if allowBool, ok := allow.(bool); ok && !allowBool {
				return false, nil
			}
		}
		if deny, ok := result["deny"]; ok {
			if denyBool, ok := deny.(bool); ok && denyBool {
				return false, nil
			}
		}

	}

	return true, nil
}
