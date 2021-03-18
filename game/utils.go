package game

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strconv"

	zeroMessage "github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	seed := make([]byte, 8)
	_, _ = cryptoRand.Read(seed)
	rand.Seed(int64(binary.LittleEndian.Uint64(seed)))
}

func zeroText(message string) zeroMessage.Message {
	return zeroMessage.Message{zeroMessage.Text(message)}
}

func findAt(message zeroMessage.Message) (int64, bool) {
	for _, segment := range message {
		if segment.Type == "at" {
			atStr := segment.Data["qq"]
			if atStr == "all" {
				return 0, false
			}
			qq, err := strconv.ParseInt(atStr, 10, 64)
			if err != nil {
				return 0, false
			}
			return qq, true
		}
	}
	return 0, false
}

func randInt32(min, max int32) int32 {
	v := max - min
	return min + rand.Int31n(v)
}

func randInt64(min, max int64) int64 {
	v := max - min
	return min + rand.Int63n(v)
}

func AbsSignInt32(i int32) (int32, string) {
	if i < 0 {
		return -i, "-"
	}
	return i, ""
}

func ThousandthStr(i int32) string {
	abs, sign := AbsSignInt32(i)
	t := abs / 1000
	r := abs % 1000
	return fmt.Sprintf("%s%d.%03d", sign, t, r)
}

func TenthStr(i int32) string {
	abs, sign := AbsSignInt32(i)
	t := abs / 10
	r := abs % 10
	return fmt.Sprintf("%s%d.%d", sign, t, r)
}

func zeroUserAvatar(userID int64) zeroMessage.Message {
	return zeroMessage.Message{zeroMessage.Image(fmt.Sprintf("https://q1.qlogo.cn/g?b=qq&nk=%d&s=140", userID))}
}
