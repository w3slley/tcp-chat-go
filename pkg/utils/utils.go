package utils

import "strings"

func GetCommandFromMessage(message string) string {
	cmd := strings.Split(message, " ")[0]
	return strings.TrimSpace(cmd)
}

func GetCommandArgument(message string, cmd string) string {
	return strings.TrimSuffix(strings.TrimPrefix(message, cmd+" "), "\n")
}
