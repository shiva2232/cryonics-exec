package auth

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	firebaseCertsURL  = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"
	firebaseProjectID = "cryonics-em"
)

var certs map[string]*x509.Certificate
var certsExpiry time.Time

type FirebaseClaims struct {
	jwt.RegisteredClaims
	Firebase map[string]interface{} `json:"firebase,omitempty"`
}

func fetchCerts() (map[string]*x509.Certificate, error) {
	if certs != nil && time.Now().Before(certsExpiry) {
		return certs, nil
	}

	resp, err := http.Get(firebaseCertsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rawCerts map[string]string
	if err := json.Unmarshal(body, &rawCerts); err != nil {
		return nil, err
	}

	certs = make(map[string]*x509.Certificate)
	for kid, certPEM := range rawCerts {
		block, _ := pem.Decode([]byte(certPEM))
		if block == nil {
			return nil, fmt.Errorf("failed to decode cert PEM for kid %s", kid)
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse cert for kid %s: %v", kid, err)
		}
		certs[kid] = cert
	}
	certsExpiry = time.Now().Add(time.Hour)
	return certs, nil
}

func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func VerifyFirebaseIDToken(ctx context.Context, idToken string) (*FirebaseClaims, error) {
	certs, err := fetchCerts()
	if err != nil {
		return nil, err
	}

	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, _, err := parser.ParseUnverified(idToken, &FirebaseClaims{})
	if err != nil {
		return nil, err
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("token header missing kid")
	}

	cert, ok := certs[kid]
	if !ok {
		return nil, errors.New("no cert found for kid " + kid)
	}

	key := cert.PublicKey

	claims := &FirebaseClaims{}
	token, err = jwt.ParseWithClaims(idToken, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}, jwt.WithValidMethods([]string{"RS256"}))
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	expectedIssuer := "https://securetoken.google.com/" + firebaseProjectID
	if claims.Issuer != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer: %s", claims.Issuer)
	}

	if !containsString(claims.Audience, firebaseProjectID) {
		return nil, fmt.Errorf("invalid audience: %v", claims.Audience)
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
