package translator

import (
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"strings"
)

func InitTranslator() (v *validator.Validate, t ut.Translator) {
	// 中文翻译器
	zhCh := zh.New()
	uni := ut.New(zhCh)
	trans, _ := uni.GetTranslator("zh")
	// 实例化验证对象
	validate := validator.New()
	// 验证器注册翻译器
	_ = zhtranslations.RegisterDefaultTranslations(validate, trans)
	return validate, trans
}
func ReErr(structs any) error {
	validate, trans := InitTranslator()
	errs := validate.Struct(structs)
	if errs != nil {
		var sliceErrs []string
		for _, err := range errs.(validator.ValidationErrors) {
			//翻译错误信息
			sliceErrs = append(sliceErrs, err.Translate(trans))
		}
		return errors.New(strings.Join(sliceErrs, ","))
	}
	return nil
}
