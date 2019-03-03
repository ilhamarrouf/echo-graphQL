package middlewares

import (
	"fmt"
	"time"
	"github.com/ilhamarrouf/echo-graphql/libs"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type MyClaim struct {
	UserId int64
	IsAdmin bool
	RefreshJti string
	jwt.StandardClaims
}

func createRefreshTokenString(userid int64) (refreshTokenString string, err error) {
	refreshJti, err := libs.StoreRefreshToken()
	if err != nil {
		return "", err
	}

	if userid != 0 {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &MyClaim{
			UserId: userid,
			IsAdmin: false,
			RefreshJti: refreshJti,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			}})

		t, err := token.SignedString([]byte("secret2"))
		if err != nil {
			return "", err
		}

		return t, err
	}

	return "", echo.ErrUnauthorized
}

func CreateAuthTokenString(userid int64) (authTokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &MyClaim{
		UserId: userid,
		IsAdmin: true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		}})

	t, err := token.SignedString([]byte("secret1"))
	if err != nil{
		return "", nil
	}

	return t, err
}

func CreateNewTokens(username string, password string) (authTokenString string, refreshTokenString string, err error) {
	user := libs.FetchUser(username, password)
	fmt.Println(user)

	refreshTokenString, err = createRefreshTokenString(user.Id)
	if err != nil {
		return "", "", nil
	}

	authTokenString, err = CreateAuthTokenString(user.Id)
	if err != nil {
		return "", "", nil
	}

	return
}

func UpdateRefreshTokenExp(myClaim *MyClaim, oldTokenString string) (newTokenString, newRefreshTokenString string, err error) {
	myClaim2 := MyClaim{}

	_, err = jwt.ParseWithClaims(oldTokenString, &myClaim2, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret1"), nil
	})
	if err != nil {
		return "", "", err
	}

	if !libs.CheckRefreshToken(myClaim.RefreshJti) || myClaim.UserId != myClaim2.UserId {
		return "", "", fmt.Errorf("error: %s", "old token is invalid")
	}

	libs.DeleteRefreshToken(myClaim2.RefreshJti)

	newRefreshTokenString, err = createRefreshTokenString(myClaim2.UserId)
	if err != nil {
		return "", "", nil
	}

	newTokenString, err = CreateAuthTokenString(myClaim2.UserId)
	if err != nil {
		return "", "", nil
	}

	return
}