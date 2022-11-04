package domain

import "strings"

const (
	SegmentNew = "Новостройка"
	SegmentMid = "Современное жилье"
	SegmentOld = "Старый жилой фонд"
)

const (
	WallBrick = "Кирпич"
	WallMono  = "Монолит"
	WallPanel = "Панель"
)

const (
	StateOff = "Без отделки"
	StateMun = "Муниципальный ремонт"
	StateNew = "Современная отделка"
)

const (
	Yes = "Да"
	No  = "Нет"
)

const Studio = "Студия"

var (
	LowerSegmentNew = strings.ToLower(SegmentNew)
	LowerSegmentMid = strings.ToLower(SegmentMid)
	LowerSegmentOld = strings.ToLower(SegmentOld)

	LowerWallBrick = strings.ToLower(WallBrick)
	LowerWallMono  = strings.ToLower(WallMono)
	LowerWallPanel = strings.ToLower(WallPanel)

	LowerStateOff = strings.ToLower(StateOff)
	LowerStateMun = strings.ToLower(StateMun)
	LowerStateNew = strings.ToLower(StateNew)

	LowerYes = strings.ToLower(Yes)
	LowerNo  = strings.ToLower(No)

	LowerStudio = strings.ToLower(Studio)
)
