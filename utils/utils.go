package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"umn-technology/config"
	"umn-technology/constants"
	"umn-technology/models"

	"github.com/dgrijalva/jwt-go"

	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type TxFn func(Transaction) error

func WithTransaction(db *sql.DB, fn TxFn) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

func BindValidateStruct(ctx echo.Context, i interface{}) error {
	if err := ctx.Bind(i); err != nil {
		return err
	}

	if err := ctx.Validate(i); err != nil {
		return err
	}

	return nil
}

func ResponseJSON(success bool, code string, msg string, result interface{}) models.Response {
	tm := time.Now()
	response := models.Response{
		Success:          success,
		StatusCode:       code,
		Result:           result,
		Message:          msg,
		ResponseDatetime: tm,
	}

	return response
}
func ResponseListJSON(success bool, code string, msg string, countData int, result interface{}) models.ResponseList {
	tm := time.Now()
	response := models.ResponseList{
		Success:          success,
		StatusCode:       code,
		CountData:        countData,
		Result:           result,
		Message:          msg,
		ResponseDatetime: tm,
	}

	return response
}

func TimeStampNow() string {
	return time.Now().Format(constants.LAYOUT_TIMESTAMP)
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func DBTransaction(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Rollback Panic
		} else if err != nil {
			tx.Rollback() // err is not nill
		} else {
			err = tx.Commit() // err is nil
		}
	}()
	err = txFunc(tx)
	return err
}

func Stringify(input interface{}) string {
	bytes, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	strings := string(bytes)
	bytes, err = json.Marshal(strings)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func JSONPrettyfy(data interface{}) string {
	bytesData, _ := json.MarshalIndent(data, "", "  ")
	return string(bytesData)
}

func ToString(i interface{}) string {
	log, _ := json.Marshal(i)
	logString := string(log)

	return logString
}

func GenerateToken(request models.Login) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = request.Username
	claims["password"] = request.Password
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	// Generate encoded token and send it as response.
	restoken, err := token.SignedString([]byte(config.GetEnv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return restoken, nil
}

func GenerateRefreshToken(request models.Login) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// Set claims for refresh token
	claims["username"] = request.Username
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix() // Refresh token expires in 7 days

	refreshToken, err := token.SignedString([]byte(config.GetEnv("JWT_REFRESH_KEY")))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

// Validate Token
// ValidateToken validates the given JWT token string
func ValidateToken(signedToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetEnv("JWT_KEY")), nil
	})
	if err != nil {
		// Log the error for debugging
		fmt.Println("Error parsing token:", err)
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	// Extract and return claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error parsing claims")
	}

	return claims, nil
}

func Encrypt(key, text string) (string, error) {
	// Ensure the key is 32 bytes long (256-bit AES)
	if len(key) != 32 {
		return "", fmt.Errorf("invalid key size: key must be 32 bytes long")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plaintext := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(key, cryptoText string) (string, error) {
	// Ensure the key is 32 bytes long (256-bit AES)
	if len(key) != 32 {
		return "", fmt.Errorf("invalid key size: key must be 32 bytes long")
	}

	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func Hash(text string) string {
	h := sha1.New()
	h.Write([]byte(text))

	bs := h.Sum(nil)

	res := fmt.Sprintf("%x\n", bs)

	return res
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RefreshAccessToken(refreshToken string) (string, string, error) {
	// Parse the refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetEnv("JWT_REFRESH_KEY")), nil
	})
	if err != nil {
		return "", "", err
	}

	// Verify token validity
	if !token.Valid {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("error parsing claims")
	}

	// Generate new access token
	accessTokenClaims := jwt.MapClaims{
		"username": claims["username"],
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Access token expires in 10 minutes
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(config.GetEnv("JWT_ACCESS_KEY")))
	if err != nil {
		return "", "", err
	}

	// Generate new refresh token
	newRefreshTokenClaims := jwt.MapClaims{
		"username": claims["username"],
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Refresh token expires in 24 hours
	}
	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newRefreshTokenClaims)
	newRefreshTokenString, err := newRefreshToken.SignedString([]byte(config.GetEnv("JWT_REFRESH_KEY")))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, newRefreshTokenString, nil
}
