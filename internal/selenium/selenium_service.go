package selenium

import (
	"errors"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium/firefox"
	"log"
	"os"
	"strings"
)

const (
	CHROME  = "chrome"
	FIREFOX = "firefox"
)

type Element struct {
	webElement selenium.WebElement
}

type SeleniumService struct {
	debug            bool
	output           bool
	capabilities     []string
	browserName      string
	seleniumPath     string
	geckoDriverPath  string
	chromeDriverPath string
	port             int
	service          *selenium.Service
	webDriver        selenium.WebDriver
}

func NewSeleniumService(browserName string, seleniumPath string, geckoDriverPath string, chromeDriverPath string, port int, capabilities string, debug bool, output bool) *SeleniumService {
	return &SeleniumService{
		debug:            debug,
		output:           output,
		capabilities:     strings.Split(capabilities, ","),
		browserName:      browserName,
		seleniumPath:     seleniumPath,
		geckoDriverPath:  geckoDriverPath,
		chromeDriverPath: chromeDriverPath,
		port:             port,
		//service
		//webDriver
	}
}

func (s *SeleniumService) Start() error {
	opts := []selenium.ServiceOption{
		//selenium.StartFrameBuffer(), // Start an X frame buffer for the browser to run in.
		s.getDriver(),
	}

	if s.output {
		opts = append(opts, selenium.Output(os.Stderr)) // Output debug information to STDERR.
	}

	selenium.SetDebug(s.debug)
	service, err := selenium.NewSeleniumService(s.seleniumPath, s.port, opts...)
	if err != nil {
		log.Println(err)
		return err
	}
	s.service = service

	wd, err := selenium.NewRemote(s.getCaps(), fmt.Sprintf("http://localhost:%d/wd/hub", s.port))
	if err != nil {
		log.Println(err)
		return err
	}
	s.webDriver = wd

	return nil
}

func (s *SeleniumService) Stop() {
	defer s.service.Stop()
	defer s.webDriver.Quit()
}

func (s *SeleniumService) OpenUrl(url string) error {
	if err := s.webDriver.Get(url); err != nil {
		return err
	}
	return nil
}

func (s *SeleniumService) FindElementByCssSelector(selector string) (*Element, error) {
	el, err := s.webDriver.FindElement(selenium.ByCSSSelector, selector)
	if err != nil {
		log.Println(err)

		if s.IsRecaptcha() {
			return &Element{webElement: nil}, errors.New("ReCaptcha")
		}
	}

	return &Element{webElement: el}, err
}

func (s *SeleniumService) IsRecaptcha() bool {
	const captcha = "#nocaptcha-form"

	el, err := s.webDriver.FindElement(selenium.ByCSSSelector, captcha)
	if err != nil {
		log.Println(err)
	}

	return el != nil
}

func (s *SeleniumService) getDriver() selenium.ServiceOption {
	var driver selenium.ServiceOption

	switch s.browserName {
	case CHROME:
		driver = selenium.ChromeDriver(s.chromeDriverPath)
	case FIREFOX:
		driver = selenium.GeckoDriver(s.geckoDriverPath)
	default:
		log.Fatal(fmt.Sprintf("Wrong browser name in .env file. Use '%s' or '%s'", CHROME, FIREFOX))
	}

	return driver
}

func (s *SeleniumService) getCaps() selenium.Capabilities {
	caps := selenium.Capabilities{"browserName": s.browserName}

	switch s.browserName {
	case CHROME:
		caps.AddChrome(
			chrome.Capabilities{
				Args: s.capabilities,
			},
		)
	case FIREFOX:
		caps.AddFirefox(
			firefox.Capabilities{
				Args: s.capabilities,
			},
		)
	default:
		log.Fatal(fmt.Sprintf("Wrong browser name in .env file. Use '%s' or '%s'", CHROME, FIREFOX))
	}

	return caps
}

func (e *Element) SendKeys(keys string) error {
	return e.webElement.SendKeys(keys)
}

func (e *Element) Click() error {
	return e.webElement.Click()
}
