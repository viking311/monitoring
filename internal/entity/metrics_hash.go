package entity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func MetricsHash(data Metrics, hashKey string) string {
	hash := ""
	if len(hashKey) > 0 {
		src := ""
		hasher := hmac.New(sha256.New, []byte(hashKey))
		if data.MType == "counter" {
			src = fmt.Sprintf("%s:counter:%d", data.ID, *data.Delta)
		}

		if data.MType == "gauge" {
			src = fmt.Sprintf("%s:gauge:%f", data.ID, *data.Value)
		}
		if len(src) > 0 {
			hasher.Write([]byte(src))
			hash = hex.EncodeToString(hasher.Sum(nil))
		}
	}
	return hash
}
