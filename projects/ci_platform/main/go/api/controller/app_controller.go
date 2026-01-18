package controller

type App struct {
	User      interface{ User } `validate:"required"`
	Frontend  Frontend          `validate:"required"`
	Lifecycle Lifecycle         `validate:"required"`
	Project   Project           `validate:"required"`
	Branch    Branch            `validate:"required"`
	Pipeline  Pipeline          `validate:"required"`
	Session   Session           `validate:"required"`
	JWT       JWT               `validate:"required"`
}
