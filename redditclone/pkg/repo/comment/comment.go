package comment

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/user"
)

type Comment struct {
	TimeCreated string     `json:"created"`
	ID          string     `json:"id"`
	Author      *user.User `json:"author"`
	Body        string     `json:"body"`
}

type BodyComment struct {
	Body string `json:"comment" valid:"required,length(1|1000)"`
}

func (cb *BodyComment) UnmarshalJSON(data []byte) error {
	type Alias BodyComment
	comment := &Alias{}
	if err := json.Unmarshal(data, comment); err != nil {
		return err
	}
	*cb = BodyComment(*comment)
	return nil
}

func (cb *BodyComment) Validate() []string {
	_, err := govalidator.ValidateStruct(cb)
	if err == nil {
		return nil
	}
	errors := make([]string, 0)
	if errs, ok := err.(govalidator.Errors); ok {
		for _, e := range errs {
			errors = append(errors, e.Error())
		}
	}
	return errors
}
