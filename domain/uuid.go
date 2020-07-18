package domain

import (
	gouuid "github.com/google/uuid"
)

func NewUUID() string {
	return gouuid.New().String()
}

//IsValidUUID retorna um UUID válido
func IsValidUUID(uuid string) bool {
	_, err := gouuid.Parse(uuid)
	return err == nil
}
