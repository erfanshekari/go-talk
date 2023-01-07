package genkey

import (
	"net/http"

	uencrypt "github.com/erfanshekari/go-talk/utils/encrypt"
	"github.com/labstack/echo/v4"
)

type GenKeyResponse struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

func GenKey(g *echo.Group) {
	g.POST("/genkey", func(c echo.Context) error {
		privateKey, err := uencrypt.GenerateKey()
		if err != nil {
			return c.JSON(http.StatusBadRequest, nil)
		}
		return c.JSON(200, GenKeyResponse{
			PublicKey:  uencrypt.ExportPubKeyAsPEMStr(&privateKey.PublicKey),
			PrivateKey: uencrypt.ExportPrvKeyAsPEMStr(privateKey),
		})
	})
}
