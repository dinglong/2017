package token

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"anthserver/models"

	"github.com/docker/distribution/registry/auth/token"
	"github.com/docker/libtrust"
)

const (
	// privateKey = "/etc/ui/private_key.pem"
	privateKey = "e:/data/private_key.pem"
	issuer     = "registry-token-issuer"
)

func ParseScopes(u *url.URL) []string {
	var sector string
	var result []string
	for _, sector = range u.Query()["scope"] {
		result = append(result, strings.Split(sector, " ")...)
	}
	return result
}

// GetResourceActions
func GetResourceActions(scopes []string) []*token.ResourceActions {
	log.Printf("scopes: %+v", scopes)
	var res []*token.ResourceActions
	for _, s := range scopes {
		if s == "" {
			continue
		}
		items := strings.Split(s, ":")
		length := len(items)

		typee := ""
		name := ""
		actions := []string{}

		if length == 1 {
			typee = items[0]
		} else if length == 2 {
			typee = items[0]
			name = items[1]
		} else {
			typee = items[0]
			name = strings.Join(items[1:length-1], ":")
			if len(items[length-1]) > 0 {
				actions = strings.Split(items[length-1], ",")
			}
		}

		res = append(res, &token.ResourceActions{
			Type:    typee,
			Name:    name,
			Actions: actions,
		})
	}
	return res
}

// MakeToken makes a valid jwt token based on parms.
func MakeToken(username, service string, access []*token.ResourceActions) (*models.Token, error) {
	pk, err := libtrust.LoadKeyFile(privateKey)
	if err != nil {
		return nil, err
	}

	tk, expiresIn, issuedAt, err := makeTokenCore(issuer, username, service, 30, access, pk)
	if err != nil {
		return nil, err
	}
	rs := fmt.Sprintf("%s.%s", tk.Raw, base64UrlEncode(tk.Signature))
	return &models.Token{
		Token:     rs,
		ExpiresIn: expiresIn,
		IssuedAt:  issuedAt.Format(time.RFC3339),
	}, nil
}

//make token core
func makeTokenCore(issuer, subject, audience string, expiration int, access []*token.ResourceActions, signingKey libtrust.PrivateKey) (t *token.Token, expiresIn int, issuedAt *time.Time, err error) {
	joseHeader := &token.Header{
		Type:       "JWT",
		SigningAlg: "RS256",
		KeyID:      signingKey.KeyID(),
	}

	jwtID, err := randString(16)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("Error to generate jwt id: %s", err)
	}

	now := time.Now().UTC()
	issuedAt = &now
	expiresIn = expiration * 60

	claimSet := &token.ClaimSet{
		Issuer:     issuer,
		Subject:    subject,
		Audience:   audience,
		Expiration: now.Add(time.Duration(expiration) * time.Minute).Unix(),
		NotBefore:  now.Unix(),
		IssuedAt:   now.Unix(),
		JWTID:      jwtID,
		Access:     access,
	}

	var joseHeaderBytes, claimSetBytes []byte

	if joseHeaderBytes, err = json.Marshal(joseHeader); err != nil {
		return nil, 0, nil, fmt.Errorf("unable to marshal jose header: %s", err)
	}
	if claimSetBytes, err = json.Marshal(claimSet); err != nil {
		return nil, 0, nil, fmt.Errorf("unable to marshal claim set: %s", err)
	}

	encodedJoseHeader := base64UrlEncode(joseHeaderBytes)
	encodedClaimSet := base64UrlEncode(claimSetBytes)
	payload := fmt.Sprintf("%s.%s", encodedJoseHeader, encodedClaimSet)

	var signatureBytes []byte
	if signatureBytes, _, err = signingKey.Sign(strings.NewReader(payload), crypto.SHA256); err != nil {
		return nil, 0, nil, fmt.Errorf("unable to sign jwt payload: %s", err)
	}

	signature := base64UrlEncode(signatureBytes)
	t, err = token.NewToken(fmt.Sprintf("%s.%s", payload, signature))
	return
}

func randString(length int) (string, error) {
	const alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rb := make([]byte, length)
	if _, err := rand.Read(rb); err != nil {
		return "", err
	}
	for i, b := range rb {
		rb[i] = alphanum[int(b)%len(alphanum)]
	}

	return string(rb), nil
}

func base64UrlEncode(b []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}
