package models

import (
	"errors"
	"gorm.io/gorm"
	"online-voice-channel/dao"
)

// Todo Model
type Todo struct {
	ID      uint           `json:"id"`
	Title   string         `json:"title"`
	Status  bool           `json:"status"`
	Deleted gorm.DeletedAt `json:"deleted"`
}

/*
	Todo这个Model的增删改查操作都放在这里
*/
// CreateATodo 创建todo
func CreateTodo(todo *Todo) (err error) {
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
