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
	IsSuccessful bool   `json:"isSuccessful" binding:"required"`
	Policy       string `json:"policy,omitempty" binding:"required"`
	ErrMsg       string `json:"err_msg,omitempty" binding:"required"`
}

// PoliciesValidateHandler
// @Tags			Policies Validate
// @Summary		Validate the policies
// @description	Validate the policies
// @Accept			json
// @Produce		json
// @Router			/validate/policies [POST]
// @Param			policyList	body		validatePolicyList	true	"policyList"
// @Success		200			{object}	ValidateResponse	"OK"
// @Failure		400			{object}	ErrorResult			"Bad Request"
// @Failure		500			{object}	ErrorResult			"Internal Server Error"
// @Failure		502			{object}	ErrorResult			"Bad Gateway"
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
