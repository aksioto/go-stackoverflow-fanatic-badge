package usecase

import (
	"github.com/aksioto/go-stackoverflow-fanatic-badge/internal/selenium"
	"github.com/aksioto/go-stackoverflow-fanatic-badge/utils"
	"log"
)

type BadgeUsecase struct {
	seleniumService *selenium.SeleniumService
	url             string
	email           string
	pass            string
}

func NewBadgeUsecase(seleniumService *selenium.SeleniumService, url, email, pass string) *BadgeUsecase {
	return &BadgeUsecase{
		seleniumService: seleniumService,
		url:             url,
		email:           email,
		pass:            pass,
	}
}

func (u *BadgeUsecase) GoBrrr() {
	u.seleniumService.Start()

	simpleFlowJobs := []PipelineJob{
		PipelineJob(func(el *selenium.Element) (*selenium.Element, error) {
			return nil, u.seleniumService.OpenUrl(u.url)
		}),
		PipelineJob(func(el *selenium.Element) (*selenium.Element, error) {
			return u.seleniumService.FindElementByCssSelector("#email")
		}),
		PipelineJob(func(el *selenium.Element) (*selenium.Element, error) {
			return nil, el.SendKeys(u.email)
		}),
		PipelineJob(func(el *selenium.Element) (*selenium.Element, error) {
			return u.seleniumService.FindElementByCssSelector("#password")
		}),
		PipelineJob(func(el *selenium.Element) (*selenium.Element, error) {
			return nil, el.SendKeys(u.pass)
		}),
		PipelineJob(func(el *selenium.Element) (*selenium.Element, error) {
			return u.seleniumService.FindElementByCssSelector("#submit-button")
		}),
		PipelineJob(func(el *selenium.Element) (*selenium.Element, error) {
			return nil, el.Click()
		}),
		PipelineJob(func(el *selenium.Element) (*selenium.Element, error) {
			return u.seleniumService.FindElementByCssSelector(".s-user-card")
		}),
		PipelineJob(func(el *selenium.Element) (*selenium.Element, error) {
			return nil, el.Click()
		}),
	}

	if err := executePipeline(simpleFlowJobs...); err != nil {
		if u.hereWeGoAgain(3, simpleFlowJobs...) {
			// TODO: implement email notification
			log.Fatal("Ah shit! ", err.Error())
		}
	}
	u.seleniumService.Stop()
}

func (u *BadgeUsecase) hereWeGoAgain(attempts int, jobs ...PipelineJob) bool {
	u.restartSelenium()

	for i := 0; i < attempts; i++ {
		if err := executePipeline(jobs...); err == nil {
			return false
		}
	}
	return true
}

func (u *BadgeUsecase) restartSelenium() {
	u.seleniumService.Stop()
	utils.SleepRandomTime(60, 90)
	u.seleniumService.Start()
}

//TODO: move this
type PipelineJob func(element *selenium.Element) (*selenium.Element, error)

func executePipeline(jobs ...PipelineJob) error {
	var element *selenium.Element
	for _, job := range jobs {
		if el, err := job(element); err != nil {
			return err
		} else {
			element = el
		}
		utils.SleepRandomTime(1, 10)
	}
	return nil
}
