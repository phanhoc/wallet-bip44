package common

import (
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcwallet/snacl"
	"io/ioutil"
	"nextop/c-horde/walletengine/offlinetool/common"
	"strings"
)

const defaultPassphrase = "yeshi dolma"

type cipherOptions struct {
	N, R, P int
}

var defaultCipherOptions = cipherOptions{
	N: 262144, // 2^18
	R: 8,
	P: 1,
}

type encryptedJSON struct {
	Data   string `json:"data"`
	Params string `json:"cipherparams"`
}

type Encrypter struct {
	filePath   string
	passPhrase *string
	option     *cipherOptions
}

// NewEncrypter creates new encryter and returns its. If the
// passPhrase is nil, encrypter will use default passphrase of system
func NewEncrypter(filePath string, passPhrase *string) *Encrypter {
	return &Encrypter{filePath: filePath, passPhrase: passPhrase, option: &defaultCipherOptions}
}

// NewEncrypterToBuffer creates new encryter using for encrypt data to buffer and returns its. If the
// passPhrase is nil, encrypter will use default passphrase of system
func NewEncrypterWriteBuffer(passPhrase *string) *Encrypter {
	return &Encrypter{passPhrase: passPhrase, option: &defaultCipherOptions}
}

// WriteEncryptedFile writes a struct to encrypted file.
func (e *Encrypter) WriteEncryptedFile(v []byte) error {
	return e.writeEncryptedToFile(v)
}

// WriteToBuffer write a struct to buffer with data encrypted
func (e *Encrypter) WriteToBuffer(v []byte) ([]byte, error) {
	return e.encyptedData(v)
}

// ReadEncrytedFile reads struct from encrypted file.
func (e *Encrypter) ReadEncrytedFile() ([]byte, error) {
	return e.getRawDataFromEncryptedFile()
}

// ReadEncrytedFromBuffer reads struct from buffer.
func (e *Encrypter) ReadEncrytedFromBuffer(data []byte) ([]byte, error) {
	return e.decryptedData(data)
}

// writeEncryptedToFile is helper to encrypt data and write to file keystore
func (e *Encrypter) writeEncryptedToFile(data []byte) error {
	encrypted, err := e.encyptedData(data)
	if err != nil {
		return err
	}

	return common.WriteDataToFile(e.filePath, encrypted)
}

// getRawDataFromEncryptedFile is helper to decrypted keystore file  and return raw data
func (e *Encrypter) getRawDataFromEncryptedFile() ([]byte, error) {
	e.filePath = strings.TrimSpace(e.filePath)
	data, err := ioutil.ReadFile(e.filePath)
	if err != nil {
		return nil, err
	}

	return e.decryptedData(data)
}

// encyptedData return []byte after encrypted data
func (e *Encrypter) encyptedData(data []byte) ([]byte, error) {
	pass := e.getPassPhrase(e.passPhrase)
	secretKey, err := snacl.NewSecretKey(
		&pass,
		e.option.N,
		e.option.R,
		e.option.P)
	if err != nil {
		return nil, err
	}

	encryptedData, err := secretKey.Encrypt(data)
	if err != nil {
		return nil, err
	}

	cipherParams := secretKey.Marshal()

	return json.Marshal(&encryptedJSON{
		Data:   hex.EncodeToString(encryptedData),
		Params: hex.EncodeToString(cipherParams)})

}

// decryptedData return []byte after decrypted data
func (e *Encrypter) decryptedData(data []byte) ([]byte, error) {
	var encrypted encryptedJSON
	if err := json.Unmarshal([]byte(data), &encrypted); err != nil {
		return nil, err
	}

	var secretKey snacl.SecretKey
	cipherParams, err := hex.DecodeString(encrypted.Params)
	if err != nil {
		return nil, err
	}

	if err = secretKey.Unmarshal(cipherParams); err != nil {
		return nil, err
	}

	pass := e.getPassPhrase(e.passPhrase)

	if err = secretKey.DeriveKey(&pass); err != nil {
		return nil, err
	}

	encryptedData, err := hex.DecodeString(encrypted.Data)
	if err != nil {
		return nil, err
	}

	return secretKey.Decrypt(encryptedData)
}

// getPassPhrase return default passphrase if passphrase input is nil, otherwise return passphrase input
func (e *Encrypter) getPassPhrase(passPhrase *string) []byte {
	if passPhrase != nil {
		return []byte(*passPhrase)
	}
	return []byte(defaultPassphrase)
}
