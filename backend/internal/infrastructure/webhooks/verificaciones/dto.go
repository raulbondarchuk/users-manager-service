package verificaciones

import (
	"app/internal/application/ports"
	"app/pkg/config"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type verificacionesClient struct {
	client *resty.Client

	baseURL string

	loginRoute                 string
	getCompanyByCompanyIdRoute string
	getCompanyByICCIDRoute     string
	checkUserExistsRoute       string

	tokenCacheDuration time.Duration

	systemUsername string
	systemPassword string

	// Internal appToken cache
	token    string
	tokenExp time.Time
}

var _ ports.VerificacionesService = (*verificacionesClient)(nil)

func NewVerificacionesClient() ports.VerificacionesService {
	return &verificacionesClient{
		client: resty.New(),

		baseURL:                    viper.GetString("webhooks.verificaciones.url"),
		loginRoute:                 viper.GetString("webhooks.verificaciones.api.login"),
		getCompanyByCompanyIdRoute: viper.GetString("webhooks.verificaciones.api.getCompanyByCompanyId"),
		getCompanyByICCIDRoute:     viper.GetString("webhooks.verificaciones.api.getCompanyByICCID"),
		checkUserExistsRoute:       viper.GetString("webhooks.verificaciones.api.checkIfUserExists"),

		tokenCacheDuration: viper.GetDuration("webhooks.verificaciones.token_cache_duration"),

		systemUsername: config.ENV().VERIFICACIONES_USERNAME,
		systemPassword: config.ENV().VERIFICACIONES_PASSWORD,
	}
}
