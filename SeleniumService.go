package main

import (
	"errors"
	"fmt"
	"github.com/tebeka/selenium"
	"os"
)

type SeleniumService struct {
	service   *selenium.Service
	webDriver selenium.WebDriver
}

type Element struct {
	webElement selenium.WebElement
}

func StartSelenium(seleniumPath, geckoDriverPath string, port int) *SeleniumService {
	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),            // Output debug information to STDERR.
	}

	//selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		fmt.Println(err)
	}
	//defer service.Stop()
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		fmt.Println(err)
	}
	//defer wd.Quit()
	return &SeleniumService{
		service:   service,
		webDriver: wd,
	}
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
		fmt.Println(err)

		if s.IsRecaptcha() {
			return &Element{webElement: nil}, errors.New("ReCaptcha")
		}
	}

	return &Element{webElement: el}, err
}

func (e *Element) SendKeys(keys string) error {
	return e.webElement.SendKeys(keys)
}
func (e *Element) Click() error {
	return e.webElement.Click()
}

func (s *SeleniumService) IsRecaptcha() bool {
	const captcha = "#nocaptcha-form"

	el, err := s.webDriver.FindElement(selenium.ByCSSSelector, captcha)
	if err != nil {
		fmt.Println(err)
	}

	return el != nil
}
