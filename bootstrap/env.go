package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct{
	MONGO_URI string `mapstructure:"MONGO_URI"`
	PORT string `mapstructure:"PORT"`
	DB_NAME string `mapstructure:"DB_NAME"`
	SECRET_KEY string `mapstructure:"SECRET_KEY"`
}

func LoadEnv() *Env{
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&env)
	if err != nil{
		log.Fatal(err)
	}

	return &env
}