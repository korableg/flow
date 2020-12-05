package cmn

import (
	"math/rand"
	"regexp"
	"testing"
	"time"
)

func TestCmn(t *testing.T) {

	const contests = 10

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < contests; i++ {
		b := make([]byte, r.Intn(100))
		for k := range b {
			b[k] = byte(r.Intn(255))
		}

		testName := string(b)
		went, _ := regexp.MatchString("^[a-zA-Z0-9_.-]*$", testName)
		got := NameMatchedPattern(testName)

		if went != got {
			t.Errorf("%s, must be computed as %t", testName, went)
		}
	}

	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_.-"

	for i := 0; i < contests; i++ {
		b := make([]byte, r.Intn(100))
		for k := range b {
			b[k] = charset[r.Intn(len(charset))]
		}

		testName := string(b)
		went, _ := regexp.MatchString("^[a-zA-Z0-9_.-]*$", testName)
		got := NameMatchedPattern(testName)

		if went != got {
			t.Errorf("%s, must be computed as %t", testName, went)
		}
	}

}
