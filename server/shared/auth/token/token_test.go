package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const privateKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu1SU1LfVLPHCozMxH2Mo
4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0/IzW7yWR7QkrmBL7jTKEn5u
+qKhbwKfBstIs+bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyeh
kd3qqGElvW/VDL5AaWTg0nLVkjRo9z+40RQzuVaE8AkAFmxZzow3x+VJYKdjykkJ
0iT9wCS0DRTXu269V264Vf/3jvredZiKRkgwlL9xNAwxXFg0x/XFw005UWVRIkdg
cKWTjpBP2dPwVZ4WWC+9aGVd+Gyn1o0CLelf4rEjGoXbAAEgAqeGUxrcIlbjXfbc
mwIDAQAB
-----END PUBLIC KEY-----`

func TestVerify(t *testing.T) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(privateKey))
	if err != nil {
		t.Fatalf("cannot parse public key: %v", err)
	}

	v := &JWTTokenVerifier{
		PublicKey: pubKey,
	}

	tkn := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiMTIzNDU2Nzg5MCJ9.GcHiRgOuFiiQJAMKJemV2j5Vr8uZslvOJksONETcsXxTpDqEJPHiwLsc94W3cvVpYrJO6O6c8mywxYjOWkk7iBoyEWMmbapsE8T3dDyFRq2xnV-1DZerlTNVuO4gT2fq3eNOEE-XXu0y0zlnCW7LMnOZdstHAkMD-ZQP0vKZuLJjP_AMhfd3BcsVXTMLVKjW0aG-UwkAhsathBa24NaLy2AsCIljSGNjmQ4gp9CihlHDRyUCRxBPuKDf0ym-tBSUgWk9zFugKlx-nSCYLSXgMPJ0CzgSDvmoXkC3HNM1VWOo-qd-QtInMYYuQs_RSPK8VVDj7EV7llHcjbki-OtRPA"
	cases := []struct {
		name    string
		tkn     string
		now     time.Time
		want    string
		wantErr bool
	}{
		{
			name: "valid_token",
			tkn:  tkn,
			now:  time.Unix(1516239122, 0),
			want: "1234567890",
		},
		{
			name:    "token_expired",
			tkn:     tkn,
			now:     time.Unix(1517239122, 0),
			wantErr: true,
		},
		{
			name:    "bad_token",
			tkn:     "bad_token",
			now:     time.Unix(1516239122, 0),
			wantErr: true,
		},
		{
			name:    "wrong_signature",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYsImlhdCI6Mc3ViIjoiMTIzNDU2Nzg5MCJ8.GcHiRgOuFiiQJAMKJemV2j5Vr8uZslvOJksONETcsXxTpDqEJPHiwLsc94W3cvVpYrJO6O6c8mywxYjOWkk7iBoyEWMmbapsE8T3dDyFRq2xnV-1DZerlTNVuO4gT2fq3eNOEE-XXu0y0zlnCW7LMnOZdstHAkMD-ZQP0vKZuLJjP_AMhfd3BcsVXTMLVKjW0aG-UwkAhsathBa24NaLy2AsCIljSGNjmQ4gp9CihlHDRyUCRxBPuKDf0ym-tBSUgWk9zFugKlx-nSCYLSXgMPJ0CzgSDvmoXkC3HNM1VWOo-qd-QtInMYYuQs_RSPK8VVDj7EV7llHcjbki-OtRPA",
			now:     time.Unix(1516239122, 0),
			wantErr: true,
		},
	}

	for _, c := range cases {
		// t.Run不会并行执行，这样才能每次重置时间
		t.Run(c.name, func(t *testing.T) {
			jwt.TimeFunc = func() time.Time {
				return c.now
			}

			accountID, err := v.Verify(c.tkn)
			if !c.wantErr && err != nil {
				t.Errorf("verification failed: %v", err) // Errorf下面的代码还会继续执行，才可以打印出accountID
			}

			if c.wantErr && err == nil {
				t.Errorf("want err; got no error")
			}

			if accountID != c.want {
				t.Fatalf("wrong account id. want: %q, got: %q", c.want, accountID)
			}
		})
	}
}
