package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const secret = "sysu_activity_2018_activity_sysu"

// GetPoster judge if the poster and returns accurate one with given type
func GetPoster(raw string, actType int) string {
	if len(raw) == 0 {
		switch actType {
		// physics
		case 0:
			return "b6f487c6d08921463a6ebc0612d9fe1f.gif"
		// volunteer
		case 1:
			return "ccc55f553829fabb7c15227d79450dae.gif"
		// match
		case 2:
			return "2bee829b10b0a84002cf5cb5c4a3c8f3.gif"
		// show
		case 3:
			return "68dac067d05a98995a353ad8265b1f09.png"
		// speech
		case 4:
			return "a90dc26fbd5299e4053a3bbc39b5afc8.gif"
		// outdoor
		case 5:
			return "e8ae3078dfa14c62ff1e71104ec0b11f.png"
		// relax
		case 6:
			return "b2b71f5f39d3a4389d34ce1b248e9fee.png"
		}
	}
	return raw
}

// CheckToken Check token and return token status code with openId
// status code: 0 -> check error; 1 -> timeout; 2 -> ok
func CheckToken(tokenString string) (int, string) {
	var hmacSampleSecret = []byte(secret)
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing my secret
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expTime := claims["exp"]
		openId := claims["sub"]

		if (int64)(expTime.(float64)) <= time.Now().Unix() {
			return 1, openId.(string)
		}
		return 2, openId.(string)
	} else {
		return 0, ""
	}
}

// GetUserOpenId Get user openId from Wechat server
func GetUserOpenId(code string) (string, error) {
	var retData map[string]interface{}
	// Get from Wechat api
	resp, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=wxe8db5a32e4ca30e9&secret=bf785281b28fc2fba45b7613965bcbb1&js_code=" + code + "&grant_type=authorization_code")
	if err != nil {
		return "", err
	}

	// Read response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// decode json string
	if err = json.Unmarshal(body, &retData); err != nil {
		return "", err
	}

	// Wechat reponse uses errcode in json and status code will always be 200
	if _, ok := retData["errcode"]; ok {
		return "", errors.New(retData["errmsg"].(string))
	}

	openId := retData["openid"].(string)
	return openId, nil
}

// GenerateJWT Generate jwt with openid(sub), issuance time(iat) and expiration time(exp)
func GenerateJWT(openId string) (string, error) {
	// expire in two weeks
	var exp = time.Hour * 24 * 14
	var hmacSampleSecret = []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": openId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(exp).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	return tokenString, err
}
