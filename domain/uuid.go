package domain

import (
	gouuid "github.com/satori/go.uuid"
)

func NewUUID() string {
	return gouuid.NewV4().String()
}

//IsValidUUID retorna um UUID válido
func IsValidUUID(uuid string) bool {
	_, err := gouuid.FromString(uuid)
	return err == nil
}
