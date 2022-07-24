package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/kmx0/devops/internal/types"
)

// CheckHash - проверка хеша при приеме сервером метрик
func CheckHash(metrics types.Metrics, key string) error {

	if key == "" {
		return nil
	}
	switch metrics.MType {
	case "counter":
		if metrics.Delta == nil {
			return errors.New("recieved nil pointer on Delta")
		}

		hash := Hash(fmt.Sprintf("%s:counter:%d", metrics.ID, *metrics.Delta), key)
		if hash != metrics.Hash {
			return errors.New("hash sum not matched")
		}

	case "gauge":
		if metrics.Value == nil {
			return errors.New("recieved nil pointer on Value")
		}
		hash := Hash(fmt.Sprintf("%s:gauge:%f", metrics.ID, *metrics.Value), key)
		if hash != metrics.Hash {
			return errors.New("hash sum not matched")
		}
	}
	return nil

}

// Hash - подписывает алгоритмом HMAC, используя SHA256
func Hash(src, key string) (hash string) {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(src))
	dst := h.Sum(nil)
	hash = fmt.Sprintf("%x", dst)
	return
}
