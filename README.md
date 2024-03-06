# Papermache

This project is aimed to simplify storage of encrypted secrets (e.g. passwords and pass phrases etc.) in paper form or in PDF format.

This is useful when such secrets are to be available, but storage locations (e.g. personal files) can be accessed by unauthorized party.

Secrets are encrypted and stored as QR codes (possibly with hex representation of cypher). Reading of QR codes is done with companion android application (see [Papermache-android](https://github.com/vaal12-go/papermache-android) sister project.

## HOW-TO

Download application executable at [Releases](https://github.com/vaal12-go/papermache/releases/tag/Release_0.1.2_03Feb2024) page.

Run the executable file, which should look like this
![Application in windows explorer](/README-images/Explorer%20application.png "App in Windows explorer")

First it should open console window like this
![Application console window](/README-images/AppConsoleWindow.png)

And in about 2-3 second in default browser it should open window with main interface
![Application in browser](/README-images/Application%20in%20browser.png)

Enter text to be encrypted below "Text to encrypt"
![Text entered](/README-images/TextEntered.png)

> Please note that max number of characters which QR image can hold is about 2000. Encryption adds some overhead (30-50%), so actual number of charactes can be about 1200. But QR images with such great number of characters produce very high detailed QR codes, which are harder to scan. For this recommendation is to stick to text length up to 600-700 latin characters. If you use non latin alphabet, then it should be about 300-500 characters for best results.

Number of characters in the text to be encoded will be displayed below text. Number of bytes should not exceed 2000.
![Number of characters](/README-images/NumberOfCharacters.png)

If you need to have a long strong password, you can press "Generate" button to the right of the password field (red arror on the image above). 
![Generated password](/README-images/GeneratedPassword.png)

We will use in this example a simple password: qwe1
![Encrypt data](/README-images/EncryptData.png)

press "Encrypt data" button below the key field (red arrow on the image above)

Result of encryption will look like this:
![Encrypted text](/README-images/Encrypted%20text.png)

You can print this page to your printer or to PDF using your browser's facilities (Ctrl-P is the usual shortcut for such printing dialog)
![PrintDialog](/README-images/PrintDialog.png)

QR code contains encrypted text, which cannot be decyphered without your key. Encryption used is [Advanced Encryption Standard](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard).

> Wiki: The Advanced Encryption Standard (AES), also known by its original name Rijndael (Dutch pronunciation: [ˈrɛindaːl]), is a specification for the encryption of electronic data established by the U.S. National Institute of Standards and Technology (NIST) in 2001.

### Exiting application

To stop application at the moment in console window press Ctrl-C and close browser window. No information will be stored on your computer.

## Technology

Papermache uses simple local http server and locally connected webbrowser. No data is ever sent to external machines on network and papermache can perfectly work on computers/phones, which are not connected to the network.

AES is used for encryption.

### Supported Operating systems

#### Generating encrypted QR code

- Windows (tested on 10, should work on 11)
- Linux support is in immediate plans
- MacOS maybe will be provided, but no immediate plans

#### Reading QR code

For [Papermache-android](https://github.com/vaal12-go/papermache-android) only Android is supported at the moment with no plans for iOS application.

Reading of QR code on notebook/desktop is not possible at the moment, but this will be a medium priority (e.g. in several weeks) goal.



## Used libraries

- [Dongle](https://github.com/golang-module/dongle) library for encryption

- [Go-qrcode](https://github.com/yeqown/go-qrcode/) for generating qr-codes

- [Godotenv](https://github.com/joho/godotenv) for testing environment during development.

- Favicon - Key by randi from [Noun Project](https://thenounproject.com/browse/icons/term/key/) (CC BY 3.0) 