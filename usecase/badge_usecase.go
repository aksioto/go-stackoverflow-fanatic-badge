package usecase

import (
	"fmt"
	"github.com/aksioto/go-stackoverflow-fanatic-badge/internal/selenium"
	"github.com/aksioto/go-stackoverflow-fanatic-badge/utils"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type BadgeUsecase struct {
	seleniumService *selenium.SeleniumService
	pipeline        *PipelineConfig
}

func NewBadgeUsecase(seleniumService *selenium.SeleniumService, jobsConfigPath string) *BadgeUsecase {
	pipeline := &PipelineConfig{}
	if err := parseConfig(pipeline, jobsConfigPath); err != nil {
		log.Fatal(err.Error())
	}

	return &BadgeUsecase{
		seleniumService: seleniumService,
		pipeline:        pipeline,
	}
}

func (u *BadgeUsecase) StartEarnBadge() {
	if err := u.seleniumService.Start(); err != nil {
		log.Fatal("Selenium service not started. Error: ", err.Error())
		return
	}

	if err := u.executePipeline(); err != nil {
		if u.hereWeGoAgain(3) {
			// TODO: implement email notification
			log.Fatal("Ah shit! ", err.Error())
		}
	}
	u.seleniumService.Stop()
}

func (u *BadgeUsecase) hereWeGoAgain(attempts int) bool {
	u.restartSelenium()

	for i := 0; i < attempts; i++ {
		log.Println(fmt.Printf("[hereWeGoAgain] attempts: %o\n", i))
		if err := u.executePipeline(); err == nil {
			return false
		}
	}
	return true
}

func (u *BadgeUsecase) restartSelenium() {
	u.seleniumService.Stop()
	utils.SleepRandomTime(u.pipeline.RestartTimeout.Min, u.pipeline.RestartTimeout.Max)
	_ = u.seleniumService.Start()
}

type PipelineJob func(element *selenium.Element) (*selenium.Element, error)

func (u *BadgeUsecase) executePipeline() error {
	var element *selenium.Element
	for _, job := range u.pipeline.Jobs {
		pj := u.prepareJob(job.Method, job.Args)
		if el, err := pj(element); err != nil {
			return err
		} else {
			element = el
		}
		utils.SleepRandomTime(u.pipeline.JobTimeout.Min, u.pipeline.JobTimeout.Max)
	}
	return nil
}

func (u *BadgeUsecase) prepareJob(method, args string) PipelineJob {
	args = parseArgs(args)
	var pipelineJob PipelineJob
	switch method {
	case "OpenUrl":
		pipelineJob = func(el *selenium.Element) (*selenium.Element, error) {
			return nil, u.seleniumService.OpenUrl(args)
		}
	case "FindElementByCssSelector":
		pipelineJob = func(el *selenium.Element) (*selenium.Element, error) {
			return u.seleniumService.FindElementByCssSelector(args)
		}
	case "Click":
		pipelineJob = func(el *selenium.Element) (*selenium.Element, error) {
			return nil, el.Click()
		}
	case "SendKeys":
		pipelineJob = func(el *selenium.Element) (*selenium.Element, error) {
			return nil, el.SendKeys(args)
		}
	}

	return pipelineJob
}

type PipelineConfig struct {
	JobTimeout struct {
		Min int `yaml:"min"`
		Max int `yaml:"max"`
	} `yaml:"jobTimeout"`

	RestartTimeout struct {
		Min int `yaml:"min"`
		Max int `yaml:"max"`
	} `yaml:"restartTimeout"`

	Jobs []struct {
		Method string `yaml:"method"`
		Args   string `yaml:"args"`
	} `yaml:"jobs"`
}

func parseConfig(config interface{}, path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Failed to load config, error: ", err.Error())
		return errors.WithStack(err)
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Fatal("Failed to parse config, error: ", err.Error())
		return errors.WithStack(err)
	}
	return nil
}

func parseArgs(args string) string {
	re := regexp.MustCompile(`\$\{([a-zA-Z\d_-]+?)\}`)
	subMatch := re.FindStringSubmatch(strings.ReplaceAll(args, " ", ""))
	if len(subMatch) > 0 {
		return os.Getenv(subMatch[1])
	}

	return args
}
