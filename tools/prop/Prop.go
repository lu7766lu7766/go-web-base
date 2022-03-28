package prop

import (
	"fmt"
	"p4_web/constant"
	"p4_web/tools/exception"

	"github.com/goldeneggg/structil"
)

func Get(obj interface{}, prop string) interface{} {
	finder, err := structil.NewFinder(obj)
	if err != nil {
		panic(exception.ApiException{
			Code:    []int{constant.SYSTEM_ERROR},
			Message: fmt.Sprintf("finder cannot created: %v: %v", obj, prop),
		})
	}
	mapper, err := finder.FindTop(prop).ToMap()
	if err != nil {
		panic(exception.ApiException{
			Code:    []int{constant.SYSTEM_ERROR},
			Message: fmt.Sprintf("property not found: %v", prop),
		})
	}
	return mapper[prop]
}
