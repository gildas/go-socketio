package socketio

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

func base64ID() string {
	buffer := &bytes.Buffer{}
	blob  := sha512.Sum512([]byte(fmt.Sprintf("%s%d", time.Now(), rand.Uint64())))
	encoder := base64.NewEncoder(base64.URLEncoding, buffer)
	_, _ = encoder.Write(blob[:])
	encoder.Close()
	return buffer.String()[:12]
}
