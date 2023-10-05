package actionsrv

import (
	"errors"
	"log"
	"regexp"

	"github.com/dieklingel/core/internal/core"
	"gorm.io/gorm"
)

type ActionService struct {
	database *gorm.DB
}

func NewService(db *gorm.DB) core.ActionService {
	db.AutoMigrate(&core.Action{})

	return &ActionService{
		database: db,
	}
}

func (service *ActionService) Actions() []core.Action {
	var actions []core.Action
	if res := service.database.Find(&actions); res.Error != nil {
		log.Print(res.Error.Error())
	}

	return actions
}

func (service *ActionService) GetActionById(id int) *core.Action {
	var action core.Action
	err := service.database.First(&action, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return &action
}

func (service *ActionService) SaveAction(action core.Action) core.Action {
	service.database.Save(&action)
	return action
}

func (service *ActionService) RemoveAction(action core.Action) core.Action {
	service.database.Delete(&action)
	return action
}

func (service *ActionService) OnActionSaved(handler func(action core.Action)) {
	panic("not implemented")
}

func (service *ActionService) OnActionRemoved(handler func(action core.Action)) {
	panic("not implemented")
}

func (service *ActionService) Execute(pattern string, env map[string]string) []core.ActionExecutionResult {
	var actions []core.Action
	service.database.Find(&actions)

	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Printf("error while compiling regex: %s", err.Error())
		return make([]core.ActionExecutionResult, 0)
	}

	results := make([]core.ActionExecutionResult, 0)
	for _, action := range actions {
		match := regex.MatchString(action.Trigger)
		if match {
			result := action.Execute(env)
			results = append(results, result)
		}
	}

	return results
}
