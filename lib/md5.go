package lib

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

func RamdomTokenGenerator() string {
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(1e11)
	data := []byte(strconv.Itoa(x))
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
