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

type Config struct {
	*SeleniumConfig
	JobsConfigPath string `env:"JOBS_CONFIG,required"`
}

func main() {
	cfg := &Config{
		SeleniumConfig: &SeleniumConfig{},
	}
	err := igniteConfig(cfg)
	if err != nil {
		log.Fatal("Error happened on IgniteConfig", err)
	}

	// services
	seleniumService := selenium.NewSeleniumService(cfg.BrowserName, cfg.SeleniumPath, cfg.GeckoDriverPath, cfg.ChromeDriverPath, cfg.Port, cfg.Capabilities, cfg.Debug, cfg.Output)
	// usecase
	badgeUsecase := usecase.NewBadgeUsecase(seleniumService, cfg.JobsConfigPath)

	badgeUsecase.StartEarnBadge()

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
