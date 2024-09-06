package encrypt

import (
	"testing"

	"github.com/golang-module/dongle"
)

//TO run specific test: go test -run TestStretchKey -v
//To run from root dir: go test  ./encrypt/. -v

var stretchKeyTestSet = map[string]string{
	"qwe123":       "2f7bf3fbfe52de9898f2bc04d042805318bd9f53d664e5b0c289550b0264ac78",
	"key#23":       "ea09e743d2b3220eb50b04584eb3b9a73da4de8aa32aef3d49245a5ac20b8e5a",
	"نشرة الاخبار": "594db146b038f2c6159564fb61bf6496115b80c921496555b371cb563a1f6195",
	"澳大利亚对待残障人士的移民政": "dff1522154996dbb2665dc6c833585070752de0d8ab2c61db32d8b337ed35bde",
	"トランプ氏の機密文":      "f85f01abad69a89f80881f9fefa1148339890ccb30896ad97c20561749552e81",
}

// TODO: add test function for stretchKey

func TestStretchKey(t *testing.T) {
	for key, val := range stretchKeyTestSet {
		res := stretchKey([]byte(key))
		// fmt.Printf("key: %v\n", key)
		// fmt.Printf("res: %v\n", res)
		res_hex := dongle.Encode.FromBytes(res).ByHex().ToString()
		if val != res_hex {
			// fmt.Printf("res_hex: %v\n", res_hex)
			t.Errorf("Key:%s required hex \n(%s) \nis not matching returned hex:\n%s",
				key, val, res_hex)
		}
	}
}
