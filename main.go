package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	seleniumPath    = "binaries/selenium-server.jar"
	geckoDriverPath = "binaries/geckodriver.exe"
	port            = 8080
)
const (
	soUrl  = "https://stackoverflow.com/users/login"
	soUser = ""
	soPass = ""
)

var service *SeleniumService

func main() {
	//TODO: implement env vars
	Init()
}

func Init() {
	service = StartSelenium(seleniumPath, geckoDriverPath, port)

	simpleFlowJobs := []PipelineJob{
		PipelineJob(func(el *Element) (*Element, error) {
			return nil, service.OpenUrl(soUrl)
		}),
		PipelineJob(func(el *Element) (*Element, error) {
			return service.FindElementByCssSelector("#email")
		}),
		PipelineJob(func(el *Element) (*Element, error) {
			return nil, el.SendKeys(soUser)
		}),
		PipelineJob(func(el *Element) (*Element, error) {
			return service.FindElementByCssSelector("#password")
		}),
		PipelineJob(func(el *Element) (*Element, error) {
			return nil, el.SendKeys(soPass)
		}),
		PipelineJob(func(el *Element) (*Element, error) {
			return service.FindElementByCssSelector("#submit-button")
		}),
		PipelineJob(func(el *Element) (*Element, error) {
			return nil, el.Click()
		}),
		PipelineJob(func(el *Element) (*Element, error) {
			return service.FindElementByCssSelector(".s-user-card")
		}),
		PipelineJob(func(el *Element) (*Element, error) {
			return nil, el.Click()
		}),
	}

	if err := ExecutePipeline(simpleFlowJobs...); err != nil {
		if HereWeGoAgain(3, simpleFlowJobs...) {
			// TODO: implement email notification
			fmt.Println("Ah shit!")
		}
	}
	service.Stop()
}

func RestartSelenium() {
	service.Stop()
	SleepRandomTime(60, 90)
	service = StartSelenium(seleniumPath, geckoDriverPath, port)
}

type PipelineJob func(element *Element) (*Element, error)

func ExecutePipeline(jobs ...PipelineJob) error {
	var element *Element
	for _, job := range jobs {
		if el, err := job(element); err != nil {
			return err
		} else {
			element = el
		}
		SleepRandomTime(1, 10)
	}
	return nil
}

func HereWeGoAgain(attempts int, jobs ...PipelineJob) bool {
	RestartSelenium()

	for i := 0; i < attempts; i++ {
		if err := ExecutePipeline(jobs...); err == nil {
			return false
		}
	}
	return true
}

func SleepRandomTime(min, max int) {
	rand.Seed(time.Now().UnixNano())
	duration := rand.Intn(max-min+1) + min
	time.Sleep(time.Duration(duration) * time.Second)
}
