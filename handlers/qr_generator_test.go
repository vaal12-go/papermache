package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/golang-module/dongle"
	"github.com/makiuchi-d/gozxing"
	zxingqrcode "github.com/makiuchi-d/gozxing/qrcode"
	"github.com/vaal12-go/papermache/encrypt"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"golang.org/x/image/draw"
)

//TO run: go test -run TestQRGenerationRand -v

const STRING_LEN = 1600
const KEY_LEN = 10
const NUMBER_OF_QR_GENERATIONS = 20

const TEST_KEY = "qwe1"
const TIME_FOR_ANIMATION = 0

func generateTestPlainString(strLen int) string {
	testStringBuilder := strings.Builder{}
	for testStringBuilder.Len() < strLen {
		currAlphabetIdx := rand.Intn(len(alphabets))
		currCharIdx := rand.Intn(len([]rune(alphabets[currAlphabetIdx])))
		testStringBuilder.WriteRune([]rune(alphabets[currAlphabetIdx])[currCharIdx])
	}
	// fmt.Printf("test_string.String(): %v\n", testStringBuilder.String())
	// fmt.Printf("\tlen: %v\n", len(testStringBuilder.String()))

	return testStringBuilder.String()
}

func generateTestKey(keyLen int) []byte {
	testStringBuilder := strings.Builder{}
	for len([]rune(testStringBuilder.String())) < keyLen {
		currAlphabetIdx := rand.Intn(2) //This will only include ENG and punctuations
		//TODO: to extend to all alphabets
		currCharIdx := rand.Intn(len([]rune(alphabets[currAlphabetIdx])))
		testStringBuilder.WriteRune([]rune(alphabets[currAlphabetIdx])[currCharIdx])
	}
	// fmt.Printf("test_string.String(): %v\n", testStringBuilder.String())
	// fmt.Printf("\tlen: %v\n", len(testStringBuilder.String()))
	// fmt.Printf("\tgeneratedKey: %v\n", testStringBuilder.String())
	return []byte(testStringBuilder.String())
} //func generateTestKey(keyLen int) []byte {

// func yencodeByteArray(arrayToEncode []byte) ([]byte, error) {
// 	srcReader := bytes.NewReader(arrayToEncode)
// 	encodedBytes := bytes.Buffer{}
// 	yencWriter := NewYEncoder(&encodedBytes)

// 	_, err := io.Copy(yencWriter, srcReader)
// 	if err != nil && !(errors.Is(err, io.ErrUnexpectedEOF)) {
// 		return nil, err
// 	}
// 	return encodedBytes.Bytes(), nil
// }

// func ydecodeByteArray(arrayToDecode []byte) ([]byte, error) {
// 	decodedBuf := &bytes.Buffer{}
// 	decReader := NewYDecoder(
// 		bytes.NewReader(arrayToDecode))
// 	_, err := io.Copy(decodedBuf, decReader)
// 	if err != nil && !(errors.Is(err, io.ErrUnexpectedEOF)) {
// 		return nil, err
// 	}

// 	return decodedBuf.Bytes(), nil
// }

func saveQRCode(code *qrcode.QRCode, fName string) error {
	writ, err := standard.New(fName)
	if err != nil {
		return err
	}
	if err = code.Save(writ); err != nil {
		return err
	}
	return nil
}

func TestQRGenerationRand(t *testing.T) {
	for i := 0; i < NUMBER_OF_QR_GENERATIONS; i++ {
		testString := generateTestPlainString(STRING_LEN)
		// testString := "qwe1"
		testKey := generateTestKey(KEY_LEN)
		// fmt.Printf("testString: %v\n", testString)

		//TEST ENCRYPTION
		encrText, err := encrypt.Encrypt([]byte(testString),
			[]byte(testKey))
		if err != nil {
			t.Fatalf("Error encrypting:%v\n", err)
		}

		decrText, err := encrypt.Decrypt(encrText, testKey)
		if err != nil {
			t.Fatalf("Error decrypting:%v\n", err)
		}
		// fmt.Printf("string(decrText): %v\n", string(decrText))
		if bytes.Compare([]byte(testString), decrText) != 0 {
			t.Logf("testString: %v\n", testString)
			t.Logf("testKey: %v\n", testKey)
			t.Logf("string(decrText): %v\n", string(decrText))
			t.Fatalf("Encrypted and decrypted strings are different.")
		}

		// time.Sleep(time.Millisecond * TIME_FOR_ANIMATION)

		encryptedBase64String := dongle.Encode.FromBytes(encrText).ByBase64().ToString()
		qrCode, err := QRCodeFromString(encryptedBase64String)
		if err != nil {
			t.Fatalf("Error creating QRCode. Err:%v\n", err)
			return
		}

		// err = saveQRCode(qrCode,
		// 	fmt.Sprintf("test_qrcode_%03d.jpeg", i))
		// if err != nil {
		// 	t.Fatalf("Error saving QRCode as file. Err:%v\n", err)
		// 	return
		// }
		imgBuf, err := saveQRCodeToByteArray(qrCode)
		if err != nil {
			t.Fatalf("Error saving QRCode as stream. Err:%v\n", err)
			return
		}

		// fmt.Printf("\"\n\n QRCODE GENERATED \n\n\": %v\n", "\n\n QRCODE GENERATED \n\n")
		// fmt.Printf("len(imgBuf): %v\n", len(imgBuf))

		//TODO: move below to separate function (reading of QR code from byte array)
		bReader := bytes.NewReader(imgBuf)
		img, _, err := image.Decode(bReader)

		// Encode to `output`:
		// png.Encode(output, dst)

		if err != nil {
			t.Fatalf("Error decoding image in request err: %v\n", err)
			return
		}

		imgSizesArray := []int{
			1600, 1400, 1000, 800, 600, 500,
		}

		var result *gozxing.Result
		initialImage := img

		for k := 0; k < len(imgSizesArray); k++ {
			// fmt.Printf("img: %v\n", img.Bounds())

			bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
			qrReader := zxingqrcode.NewQRCodeReader()
			result, err = qrReader.Decode(bmp, map[gozxing.DecodeHintType]interface{}{
				gozxing.DecodeHintType_TRY_HARDER:   true,
				gozxing.DecodeHintType_PURE_BARCODE: true,
			})

			if err != nil {
				t.Logf("i: %v\n", i)
				t.Logf("qrReader.Decode err: %v  Will retry\n", err)

			} else {
				break
			}
			t.Logf("will rescale image to:%v\n", imgSizesArray[k])
			img, err = rescaleQRCode(initialImage, imgSizesArray[k], imgSizesArray[k])
			if err != nil {
				// t.Logf("i: %v\n", i)
				t.Logf("Error resizing image: %v\n", err)
				return
			}
		}

		if result == nil {
			t.Fatalf("Rescaling attempts failed. Aborting")
			return
		}

		// if err != nil {
		// 	t.Logf("i: %v\n", i)
		// 	t.Logf("qrReader.Decode err: %v  Will retry\n", err)

		// 	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
		// 	qrReader := zxingqrcode.NewQRCodeReader()
		// 	result, err = qrReader.Decode(bmp, nil)

		// 	if err != nil {
		// 		t.Logf("i: %v\n", i)
		// 		t.Fatalf("qrReader.Decode err: %v\n", err)
		// 		return
		// 	}
		// 	t.Logf("2nd attempt OK\n")
		// 	return
		// }
		// fmt.Printf("QR code read (bytes):%v\n", result.GetRawBytes())
		// fmt.Printf("QR code read (bytes) len:%v\n", len(result.GetRawBytes()))

		// txt := result.GetText()
		// fmt.Printf("result.GetText():%v\n", txt)
		// fmt.Printf("result.GetText():%v\n", len(txt))
		// fmt.Printf("QR code read(str):%v\n", result.String())
		// fmt.Printf("QR code read len (str):%v\n", len(result.String()))
		// fmt.Printf("QR code read len (byteArr from Str):%v\n", ([]byte(result.String())))
		// fmt.Printf("QR code read len (byteArr from Str):%v\n", len([]byte(result.String())))

		// fmt.Printf("QR code read len (runeArr from Str):%v\n", []rune(result.String()))
		// fmt.Printf("QR code read len (runeArr from Str):%v\n", len([]rune(result.String())))

		// decodedArr, err = ydecodeByteArray([]byte(result.String()))
		// fmt.Printf("\n\n\ndecodedArr: %v\n", decodedArr)
		// fmt.Printf("decodedArr: %v\n", string(decodedArr))
		// if err != nil {
		// 	t.Fatalf("Error decoding byte array (from image). Err:%v\n", err)
		// 	return
		// }
		decodedArr := dongle.Decode.FromString(result.String()).ByBase64().ToBytes()
		decrText, err = encrypt.Decrypt(decodedArr, testKey)

		if err != nil {
			t.Fatalf("Error decrypting (from image):%v\n", err)
		}
		// fmt.Printf("decrText: %s\n", string(decrText))

		if bytes.Compare([]byte(testString), decrText) != 0 {
			t.Logf("Compare (from image).testString: %v\n", testString)
			t.Logf("testKey: %v\n", testKey)
			t.Logf("string(decrText): %v\n", string(decrText))
			t.Fatalf("Encrypted and decrypted strings are different.")
		}

	}

} //func TestQRGenerationRand(t *testing.T) {

func rescaleQRCode(img image.Image, newWidth int, newHeight int) (image.Image, error) {
	//From here: https://stackoverflow.com/questions/31463756/convert-image-image-to-image-nrgba
	//A bit also from here: https://stackoverflow.com/questions/22940724/go-resizing-images

	// Set the expected size that you want:
	// dst := image.NewRGBA(image.Rect(0, 0,
	// 	img.Bounds().Max.X/2, img.Bounds().Max.Y/2))
	dst := image.NewRGBA(image.Rect(0, 0,
		800, 800))

	// Resize:
	draw.NearestNeighbor.Scale(dst, dst.Rect, img,
		img.Bounds(), draw.Over, nil)

	resizedBuf := bytes.Buffer{}
	err := png.Encode(&resizedBuf, dst)
	if err != nil {
		return nil, err
	}
	bReader := bytes.NewReader(resizedBuf.Bytes())
	retImage, _, err := image.Decode(bReader)
	if err != nil {
		return nil, err
	}

	output, _ := os.Create(
		fmt.Sprintf("your_image_resized_%d x %d.png", newWidth, newHeight))
	defer output.Close()

	bReader = bytes.NewReader(resizedBuf.Bytes())
	io.Copy(output, bReader)

	return retImage, nil
}
