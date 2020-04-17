package domain

import (
	"strings"

	gouuid "github.com/satori/go.uuid"
)

func uuid() string {
	return strings.Replace(gouuid.NewV4().String(), "-", "", -1)
}
