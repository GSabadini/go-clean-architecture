package domain

import (
	"strings"

	gouuid "github.com/satori/go.uuid"
)

func uuid() string {
	return strings.Replace(gouuid.NewV4().String(), "-", "", -1)
}

//IsValidUUID
func IsValidUUID(uuid string) bool {
	_, err := gouuid.FromString(uuid)
	return err == nil
}
