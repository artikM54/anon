package hash

import (
	"anonymous_chat/internal/redis"
	"context"
	"crypto/rand"
	"encoding/hex"
)

func CreateUniqueModelHash(nameList string) string {
	ctx := context.Background()
	hash := generateRandomHash()

	added, err := redis.Client.SAdd(ctx, nameList, hash).Result()
	if err != nil {
		panic(err)
	}

	if added == 0 {
		hash = CreateUniqueModelHash(nameList)
	}

	return hash
}

func generateRandomHash() string {
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		println("Error 123 hash util")
	}

	randomHash := hex.EncodeToString(randomBytes)

	return randomHash
}
