package storage

import "errors"

var (
	ErrURLNotFound = errors.New("адрес с таким псевдонимом не найден")
	ErrURLExists   = errors.New("адрес с таким псевдонимом уже существует")
)
