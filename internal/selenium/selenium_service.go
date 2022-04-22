package selenium

import (
	"errors"
	"fmt"
	"github.com/tebeka/selenium"
	"log"
	"os"
)

type Element struct {
	webElement selenium.WebElement
}

type SeleniumService struct {
	debug           bool
	seleniumPath    string
	geckoDriverPath string
	port            int
	service         *selenium.Service
	webDriver       selenium.WebDriver
}

func NewSeleniumService(debug bool, seleniumPath string, geckoDriverPath string, port int) *SeleniumService {
	return &SeleniumService{
		debug:           debug,
		seleniumPath:    seleniumPath,
		geckoDriverPath: geckoDriverPath,
		port:            port,
		//service
		//webDriver
	}
}

func (s *SeleniumService) Start() {
	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(s.geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),              // Output debug information to STDERR.
	}

	selenium.SetDebug(s.debug)
	service, err := selenium.NewSeleniumService(s.seleniumPath, s.port, opts...)
	if err != nil {
		log.Println(err)
	}
	s.service = service

	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", s.port))
	if err != nil {
		log.Println(err)
	}
	s.webDriver = wd
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

func (e *Element) SendKeys(keys string) error {
	return e.webElement.SendKeys(keys)
}

func (e *Element) Click() error {
	return e.webElement.Click()
}
