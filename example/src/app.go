package main

import "gofront"

func main() {
	gofront.Run(func() {
		count := gofront.State(0)
		button := gofront.Query("#increment")
		button.OnClick(func() {
			count.Set(count.Get() + 1)
			gofront.Query("#counter").SetText(count.Get())
		})
	})
}
