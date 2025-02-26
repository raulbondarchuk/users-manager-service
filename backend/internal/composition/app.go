package composition

import (
	"app/pkg/config"
	"fmt"

	"github.com/spf13/viper"
)

// Application initialization

func Run() {

	config_init()

	fmt.Println(config.ENV().DBHost)
	fmt.Println(config.ENV().DBPort)
	fmt.Println(config.ENV().DBUser)
	fmt.Println(config.ENV().DBPassword)
	fmt.Println(viper.Get("database.schema"))
}
