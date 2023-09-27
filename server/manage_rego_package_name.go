package server

import (
	"fmt"
	"strings"
)

func handlePackageName(policyRego string) string {
	policyRegoFixed := removePackageAtTheBeginning(policyRego)
	policyRegoEx := fmt.Sprintf("package policyman.auth\n\n%v", policyRegoFixed)
	return policyRegoEx
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
