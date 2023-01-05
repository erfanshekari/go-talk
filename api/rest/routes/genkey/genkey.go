package genkey

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"

	"github.com/erfanshekari/go-talk/utils/encypt"
	"github.com/labstack/echo/v4"
)

type GenKeyResponse struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

func GenKey(g *echo.Group) {
	g.POST("/genkey", func(c echo.Context) error {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return c.JSON(http.StatusBadRequest, nil)
		}
		return c.JSON(200, GenKeyResponse{
			PublicKey:  encypt.ExportPubKeyAsPEMStr(&privateKey.PublicKey),
			PrivateKey: encypt.ExportPrvKeyAsPEMStr(privateKey),
		})
	})
}
