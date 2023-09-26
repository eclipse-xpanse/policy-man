/*
 * SPDX-License-Identifier: Apache-2.0
 * SPDX-FileCopyrightText: Huawei Inc.
 */

package server_test

import (
	"encoding/json"
	"github.com/eclipse-xpanse/policy-man/server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPolicyValidate(t *testing.T) {

	jsonRequest := `
{
    "policy_list": [
            "packagsse examplse.auth\n\n\nimport future.keywords.if\nimport future.keywords.in\n\ndefault allow := false\n\nallow if {\n    input.method == \"GET\"\n    input.path == [\"salary\", input.subject.user]\n}\n\nallow if is_admin\n\nis_admin if \"admin\" in input.subject.groups",
            "import future.keywords.if\nimport future.keywords.in\n\ndefault deny := false\n\nallow if {\n    input.method == \"GET\"\n    input.path == [\"salary\", input.subject.user]\n}\n\nallow if is_admin\n\nis_admin if \"admin\" in input.subject.groupsss"
    ]
}
`
	ctx, w, err := GinGetTestCtx("POST", jsonRequest)
	assert.Nil(t, nil, err)
	ctx.Request.Header.Set("Content-Type", "application/json")

	server.PoliciesValidateHandler(nil)(ctx)
	assert.Equal(t, http.StatusOK, w.Code)

	rsp := server.ValidateResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &rsp)
	assert.Nil(t, nil, err)
	assert.Equal(t, false, rsp.IsSuccessful)
}
