package models

import (
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"gorm.io/gorm"
	"online-voice-channel/dao"
	"strings"
)

// Todo Model
type Todo struct {
	ID      uint           `json:"id"`
	Title   string         `json:"title" validate:"required"`
	Status  bool           `json:"status" validate:"-"`
	Deleted gorm.DeletedAt `json:"deleted"`
}

func Init() (v *validator.Validate, t ut.Translator) {
	// 中文翻译器
	zh_ch := zh.New()
	uni := ut.New(zh_ch)
	trans, _ := uni.GetTranslator("zh")
	// 实例化验证对象
	validate := validator.New()
	// 验证器注册翻译器
	_ = zhtranslations.RegisterDefaultTranslations(validate, trans)
	return validate, trans
}

/*
	Todo这个Model的增删改查操作都放在这里
*/
// CreateTodo 创建todo
func CreateTodo(todo *Todo) (err error) {
	//err = dao.DB.Create(&todo).Error
	//return
	validate, trans := Init()
	errs := validate.Struct(todo)
	if errs != nil {
		var sliceErrs []string
		for _, err := range errs.(validator.ValidationErrors) {
			//翻译错误信息
			sliceErrs = append(sliceErrs, err.Translate(trans))
		}
		return errors.New(strings.Join(sliceErrs, ","))
	}
	err = dao.DB.Create(&todo).Error
	return
}

func GetAllTodo() (todoList []*Todo, err error) {
	if err = dao.DB.Find(&todoList).Error; err != nil {
		return nil, err
	}
	return
}

func GetATodo(id string) (todo *Todo, err error) {
	todo = new(Todo)
	if err = dao.DB.Debug().Where("id = ?", id).First(todo).Error; err != nil {
		return nil, err
	}
	return
}

func UpdateATodo(todo *Todo) (err error) {
	err = dao.DB.Save(todo).Error
	return
}

func DeleteATodo(id string) (err error) {
	tx := dao.DB.Debug().Where("id = ?", id).Delete(&Todo{})
	if tx.RowsAffected == 0 {
		return errors.New("删除内容不存在！")
	}
	return
}
