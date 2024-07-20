package cache

import (
	"strconv"
	"time"
)

type MockCache struct {
	Data map[string]string
}

func (m *MockCache) Get(key string) (string, error) {
	val, ok := m.Data[key]
	if !ok {
		return "", nil
	}
	return val, nil
}

func (m *MockCache) Increment(key string) error {
	val, ok := m.Data[key]
	if !ok {
		m.Data[key] = "1"
	} else {
		num, _ := strconv.Atoi(val)
		m.Data[key] = strconv.Itoa(num + 1)
	}
	return nil
}

func (m *MockCache) Decrement(key string) (int64, error) {
	val, ok := m.Data[key]
	if !ok {
		return 0, nil
	}
	num, _ := strconv.Atoi(val)
	if num == 1 {
		delete(m.Data, key)
	} else {
		m.Data[key] = strconv.Itoa(num - 1)
	}
	return int64(num - 1), nil
}

func (m *MockCache) Delete(key string) error {
	delete(m.Data, key)
	return nil
}

func (m *MockCache) ControlExpirationTime(key string) {}

func (m *MockCache) Set(key string, value string, expiration time.Duration) error {
	m.Data[key] = value
	return nil
}
