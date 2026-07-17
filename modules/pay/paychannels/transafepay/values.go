package transafepay

const (
	baseURL          = "https://pay-test.targeted.work"
	payURL           = "/cashin/india/order/create"
	queryPayURL      = "/cashin/india/order/query"
	withdrawURL      = "/cashout/india/order/create"
	queryWithdrawURL = "/cashout/india/order/query"
	mid              = 2061545184
	appID            = "R22e9c0f1x3hycRKd9"
)

var privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA2JMAuacbgswm34+AiM5bQz6tH7IqQSriRwitHu+66Ox9Q1Kx
wUkmMgmmvVOXC0GyLC3+a4TM6xs3n9mJd8dni/AQcfhBmQQuJ49C/CKIaic081Ok
iHFeoVqLLY9I/G98E+oQ/FkxljvZxAiZhV1jxh7kR5Cn+I10Z7501MiE0PSlAkk9
RWk9g6ujlM+6gyJxzWbgsw+ZNiIyUShZ+u6LqoP5AFvEvCVtNyFx5hruUqLaRzMf
eedC8RAlA98C9UVVovSOWODoFCRBBmMN+oviKdICg1KCnG3aS28z/Qu4T8dO/seS
MFkh+zebfDWNC6bBXkk0/AbWOc/A8epUGxh6wwIDAQABAoIBAQCwhAf4SOoUTh50
4foqBr69bQ0fo0/rM2+f48NuTofxpiZmQ7UbS2p1KAadks6t3wrAFXCi42c9TD6m
61+un3iWxsAwk92rQIv5tDstU1wJ4RC4LfkROA2fxzNem8Up8heYB5N2iZ/vybIw
o6201Cwo4y694Xn2DBzPHZ0NfK07vFn+/1QJz3OnRHXpiclEltxfR5qa2QTFtBOF
s8bxvSYRd451cQYRv3G2/JLPGu53grJ3gOACfg3BXlSBur/3YxQa1TKxYgxrOCgh
qGYwKk2J4KrchQYvodlrD7NckbwXoP16EPEKlnbxC3Z5fHk6JWCsbyvbokrXB0nj
7c/2IfshAoGBAPD2O2knUTXznVSzA6Rph8IiS3jvVKElkTfK2dkLMxqc+Rgo91vm
Sz4ls+pSWsO8mZ/xeFVB4UxilHbo2eA898CjrTVfE+3ovLLUkVIfQZO6a8g/SyGj
S0IooiFENEKhYl7VuhrS6YLfxXpJ3nMAcjs30001EOB9UabzgVYBKW9xAoGBAOYX
I0yuWvxPltmmE8pOnEcNsm6mQxeMjZ4tCEYQ0zgLyQ8CU84nP5Wzax11p2tKKGi2
Zpl9dr9kGekQiFEnxn0p6HDV7cnqbrRUqMzpNekFhg6R9UJYTYbVuZRQdjkXNZr0
xuwtxAqtEaL0Ecb6rCFrQoz1Cp+jr2Y0pTLSIptzAoGBAM7x2nJj9wZL9JpSqFUi
4Rpk08bIXTmNBxwQlmomiTxZBkGupUqDDaoRPWSIWDV1G8dW3g7y5qjwnStGriA8
nFcQ+2g0AOy5jp5Tm38ZIvvaBklTupXoJXFrblr9tZhn/BEQo6X3796bioGNQYQp
JPIg7tOHOLhkSKUGytRvaQOxAoGANttjXh/fr3ImTAPtiflOBx674Oeg7VdkWEby
Rwox4ZXlvL9YTblLzpw1OJPZc5n9omGCeqvQCPeLueZoKbkq+KJqOaS3icFr3XCl
wh0D7meTu5uO9t9otHH2SecZGKlw7V3Yo2yR5d23a7t3FZKkQV96QXMWslLrOZnp
nQE+oQMCgYEAjDYjbmZegOd4c7GwM8lxll1bjPfI1h4y6MozegLL2NmZixzzR0dY
bdEEOQbLDLQnX2sEouo4iD9BzipdErnS/RecFlOPFGQAxdqQPg/KlAOah6CFuXm1
PO442tSxHzlWpsIexYBttt6iPtrkhPVuJqn/61iUoy9xSRA8i3AUzHA=
-----END RSA PRIVATE KEY-----`)

var privateKeyCS8 = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDYkwC5pxuCzCbf
j4CIzltDPq0fsipBKuJHCK0e77ro7H1DUrHBSSYyCaa9U5cLQbIsLf5rhMzrGzef
2Yl3x2eL8BBx+EGZBC4nj0L8IohqJzTzU6SIcV6hWostj0j8b3wT6hD8WTGWO9nE
CJmFXWPGHuRHkKf4jXRnvnTUyITQ9KUCST1FaT2Dq6OUz7qDInHNZuCzD5k2IjJR
KFn67ouqg/kAW8S8JW03IXHmGu5SotpHMx9550LxECUD3wL1RVWi9I5Y4OgUJEEG
Yw36i+Ip0gKDUoKcbdpLbzP9C7hPx07+x5IwWSH7N5t8NY0LpsFeSTT8BtY5z8Dx
6lQbGHrDAgMBAAECggEBALCEB/hI6hROHnTh+ioGvr1tDR+jT+szb5/jw25Oh/Gm
JmZDtRtLanUoBp2Szq3fCsAVcKLjZz1MPqbrX66feJbGwDCT3atAi/m0Oy1TXAnh
ELgt+RE4DZ/HM16bxSnyF5gHk3aJn+/JsjCjrbTULCjjLr3hefYMHM8dnQ18rTu8
Wf7/VAnPc6dEdemJyUSW3F9HmprZBMW0E4WzxvG9JhF3jnVxBhG/cbb8ks8a7neC
sneA4AJ+DcFeVIG6v/djFBrVMrFiDGs4KCGoZjAqTYngqtyFBi+h2WsPs1yRvBeg
/XoQ8QqWdvELdnl8eTolYKxvK9uiStcHSePtz/Yh+yECgYEA8PY7aSdRNfOdVLMD
pGmHwiJLeO9UoSWRN8rZ2QszGpz5GCj3W+ZLPiWz6lJaw7yZn/F4VUHhTGKUdujZ
4Dz3wKOtNV8T7ei8stSRUh9Bk7pryD9LIaNLQiiiIUQ0QqFiXtW6GtLpgt/Fekne
cwByOzfTTTUQ4H1RpvOBVgEpb3ECgYEA5hcjTK5a/E+W2aYTyk6cRw2ybqZDF4yN
ni0IRhDTOAvJDwJTzic/lbNrHXWna0ooaLZmmX12v2QZ6RCIUSfGfSnocNXtyeput
FSozOk16QWGDpH1QlhNhtW5lFB2ORc1mvTG7C3ECq0RovQRxvqsIWtCjPUKn6OvZ
jSlMtIim3MCgYEAzvHacmP3Bkv0mlKoVSLhGmTTxshdOY0HHBCWaiaJPFkGQa6l
SoMNqhE9ZIhYNXUbx1beDvLmqPCdK0auIDycVxD7aDQA7LmOnlObfxki+9oGSVO6
leglcWtuWv21mGf8ERCjpffv3puKgY1BhCkk8iDu04c4uGRIpQbK1G9pA7ECgYA2
22NeH9+vciZMA+2J+U4HHrvg56DtV2RYRvJHCjHhleW8v1hNuUvOnDU4k9lzmf2i
YYJ6q9AI94u55mgpuSr4omo5pLeJwWvdcKXCHQPuZ5O7m47232i0cfZJ5xkYqXDt
XdijbJHl3bdru3cVkqRBX3pBcxayUus5memdAT6hAwKBgQCMNiNuZl6A53hzsbAz
yXGWXVuM98jWHjLoyjN6AsvY2ZmLHPNHR1ht0QQ5BssMtCdfawSi6jiIP0HOKl0S
udL9F5wWU48UZADF2pA+D8qUA5qHoIW5ebU87jja1LEfOVamwh7FgG223qI+2uSE
9W4mqf/rWJSjL3FJEDyLcBTMcA==
-----END PRIVATE KEY-----`)

var publicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCzfNscK0HjYa2USHWjdeV9pzs3
i3wRT9ucxzUeMivdjUCfwgBwU8emCXQsCYlbhg190rwupdvGGKwp1osPS19G/3Yp
mc1BY088D/shKlkK5rwHw9lK7Ch5m4a9YrSKbp+LdHGDa3hSUn9X5h4zcpseNX+k
a1IhN2D1we/w0IXRiwIDAQAB
-----END PUBLIC KEY-----`)

var myPublicKeyBase = `MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqSbw9rY58HuqDZ0tyY2xCwQfRiXUI3EasGAMhLbQiMAMBGSZQjcKUXHAnTe2YKJHD/Tn7vQpQtB6/ktuECH60lSVkEph+e6uXgE4nWbRkYur5pkosEHG/fTE+Upo0E9ID1BgElFQ/3gR9iUgOuEiy2EoyJWNyyCrT6OZB/nKSrl7M7EWFr6X84NYZWp1pzKcf8cQIxwzuw9yF3Kwa9L43ANgJItrEqZ4zxeJk0/qJEC20UaE2AmCTkl/HfY6B6j8mNeHSPSyayVcoDitWHiNeUKE3QL2sXog7G+7pljYFZz6M259w2r/zm8VkD5XJt3BGtjdZtIN0Ckf3C+ygOMY/wIDAQAB`

var myPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqSbw9rY58HuqDZ0tyY2x
CwQfRiXUI3EasGAMhLbQiMAMBGSZQjcKUXHAnTe2YKJHD/Tn7vQpQtB6/ktuECH6
0lSVkEph+e6uXgE4nWbRkYur5pkosEHG/fTE+Upo0E9ID1BgElFQ/3gR9iUgOuEi
y2EoyJWNyyCrT6OZB/nKSrl7M7EWFr6X84NYZWp1pzKcf8cQIxwzuw9yF3Kwa9L4
3ANgJItrEqZ4zxeJk0/qJEC20UaE2AmCTkl/HfY6B6j8mNeHSPSyayVcoDitWHiN
eUKE3QL2sXog7G+7pljYFZz6M259w2r/zm8VkD5XJt3BGtjdZtIN0Ckf3C+ygOMY
/wIDAQAB
-----END PUBLIC KEY-----`)

type PayReq struct {
	MchID       int64   `json:"mchId"`
	TxChannel   string  `json:"txChannel"`
	AppID       string  `json:"appId"`
	Timestamp   int64   `json:"timestamp"`
	MchOrderNo  string  `json:"mchOrderNo"`
	BankCode    string  `json:"bankCode"`
	Amount      float64 `json:"amount"`
	Name        string  `json:"name"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	Sign        string  `json:"sign"`
	ProductInfo string  `json:"productInfo"`
	NotifyUrl   string  `json:"notifyUrl"`
	ReturnUrl   string  `json:"returnUrl"`
}

type PayResp struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   struct {
		MchOrderNo  string  `json:"mchOrderNo"`
		PlatOrderNo string  `json:"platOrderNo"`
		ApplyTime   int64   `json:"applyTime"`
		Amount      float64 `json:"amount"`
		Reference   string  `json:"reference"`
		Link        string  `json:"link"`
		ExpireTime  int64   `json:"expireTime"`
	} `json:"data"`
}

type PayCallbackReq struct {
	Status      int     `json:"status"`
	MemberID    string  `json:"memberId"`
	MchOrderNo  string  `json:"mchOrderNo"`
	PlatOrderNo string  `json:"platOrderNo"`
	OrderStatus string  `json:"orderStatus"` // PENDING, SUCCESS, FAILED, REFUND, PART_SUC
	Amount      float64 `json:"amount"`
	Fee         float64 `json:"fee"`
	Msg         string  `json:"msg"`
	TimeEnd     int64   `json:"timeEnd"`
	Sign        string  `json:"sign"`
}

type QueryPayReq struct {
	MchID      int64  `json:"mchId"`
	Timestamp  int64  `json:"timestamp"`
	MchOrderNo string `json:"mchOrderNo"`
	TrxID      string `json:"trxId"`
	Sign       string `json:"sign"`
}

type QueryPayResp struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   struct {
		MchOrderNo  string  `json:"mchOrderNo"`
		PlatOrderNo string  `json:"platOrderNo"`
		OrderStatus string  `json:"orderStatus"` // PENDING, SUCCESS, FAILED, REFUND, PART_SUC
		TrxID       string  `json:"trxId"`
		Amount      float64 `json:"amount"`
		Fee         float64 `json:"fee"`
		TimeEnd     int64   `json:"timeEnd"`
	} `json:"data"`
}

type WithdrawReq struct {
	MchID      int64   `json:"mchId"`
	TxChannel  string  `json:"txChannel"`
	AppID      string  `json:"appId"`
	Timestamp  int64   `json:"timestamp"`
	MchOrderNo string  `json:"mchOrderNo"`
	Name       string  `json:"name"`
	Phone      string  `json:"phone"`
	Email      string  `json:"email"`
	BankCode   string  `json:"bankCode"`
	Account    string  `json:"account"`
	Amount     float64 `json:"amount"`
	NotifyUrl  string  `json:"notifyUrl"`
	Sign       string  `json:"sign"`
	IFSC       string  `json:"ifsc"`
}

type WithdrawResp struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   struct {
		PlatOrderNo string  `json:"platOrderNo"`
		MchOrderNo  string  `json:"mchOrderNo"`
		ApplyTime   int64   `json:"applyTime"`
		Amount      float64 `json:"amount"`
		OrderStatus string  `json:"orderStatus"` // PENDING, SUCCESS, FAILED, REFUND
	} `json:"data"`
}

type WithdrawCallbackReq struct {
	Status      int     `json:"status"`
	MemberID    string  `json:"memberId"`
	MchOrderNo  string  `json:"mchOrderNo"`
	PlatOrderNo string  `json:"platOrderNo"`
	OrderStatus string  `json:"orderStatus"`
	Amount      float64 `json:"amount"`
	Fee         float64 `json:"fee"`
	TimeEnd     int64   `json:"timeEnd"`
	Msg         string  `json:"msg"`
	Sign        string  `json:"sign"`
	UTR         string  `json:"utr"`
}

type QueryWithdrawReq struct {
	MchID      int64  `json:"mchId"`
	Timestamp  int64  `json:"timestamp"`
	MchOrderNo string `json:"mchOrderNo"`
	TrxID      string `json:"trxId"`
	Sign       string `json:"sign"`
}

type QueryWithdrawResp struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   struct {
		MchOrderNo  string  `json:"mchOrderNo"`
		PlatOrderNo string  `json:"platOrderNo"`
		TrxID       string  `json:"trxId"`
		OrderStatus string  `json:"orderStatus"` // PENDING, SUCCESS, FAILED, REFUND
		Amount      float64 `json:"amount"`
		Fee         float64 `json:"fee"`
		TimeEnd     int64   `json:"timeEnd"`
	} `json:"data"`
}
