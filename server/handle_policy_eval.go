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

type evalCmd struct {
	Policy string `json:"policy" binding:"required"`
	Input  string `json:"input" binding:"required"`
}

type evalCmdList struct {
	Input      string   `json:"input" binding:"required"`
	PolicyList []string `json:"policy_list" binding:"required"`
}

type evalResult struct {
	Input        string `json:"input,omitempty" binding:"required"`
	Policy       string `json:"policy,omitempty" binding:"required"`
	IsSuccessful bool   `json:"isSuccessful" binding:"required"`
}

// @Tags			Policies Evaluate
// @Summary		Evaluate the policies
// @description	Evaluate whether the input meets all the policies
// @Accept			json
// @Produce		json
// @Param			cmdList	body	evalCmdList	true	"evalCmdList"
// @Router			/evaluate/policies [POST]
// @Success		200	{object}	evalResult	"OK"
// @Failure		400	{object}	ErrorResult	"Bad Request"
// @Failure		500	{object}	ErrorResult	"Internal Server Error"
// @Failure		502	{object}	ErrorResult	"Bad Gateway"
func policiesEvaluateHandler(_ *config.Conf) gin.HandlerFunc {
	return func(c *gin.Context) {

		var cmdList evalCmdList

		if err := c.ShouldBindWith(&cmdList, binding.JSON); err != nil {
			log.Debug(err)
			abortWithError(c, http.StatusBadRequest, fmt.Sprintf("Missing required field. Details:%v", err))
			return
		}

		for _, policy := range cmdList.PolicyList {
			decision, err := policyQuery(policy, cmdList.Input)
			if err != nil {
				abortWithError(c, http.StatusInternalServerError, err.Error())
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

		c.JSON(http.StatusOK, evalResult{
			IsSuccessful: true,
		})
	}
}

// @Tags			Policies Evaluate
// @Summary		Evaluate the policy
// @description	Evaluate whether the input meets the policy
// @Accept			json
// @Produce		json
// @Router			/evaluate/policy [POST]
// @Param			cmd	body		evalCmd		true	"evalCmd"
// @Success		200	{object}	evalResult	"OK"
// @Failure		400	{object}	ErrorResult	"Bad Request"
// @Failure		500	{object}	ErrorResult	"Internal Server Error"
// @Failure		502	{object}	ErrorResult	"Bad Gateway"
func policyEvaluateHandler(_ *config.Conf) gin.HandlerFunc {
	return func(c *gin.Context) {

		var cmd evalCmd

		if err := c.ShouldBindWith(&cmd, binding.JSON); err != nil {
			log.Debug(err)
			abortWithError(c, http.StatusBadRequest, fmt.Sprintf("Missing required field. Details:%v", err))
			return
		}

		decision, err := policyQuery(cmd.Policy, cmd.Input)
		if err != nil {
			abortWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, evalResult{
			IsSuccessful: decision,
		})
	}
}

func policyQuery(policyRego string, input interface{}) (decision bool, err error) {

	policyRegoEx := handlePackageName(policyRego)
	policyQuery := "data.policyman.auth"
	return PolicyEval(policyRegoEx, policyQuery, input)
}

func PolicyEval(policyRego string, policyQuery string, input interface{}) (decision bool, err error) {

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
		return false, fmt.Errorf("prepare for rego failed as fllow:\n %v", err)
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
