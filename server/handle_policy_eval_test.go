/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPolicyEval(t *testing.T) {
	rego := `
package example.auth

import future.keywords.if
import future.keywords.in

default allow := false

allow if {
    input.method == "GET"
    input.path == ["salary", input.subject.user]
}

allow if is_admin

is_admin if "admin" in input.subject.groups
`
	query := "data.example.auth"

	input := map[string]interface{}{
		"method": "GET",
		"path":   []interface{}{"salary", "bob"},
		"subject": map[string]interface{}{
			"user":   "bob",
			"groups": []interface{}{"sales", "marketing"},
		},
	}

	decision, err := policyEval(rego, query, input)
	assert.Nil(t, err)
	assert.Equal(t, true, decision)

}
