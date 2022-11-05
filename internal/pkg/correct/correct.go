package correct

import (
	"leaders_apartments/internal/pkg/domain"
)

func Do(target *domain.Row, analogs []*domain.Row) []*domain.CorrectCoefs {
	corrects := make([]*domain.CorrectCoefs, len(analogs))
	for i := range analogs {
		corrects[i] = new(domain.CorrectCoefs)
		Sale(analogs[i], corrects[i])
		Total(target, analogs[i], corrects[i])
		Metro(target, analogs[i], corrects[i])
		Floor(target, analogs[i], corrects[i])
		Kitchen(target, analogs[i], corrects[i])
		Balcony(target, analogs[i], corrects[i])
		State(target, analogs[i], corrects[i])
	}
	return corrects
}

func Sale(analog *domain.Row, coef *domain.CorrectCoefs) {
	coef.Sale = -0.045
	analog.AvgCost *= (1 + coef.Sale)
}

func Floor(target, analog *domain.Row, coef *domain.CorrectCoefs) {
	switch {
	case target.CFloor == 1 && analog.CFloor > 1 && analog.CFloor < analog.Floors:
		coef.Floor = -0.07
	case target.CFloor == 1 && analog.CFloor > 1 && analog.CFloor == analog.Floors:
		coef.Floor = -0.031
	case target.CFloor < target.Floors && analog.CFloor == 1:
		coef.Floor = 0.075
	case target.CFloor < target.Floors && analog.CFloor == analog.Floors:
		coef.Floor = 0.042
	case target.CFloor == target.Floors && analog.CFloor == 1 && analog.CFloor < analog.Floors:
		coef.Floor = 0.032
	case target.CFloor == target.Floors && analog.CFloor > 1 && analog.CFloor < analog.Floors:
		coef.Floor = -0.04
	default:
		coef.Floor = 0
	}
	analog.AvgCost *= (1 + coef.Floor)
}

func Total(target, analog *domain.Row, coef *domain.CorrectCoefs) {
	switch {
	case target.Total < 30 && analog.Total >= 30 && analog.Total < 50:
		coef.Total = 0.06
	case target.Total < 30 && analog.Total >= 50 && analog.Total < 65:
		coef.Total = 0.14
	case target.Total < 30 && analog.Total >= 65 && analog.Total < 90:
		coef.Total, analog.Good = 0.21, -1
	case target.Total < 30 && analog.Total >= 90 && analog.Total < 120:
		coef.Total, analog.Good = 0.28, -1
	case target.Total < 30 && analog.Total >= 120:
		coef.Total, analog.Good = 0.31, -1
	case target.Total >= 30 && target.Total < 50 && analog.Total < 30:
		coef.Total = -0.06
	case target.Total >= 30 && target.Total < 50 && analog.Total >= 50 && analog.Total < 65:
		coef.Total = 0.07
	case target.Total >= 30 && target.Total < 50 && analog.Total >= 65 && analog.Total < 90:
		coef.Total = 0.14
	case target.Total >= 30 && target.Total < 50 && analog.Total >= 90 && analog.Total < 120:
		coef.Total, analog.Good = 0.21, -1
	case target.Total >= 30 && target.Total < 50 && analog.Total >= 120:
		coef.Total, analog.Good = 0.24, -1
	case target.Total >= 50 && target.Total < 65 && analog.Total < 30:
		coef.Total = -0.12
	case target.Total >= 50 && target.Total < 65 && analog.Total >= 30 && analog.Total < 50:
		coef.Total = -0.07
	case target.Total >= 50 && target.Total < 65 && analog.Total >= 65 && analog.Total < 90:
		coef.Total = 0.06
	case target.Total >= 50 && target.Total < 65 && analog.Total >= 90 && analog.Total < 120:
		coef.Total = 0.13
	case target.Total >= 50 && target.Total < 65 && analog.Total >= 120:
		coef.Total, analog.Good = 0.16, -1
	case target.Total >= 65 && target.Total < 90 && analog.Total < 30:
		coef.Total, analog.Good = -0.17, -1
	case target.Total >= 65 && target.Total < 90 && analog.Total >= 30 && analog.Total < 50:
		coef.Total = -0.12
	case target.Total >= 65 && target.Total < 90 && analog.Total >= 50 && analog.Total < 65:
		coef.Total = -0.06
	case target.Total >= 65 && target.Total < 90 && analog.Total >= 90 && analog.Total < 120:
		coef.Total = 0.06
	case target.Total >= 65 && target.Total < 90 && analog.Total >= 120:
		coef.Total = 0.09
	case target.Total >= 90 && target.Total < 120 && analog.Total < 30:
		coef.Total, analog.Good = -0.22, -1
	case target.Total >= 90 && target.Total < 120 && analog.Total >= 30 && analog.Total < 50:
		coef.Total, analog.Good = -0.17, -1
	case target.Total >= 90 && target.Total < 120 && analog.Total >= 50 && analog.Total < 65:
		coef.Total = -0.11
	case target.Total >= 90 && target.Total < 120 && analog.Total >= 65 && analog.Total < 90:
		coef.Total = -0.06
	case target.Total >= 90 && target.Total < 120 && analog.Total >= 120:
		coef.Total = 0.03
	case target.Total > 120 && analog.Total < 30:
		coef.Total, analog.Good = -0.24, -1
	case target.Total > 120 && analog.Total >= 30 && analog.Total < 50:
		coef.Total, analog.Good = -0.19, -1
	case target.Total > 120 && analog.Total >= 50 && analog.Total < 65:
		coef.Total, analog.Good = -0.13, -1
	case target.Total > 120 && analog.Total >= 65 && analog.Total < 90:
		coef.Total = -0.08
	case target.Total > 120 && analog.Total >= 90 && analog.Total < 120:
		coef.Total = -0.03
	default:
		coef.Total = 0
	}
	analog.AvgCost *= (1 + coef.Total)
}

func Kitchen(target, analog *domain.Row, coef *domain.CorrectCoefs) {
	switch {
	case target.Kitchen < 7 && analog.Kitchen >= 7 && analog.Kitchen < 10:
		coef.Kitchen = -0.029
	case target.Kitchen < 7 && analog.Kitchen >= 10:
		coef.Kitchen = -0.083
	case target.Kitchen >= 7 && target.Kitchen < 10 && analog.Kitchen < 7:
		coef.Kitchen = 0.03
	case target.Kitchen >= 7 && target.Kitchen < 10 && analog.Kitchen >= 10:
		coef.Kitchen = -0.055
	case target.Kitchen >= 10 && analog.Kitchen < 7:
		coef.Kitchen = 0.09
	case target.Kitchen >= 10 && analog.Kitchen >= 7 && analog.Kitchen < 10:
		coef.Kitchen = 0.058
	default:
		coef.Kitchen = 0
	}
	analog.AvgCost *= (1 + coef.Kitchen)
}

func Balcony(target, analog *domain.Row, coef *domain.CorrectCoefs) {
	switch {
	case target.Balcony == domain.Yes && analog.Balcony == domain.No:
		coef.Balcony = 0.053
	case target.Balcony == domain.No && analog.Balcony == domain.Yes:
		coef.Balcony = -0.05
	default:
		coef.Balcony = 0
	}
	analog.AvgCost *= (1 + coef.Balcony)
}

func Metro(target, analog *domain.Row, coef *domain.CorrectCoefs) {
	switch {
	case target.Metro < 5 && analog.Metro >= 5 && analog.Metro < 10:
		coef.Metro = 0.07
	case target.Metro < 5 && analog.Metro >= 10 && analog.Metro < 15:
		coef.Metro = 0.12
	case target.Metro < 5 && analog.Metro >= 15 && analog.Metro < 30:
		coef.Metro, analog.Good = 0.17, -1
	case target.Metro < 5 && analog.Metro >= 30 && analog.Metro < 60:
		coef.Metro, analog.Good = 0.24, -1
	case target.Metro < 5 && analog.Metro >= 60:
		coef.Metro, analog.Good = 0.29, -1
	case target.Metro >= 5 && target.Metro < 10 && analog.Metro < 5:
		coef.Metro = -0.07
	case target.Metro >= 5 && target.Metro < 10 && analog.Metro >= 10 && analog.Metro < 15:
		coef.Metro = 0.04
	case target.Metro >= 5 && target.Metro < 10 && analog.Metro >= 15 && analog.Metro < 30:
		coef.Metro = 0.09
	case target.Metro >= 5 && target.Metro < 10 && analog.Metro >= 30 && analog.Metro < 60:
		coef.Metro, analog.Good = 0.15, -1
	case target.Metro >= 5 && target.Metro < 10 && analog.Metro >= 60:
		coef.Metro, analog.Good = 0.2, -1
	case target.Metro >= 10 && target.Metro < 15 && analog.Metro < 5:
		coef.Metro = -0.11
	case target.Metro >= 10 && target.Metro < 15 && analog.Metro >= 5 && analog.Metro < 10:
		coef.Metro = -0.04
	case target.Metro >= 10 && target.Metro < 15 && analog.Metro >= 15 && analog.Metro < 30:
		coef.Metro = 0.05
	case target.Metro >= 10 && target.Metro < 15 && analog.Metro >= 30 && analog.Metro < 60:
		coef.Metro = 0.11
	case target.Metro >= 10 && target.Metro < 15 && analog.Metro >= 60:
		coef.Metro, analog.Good = 0.15, -1
	case target.Metro >= 15 && target.Metro < 30 && analog.Metro < 5:
		coef.Metro, analog.Good = -0.15, -1
	case target.Metro >= 15 && target.Metro < 30 && analog.Metro >= 5 && analog.Metro < 10:
		coef.Metro = -0.08
	case target.Metro >= 15 && target.Metro < 30 && analog.Metro >= 10 && analog.Metro < 15:
		coef.Metro = -0.05
	case target.Metro >= 15 && target.Metro < 30 && analog.Metro >= 30 && analog.Metro < 60:
		coef.Metro = 0.06
	case target.Metro >= 15 && target.Metro < 30 && analog.Metro >= 60:
		coef.Metro = 0.1
	case target.Metro >= 30 && target.Metro < 60 && analog.Metro < 5:
		coef.Metro, analog.Good = -0.19, -1
	case target.Metro >= 30 && target.Metro < 60 && analog.Metro >= 5 && analog.Metro < 10:
		coef.Metro, analog.Good = -0.13, -1
	case target.Metro >= 30 && target.Metro < 60 && analog.Metro >= 10 && analog.Metro < 15:
		coef.Metro = -0.1
	case target.Metro >= 30 && target.Metro < 60 && analog.Metro >= 15 && analog.Metro < 30:
		coef.Metro = -0.06
	case target.Metro >= 30 && target.Metro < 60 && analog.Metro >= 60:
		coef.Metro = 0.04
	case target.Metro >= 60 && analog.Metro < 5:
		coef.Metro, analog.Good = -0.22, -1
	case target.Metro >= 60 && analog.Metro >= 5 && analog.Metro < 10:
		coef.Metro, analog.Good = -0.17, -1
	case target.Metro >= 60 && analog.Metro >= 10 && analog.Metro < 15:
		coef.Metro, analog.Good = -0.13, -1
	case target.Metro >= 60 && analog.Metro >= 15 && analog.Metro < 30:
		coef.Metro = -0.09
	case target.Metro >= 60 && analog.Metro >= 30 && analog.Metro < 60:
		coef.Metro = -0.04
	default:
		coef.Metro = 0
	}
	analog.AvgCost *= (1 + coef.Metro)
}

func State(target, analog *domain.Row, coef *domain.CorrectCoefs) {
	switch {
	case target.State == domain.StateOff && analog.State == domain.StateMun:
		coef.State = -13400
	case target.State == domain.StateOff && analog.State == domain.StateNew:
		coef.State = -20100
	case target.State == domain.StateMun && analog.State == domain.StateOff:
		coef.State = 13400
	case target.State == domain.StateMun && analog.State == domain.StateNew:
		coef.State = -6700
	case target.State == domain.StateNew && analog.State == domain.StateOff:
		coef.State = 20100
	case target.State == domain.StateNew && analog.State == domain.StateMun:
		coef.State = 6700
	default:
		coef.State = 0
	}
	analog.AvgCost += coef.State
}
