package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type conf struct {
	RateLimitToken     int    `mapstructure:"RATE_LIMIT_TOKEN"`
	RateLimitTokenTime int    `mapstructure:"RATE_LIMIT_TOKEN_TIME"`
	RateLimitIp        int    `mapstructure:"RATE_LIMIT_IP"`
	RateLimitIpTime    int    `mapstructure:"RATE_LIMIT_IP_TIME"`
	JWTSecret          string `mapstructure:"JWT_SECRET"`
	JwtExperesIn       int    `mapstructure:"JWT_EXPIRESIN"`
	TokenAuth          *jwtauth.JWTAuth
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
	return cfg, err
}
