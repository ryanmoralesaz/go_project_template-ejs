package json

import ()

type UserStorage struct {
	filePath string
	mutex sync.RWMutex
}