package xslx

import "github.com/plandem/xlsx"

func Parse() error {
	xl := xlsx.New()
	defer xl.Close()

	return nil
}
