package cloudkey

import (
	"crypto"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestJsonWebEncryption_Decrypt(t *testing.T) {
	const wrapPrivateKey = `{
"alg": "RSA-OAEP-256",
"d": "TqFSJ_777hW5i4Y61dErvULke4rIu_TAQyf0FBRva_EhNhayxkHVLdpiJ0mxbDgRBG1L4umPskdmJrZe2GPUsVwQEqtuymz-tvUd3e253jRARU1kLKw19zgsdOGlvC_IyAqgG5Z87ch7V-r5UZ6A5cLPg6PIK5IFv_ObdCp1ilL3GAz7MQjvBQrdtitPDPfv1X_b4zIi0pCWr7o8nDXcytx8jeP6OPuij60V4H7_yYBHpshisz7jnbaxo6M9cXhsnkhHGm0Ak1M3rRLv_SZRZl1PSZFF6rvva893Ti8wk7HrrWLNeqtssMvuaABc2OMglLlwzide0U79XL9onjIDqQ",
"dp": "2dYt1UHQ3DGRPuLB4fAbh5QnczxTx4q1ZWxaplCpk2FwpM8jHwj3-RFINUKcIjikcRy_1_omDxmGr_P39dI2WdSqU2JXfm-o2yd8xWyY-GaU5wtKNSwjNX7qTLBGGPrR_GnfijqlgwfE5Kvb7ZWaLBfP2mVpXYlX7ZFBpAR52xU",
"dq": "vC8g715ET5HBq9a30A5nhkjjmtrW5sUsLj11xL_AtWat6AKyY02VO12F8KbFlKgMAywqNSThVBwKQrK0oGKtfbVv0quQzfHCHv9LQVhYftLZLZfcWHfO84XTtRpDB3fq3jtYBglcya5XiL6gsGEECDNDSQKRoLllpFZR1hcebeE",
"e": "AQAB",
"ext": true,
"key_ops": [
"unwrapKey"
],
"kty": "RSA",
"n": "1jPtOMGVWtshNVuWSoUrZ5UjfNRAc2ie-mZT-EescxRaVdHxjH8dCq80-_XFUCC-0p_CXllyDkKbkIHhphmVDuk3bq3L6PLtxsfW31JrRd2qflFYoCzOXmussNt9yvhM_2cxwUxD9RjgBa-OJVETRMavK0v4ev4iySSrIQCu6BttooyrPf1gMyiz6Z58Dx3Flfyw2HLro3I96QlUp_yEEDcTTwK4OAfl0TRPC_F0Ie7sFobypVBMn7kFP6GFJqqaiV1nkxLE_vDdmjBHnbQHAqA-ZXFZdj17HRKLJsXBi5mTM-WhKdb3qTZ4tn_wgXpQ4Nc7ZEi-_VCRbHGekj4YSw",
"p": "9k41kCZ4SLogI0fPVMhddKk1M2H7-K3l3xg2bOURCi43YKUc9OER6gu3--Dgmdz1_-GwJkZksDTQVTNmqRkYq_s_kvI0xiWA3Wa0juUK00ZPge5SNJYrRehp-ycrr8L9_c8HDu3hsPbyCVBOjX_uw8YkI2RfYMya5bDL5cYz4qM",
"q": "3qI_sVLiumULlDvVqdnrNRwM9PUL5LKnTkq8ObLxQRVKp7yKXsEieqSPb5RMNh0GUXur75S-EXmILH-BzE_-YAZr4SVDdYWEok49QzW-ogTOXpyKavlS_rWq-2y_lmPLFk2hGptANni2l2vd-5vB29cEJFTuxMTuvjFiO0349jk",
"qi": "snHYbKLsw7VgWwkivtjfb4PP6cHdnBxNji9EiSw715izNeqXNI9Qd5FqvKN4P3tPJHgxQOpZiUMFwr_qj_yEmyfgpKUGsSrFUo75KPzr6mjBKEZ8vg5tz91w4BR5vAgioLWgHU-L7G34DqLJVcOixoeEhFmZm6-6j6yDLz1Y3YE"
}`
	jweString := strings.Join([]string{
		"eyJhbGciOiJSU0EtT0FFUC0yNTYiLCJlbmMiOiJBMjU2R0NNIn0",
		"CA7xfRbDAFqX6j7Y-ySFONDikx25AYMwcU1mW1X0OKDDDhRPpeTxISq5mfPOEDkGPvC8KrAZZwJmqJPc5eTnO7szNoZhxuUt0LZjfc-5q294kv8wpoq3ZArux0BcdGjeUSXxh-wUdQqoDSHJjk-gTCCuvQ3IHMvoarWlcYUnCzumvWrohjQad_aal-NMckTTMoy1t3NGr0h7MzhHtaNbDujpjBiYr0HVV3qPGEAwSi1kLZgdOnEVg7TRPUabv-7eWis2hCj1Zt3nkKGZ3ONm0p205EwfN1ihyAKNob2ZyWmR2iuebHRRSjAXIkY-V4SQKEqYQKPbpjzUqBv_X2uvGg",
		"A1SXjP4nxWi17M5G",
		"D4jsiQScLjBtwQ",
		"2ULr-4uoL0Z5t4v1LBVmFw",
	}, ".")

	privatekeyJwk := &JsonWebKey{}
	err := json.Unmarshal([]byte(wrapPrivateKey), privatekeyJwk)
	require.NoError(t, err)
	privateKey := privatekeyJwk.PrivateKey().(crypto.Decrypter)

	// Decrypt the JWE
	jwe, err := NewJsonWebEncryption(jweString)
	require.NoError(t, err)
	decrypted, _, err := jwe.Decrypt(func(header *JoseHeader) (crypto.PrivateKey, error) {
		return privateKey, nil
	})
	require.NoError(t, err)
	assert.Equal(t, "plain text", string(decrypted))
}

func mustDecodeBase64URL(s string) []byte {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func TestJsonWebEncryption_KDF(t *testing.T) {
	alicePrivateKey := JsonWebKey{
		Curve:   CurveNameP384,
		KeyType: KeyTypeEC,
		D:       mustDecodeBase64URL("nD9ydFl6wqTzzQRE7DY50W2HoA3Ipk-qr2GdSTIzasADDIe3fsv30XWOUlCV51-w"),
		X:       mustDecodeBase64URL("1TOe7fnQIvy_7stDWY9FTWMB2brL-nWLKxOf1N9ysC9aUriKzGVG41TR4uf1_apU"),
		Y:       mustDecodeBase64URL("g1netoRQTLqVbmYD8EVfqpDne0qcTMdvOiotYdL6Su-D7QOHopG0HgTmWwVhWT-J"),
	}
	bobPublicKey := JsonWebKey{
		Curve:   CurveNameP384,
		KeyType: KeyTypeEC,
		X:       mustDecodeBase64URL("n0xM_GPHtVgYod4XUa7KXbuiEoJek1lYRQBetjA8SraIbmKI6baEI5w-lmKAgOR7"),
		Y:       mustDecodeBase64URL("vuL2ItFJtQaxRe6lkm84m-GC-7oOqPGoSOb8PNNZdxpNvKAnVNLkWRn3K97O3IGZ"),
	}
	expetedSharedSecret := mustDecodeBase64URL("jytQHNfFwnNCpeaqsDetKP7UW7H6LFijv3zmoTGpvSs")

	bobEcdhKey, err := bobPublicKey.PublicKey().(*ecdsa.PublicKey).ECDH()
	require.NoError(t, err)
	aliceEcdhKey, err := alicePrivateKey.PrivateKey().(*ecdsa.PrivateKey).ECDH()
	require.NoError(t, err)
	z, err := aliceEcdhKey.ECDH(bobEcdhKey)
	require.NoError(t, err)
	z = z[:32]
	kdf := &ecdhesKDF{
		z:   z,
		alg: "A256GCM",
	}
	assert.DeepEqual(t, expetedSharedSecret, kdf.getAESGCM256DerivedKey())
}
