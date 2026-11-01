package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

/*
https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go

parameters from the article seem reasonable
*/

// check argon2.IDKey for more information
type hashParameters struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var generationHashParameters *hashParameters = &hashParameters{
	memory:      64 * 1024,
	iterations:  3,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

// hash password using specified parameters and the salt
func (parameters *hashParameters) generate(password string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(password),
		salt,
		parameters.iterations,
		parameters.memory,
		parameters.parallelism,
		parameters.keyLength,
	)
}

// encode the hash using the standard encoded hash representation
func (parameters *hashParameters) encode(hash []byte, salt []byte) string {
	// base64 encode the salt and hashed password
	base64Salt := base64.RawStdEncoding.EncodeToString(salt)
	base64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		parameters.memory,
		parameters.iterations,
		parameters.parallelism,
		base64Salt,
		base64Hash,
	)
}

// generate the hash, then encode it
func (parameters *hashParameters) generateEncode(
	password string, salt []byte,
) string {
	return parameters.encode(parameters.generate(password, salt), salt)
}

func EncodedHashGenerate(password string) string {
	// generate a cryptographically secure random salt
	salt := make([]byte, generationHashParameters.saltLength)
	// fill the array with random bytes
	if _, err := rand.Read(salt); err != nil {
		// should not happen
		panic(err)
	}
	return generationHashParameters.generateEncode(password, salt)
}

func EncodedHashCompare(password, encodedHash string) (
	match bool, err error,
) {
	// extract the parameters, salt and derived key from the encoded hash
	parameters, salt, hash, err := decodeEncodedHash(encodedHash)
	if err != nil {
		return false, err
	}
	// hash provided string using the same parameters
	otherHash := parameters.generate(password, salt)

	// check whether the contents of the hashed passwords are identical
	// we are using subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks
	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

func formatError(err error) error {
	return fmt.Errorf("the encoded hash is not in the correct format: %s", err)
}

func decodeEncodedHash(encodedHash string) (
	parameters *hashParameters, salt, hash []byte, err error,
) {
	values := strings.Split(encodedHash, "$")
	if len(values) != 6 {
		return nil, nil, nil, formatError(errors.New("invalid number of '$'"))
	}

	var version int
	if _, err = fmt.Sscanf(values[2], "v=%d", &version); err != nil {
		return nil, nil, nil, formatError(err)
	}
	if version != argon2.Version {
		return nil, nil, nil, fmt.Errorf(
			"incompatible argon2id version: current - %d, desired - %d ",
			version,
			argon2.Version,
		)
	}

	parameters = &hashParameters{}
	if _, err = fmt.Sscanf(
		values[3],
		"m=%d,t=%d,p=%d",
		&parameters.memory,
		&parameters.iterations,
		&parameters.parallelism,
	); err != nil {
		return nil, nil, nil, formatError(err)
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return nil, nil, nil, formatError(err)
	}
	parameters.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return nil, nil, nil, formatError(err)
	}
	parameters.keyLength = uint32(len(hash))

	return parameters, salt, hash, nil
}
