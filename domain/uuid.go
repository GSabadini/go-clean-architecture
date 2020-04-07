package domain

import (
	gouuid "github.com/satori/go.uuid"
)

func uuid() string {
	return gouuid.NewV4().String()
}
