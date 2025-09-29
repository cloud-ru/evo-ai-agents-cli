package config

type IAMConfig struct {
	BffHost      string `env:"BFF_HOST" envDefault:"<HOST>"`
	ClientID     string `env:"CLIENT_ID,required" envDefault:"<CLIENT_ID>"`
	ClientSecret string `env:"CLIENT_SECRET,required,notEmpty,unset" envDefault:"<CLIENT_SECRET>"`
}
