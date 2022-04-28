package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var key = "testtesttesttest"

//var srcDir = filepath.FromSlash("C:\\Users\\laugh\\go\\src\\github.com\\cjg3484\\Botnet\\src")
//var srcDir = filepath.FromSlash("H:\\GitHub\\Botnet\\src")

func decrypt(cipherstring string, keystring string) string {
	// Byte array of the string
	ciphertext := []byte(cipherstring)

	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		panic("Text is too short")
	}

	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]

	// Remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext)
}

func encrypt(plainstring, keystring string) string {
	// Byte array of the string
	plaintext := []byte(plainstring)

	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Empty array of 16 + plaintext length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return string(ciphertext)
}

func writeToFile(data, file string) {
	err := ioutil.WriteFile(file, []byte(data), 0666)
	if err != nil {
		panic(err)
	}
}

func readFromFile(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	return data, err
}

func handleDecrypt(decryptCmd *flag.FlagSet) {
	//fmt.Println("Printing decrypted ciphertext.txt: ")
	if ciphertext, err := readFromFile("./ciphertext.txt"); err != nil {
		fmt.Println("File is not found")
	} else {
		plaintext := decrypt(string(ciphertext), key)
		writeToFile(plaintext, "./plaintext.txt")
		err := os.Remove("./ciphertext.txt")
		if err != nil {
			panic(err)
		}
		fmt.Println("y0ur f1l3s h4v3 b33n r37urn3d")
	}
}

func handleEncrypt(encryptCmd *flag.FlagSet) {

	writeToFile("Hello World\n¯\\_(ツ)_/¯", "./plaintext.txt")

	if plaintext, err := readFromFile("./plaintext.txt"); err != nil {
		fmt.Println("Plaintext file is not found")
	} else {
		ciphertext := encrypt(string(plaintext), key)
		encryptedFile := "./ciphertext.txt"
		writeToFile(ciphertext, encryptedFile)
		err := os.Remove("./plaintext.txt")
		if err != nil {
			panic(err)
		}
		fmt.Println("all your plaintext are belong to us")
	}
}

func main() {

	encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)

	decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'encrypt' or 'decrypt' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "encrypt":
		handleEncrypt(encryptCmd)
	case "decrypt":
		handleDecrypt(decryptCmd)
	case "default":
		fmt.Println("Command not supported")
		os.Exit(1)
	}

}
