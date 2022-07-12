package securities

import (
	"fmt"
	"testing"
)

func TestAES(t *testing.T) {
	key := "testtesttesttest"
	text := "text to encrypt"
	encText, err := EncryptAES(key, []byte(text))
	if err != nil {
		fmt.Println(err)
		t.Errorf("Output expect EncryptAES is success")
	}
	decBytes, err := DecryptAES(key, encText)
	if err != nil {
		t.Errorf("Output expect DecryptAES is success")
	}
	decText := string(decBytes)
	if decText != text {
		t.Errorf("Output expect %v instead of %v", text, decText)
	}
}
