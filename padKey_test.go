package main

import "testing"

var padKeyTestSet = map[string]string{
	"qwe1": "qwe1============",
}

func TestPadKey(t *testing.T) {
	for key, val := range padKeyTestSet {
		paddedKey := padKey(key)
		if paddedKey != val {
			t.Errorf("Value returned was not as expected. Key supplied:%s\n value returned:%s\n value expected:%s",
				key, paddedKey, val)
		}
		if (len(key) <= 16) && (len(paddedKey)) != 16 {
			t.Errorf("Len of the padded key is not 16. Key supplied:%s\n value returned:%s\n value expected:%s\n Len of the padded key:%d",
				key, paddedKey, val, len(paddedKey))
		}

	}

}
