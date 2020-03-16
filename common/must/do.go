package must

import "fmt"

func Do(err error) {
	if err != nil {
		fmt.Errorf("%+v", err.Error())
	}
}

func DoF(f func() error) {
	Do(f())
}
