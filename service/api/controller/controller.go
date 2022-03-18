package controller

type Controller struct {
	Block *BlockController
}

func NewController() *Controller {
	return &Controller{
		Block: NewBlockController(),
	}
}
