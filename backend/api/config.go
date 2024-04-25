package api

type Config struct {
	Port int
	Env  string
	DB   struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Limiter struct {
		RPS     float64
		Burst   int
		Enabled bool
	}
	Cors struct {
		TrustedOrigins []string
	}
	Providers struct {
		Github struct {
			ClientId     string
			ClientSecret string
			Url          string
			RedirectUrl  string
			UserUrl      string
		}
	}
}
