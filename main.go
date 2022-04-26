package main

import (
	"github.com/aksioto/go-stackoverflow-fanatic-badge/internal/selenium"
	"github.com/aksioto/go-stackoverflow-fanatic-badge/usecase"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log"
)

type SeleniumConfig struct {
	Debug            bool   `env:"DEBUG,required"`
	Output           bool   `env:"SELENIUM_OUTPUT,required"`
	SeleniumPath     string `env:"SELENIUM_PATH,required"`
	GeckoDriverPath  string `env:"GECKO_DRIVER_PATH,required"`
	ChromeDriverPath string `env:"CHROME_DRIVER_PATH,required"`
	Port             int    `env:"PORT,required"`
	Capabilities     string `env:"CAPABILITIES,required"`
	BrowserName      string `env:"BROWSER_NAME,required"`
}

type StackoverflowConfig struct {
	Url    string `env:"SO_URL,required"`
	UrlAlt string `env:"SO_URL_ALT,required"`
	Email  string `env:"SO_EMAIL,required"`
	Pass   string `env:"SO_PASS,required"`
}

type Config struct {
	*SeleniumConfig
	*StackoverflowConfig
}

func main() {
	cfg := &Config{
		SeleniumConfig:      &SeleniumConfig{},
		StackoverflowConfig: &StackoverflowConfig{},
	}
	err := igniteConfig(cfg)
	if err != nil {
		log.Fatal("Error happened on IgniteConfig", err)
	}

	// services
	seleniumService := selenium.NewSeleniumService(cfg.BrowserName, cfg.SeleniumPath, cfg.GeckoDriverPath, cfg.ChromeDriverPath, cfg.Port, cfg.Capabilities, cfg.Debug, cfg.Output)
	// usecase
	badgeUsecase := usecase.NewBadgeUsecase(seleniumService, cfg.Url, cfg.UrlAlt, cfg.Email, cfg.Pass)

	badgeUsecase.GoBrrr()

	log.Println("All done")
}

func igniteConfig(appConfig interface{}) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load config, error: ", err.Error())
		return errors.WithStack(err)
	}

	err = env.Parse(appConfig)
	if err != nil {
		log.Fatal("Failed to parse config, error: ", err.Error())
		return errors.WithStack(err)
	}
	return nil
}
