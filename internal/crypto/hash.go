package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/kmx0/devops/internal/types"
)

func CheckHash(metrics types.Metrics, key string) error {
	switch metrics.MType {
	case "counter":
		if metrics.Delta == nil {
			return errors.New("recieved nil pointer on Delta")
		}
		if metrics.Hash != "" {
			hash := Hash(fmt.Sprintf("%s:counter:%d", metrics.ID, *metrics.Delta), key)
			if hash != metrics.Hash {
				return errors.New("hash sum not matched")
			}
		}

	case "gauge":
		if metrics.Value == nil {
			return errors.New("recieved nil pointer on Value")
		}
		if metrics.Hash != "" {
			hash := Hash(fmt.Sprintf("%s:counter:%f", metrics.ID, *metrics.Value), key)
			if hash != metrics.Hash {
				return errors.New("hash sum not matched")
			}
		}
	}
	return nil

}

func Hash(src, key string) (hash string) {
	// hash(fmt.Sprintf("%s:counter:%d", id, delta), key),
	// подписываем алгоритмом HMAC, используя SHA256
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(src))
	hash = string(h.Sum(nil))
	return
}
