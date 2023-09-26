/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package server

import (
	"fmt"
	"github.com/eclipse-xpanse/policy-man/config"
	"github.com/eclipse-xpanse/policy-man/log"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type validatePolicyList struct {
	PolicyList []string `json:"policy_list" binding:"required"`
}

type ValidateResponse struct {
	IsSuccessful bool   `json:"isSuccessful"`
	Policy       string `json:"policy,omitempty"`
	ErrMsg       string `json:"err_msg,omitempty"`
}

func PoliciesValidateHandler(_ *config.Conf) gin.HandlerFunc {
	return func(c *gin.Context) {

		var policyList validatePolicyList

		if err := c.ShouldBindWith(&policyList, binding.JSON); err != nil {
			log.Debug(err)
			abortWithError(c, http.StatusBadRequest, fmt.Sprintf("Missing required field. %v", err))
			return
		}

		for _, policy := range policyList.PolicyList {
			_, err := policyQuery(policy, map[string]any{})
			if err != nil {
				c.JSON(http.StatusOK, ValidateResponse{
					IsSuccessful: false,
					Policy:       policy,
					ErrMsg:       err.Error(),
				})
				return
			}
		}

		c.JSON(http.StatusOK, ValidateResponse{
			IsSuccessful: true,
		})
	}
}
