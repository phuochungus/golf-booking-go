package initialize

import (
	"fmt"
	"golf-booking-go/global"

	"github.com/spf13/viper"
)

func LoadConfig() {
	v := viper.New()
	v.AddConfigPath("./config/")
	v.SetConfigName("local")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error reading config: %s", err.Error()))
	}

	if err := v.Unmarshal(&global.Config); err != nil {
		fmt.Println("Error unmarshalling config:", err.Error())
	}
}
