package middleware

import (
	"context"
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"order/app/dto"
	"strings"
)

const (
	AuthorizationKey  = "Authorization"
	RS256             = "RS256"
	AcuityAuthContext = "acuity-claims"
)

type JwtAuth interface {
	generateJwtToken(userId string) (string, error)
	parseClaims(token string, claims jwt.Claims) error
}

type jwtParse struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	parser     *jwt.Parser
}

var parse JwtAuth

// InitJwtParse init jwt parse with public-key and private-key
func InitJwtParse(keyPaths ...string) error {
	var (
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
	)
	if len(keyPaths) > 0 && keyPaths[0] != "" {
		k, err := readKeyFile(keyPaths[0])
		if err != nil {
			return err
		}
		pk, err := jwt.ParseRSAPublicKeyFromPEM(k)
		if err != nil {
			return err
		}
		publicKey = pk
	}
	if len(keyPaths) > 1 && keyPaths[1] != "" {
		k, err := readKeyFile(keyPaths[1])
		if err != nil {
			return err
		}
		pk, err := jwt.ParseRSAPrivateKeyFromPEM(k)
		if err != nil {
			return err
		}
		privateKey = pk
	}

	parse = &jwtParse{
		publicKey:  publicKey,
		privateKey: privateKey,
		parser: &jwt.Parser{
			ValidMethods:         []string{RS256},
			SkipClaimsValidation: false,
		},
	}
	return nil
}

// parseClaims parsing jwt token
func (jwp *jwtParse) parseClaims(tokenStr string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwp.publicKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}

func (jwp *jwtParse) generateJwtToken(userId string) (string, error) {
	claims := dto.AuthClaims{
		UserId: userId,
	}
	tkn := jwt.NewWithClaims(jwt.GetSigningMethod(RS256), claims)
	return tkn.SignedString(jwp.privateKey)
}

func readKeyFile(filePath string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// GenerateToken using for testing local, it must not be cherry-picked to master branch for deploying on production
func GenerateToken(userId string) (string, error) {
	return parse.generateJwtToken(userId)
}

// Authorize using for application authenticate
func Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			jwtToken = ctx.GetHeader(AuthorizationKey)
			splits   = strings.Split(jwtToken, " ")
		)
		var token string
		if len(splits) > 2 {
			ctx.JSON(http.StatusBadRequest, errors.New("token is invalid"))
			ctx.Abort()
			return
		}
		if len(splits) == 2 {
			if splits[0] != "Bearer" {
				ctx.JSON(http.StatusBadRequest, errors.New("token is invalid"))
				ctx.Abort()
				return
			}
			token = splits[1]
		}
		var authClaims dto.AuthClaims
		err := parse.parseClaims(token, &authClaims)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			ctx.Abort()
			return
		}
		ctxWithClaim := context.WithValue(ctx, AcuityAuthContext, &authClaims)
		ctx.Request = ctx.Request.WithContext(ctxWithClaim)
	}
}
