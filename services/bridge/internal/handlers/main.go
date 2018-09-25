package handlers

import (
	"github.com/quantadex/stellar_go/clients/federation"
	"github.com/quantadex/stellar_go/clients/horizon"
	"github.com/quantadex/stellar_go/clients/stellartoml"
	"github.com/quantadex/stellar_go/services/bridge/internal/config"
	"github.com/quantadex/stellar_go/services/bridge/internal/db"
	"github.com/quantadex/stellar_go/services/bridge/internal/listener"
	"github.com/quantadex/stellar_go/services/bridge/internal/submitter"
	"github.com/quantadex/stellar_go/support/http"
)

// RequestHandler implements bridge server request handlers
type RequestHandler struct {
	Config               *config.Config                          `inject:""`
	Client               http.SimpleHTTPClientInterface          `inject:""`
	Horizon              horizon.ClientInterface                 `inject:""`
	Database             db.Database                             `inject:""`
	StellarTomlResolver  stellartoml.ClientInterface             `inject:""`
	FederationResolver   federation.ClientInterface              `inject:""`
	TransactionSubmitter submitter.TransactionSubmitterInterface `inject:""`
	PaymentListener      *listener.PaymentListener               `inject:""`
}

func (rh *RequestHandler) isAssetAllowed(code string, issuer string) bool {
	for _, asset := range rh.Config.Assets {
		if asset.Code == code && asset.Issuer == issuer {
			return true
		}
	}
	return false
}
