package di

import (
	"sync"
)

var (
	globalContainer *Container
	once            sync.Once
)

// GetContainer возвращает глобальный DI контейнер (singleton)
func GetContainer() *Container {
	once.Do(func() {
		globalContainer = NewContainer()
	})
	return globalContainer
}

// CloseGlobalContainer закрывает глобальный контейнер
func CloseGlobalContainer() error {
	if globalContainer != nil {
		return globalContainer.Close()
	}
	return nil
}
