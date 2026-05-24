package main

import "github.com/Alazar42/GoFront"

func main() {
	GoFront.Run(func() {
		count := GoFront.State(0)
		button := GoFront.Query("#increment")
		button.OnClick(func() {
			count.Set(count.Get() + 1)
			GoFront.Query("#counter").SetText(count.Get())
		})
	})
}
