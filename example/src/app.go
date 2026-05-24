package main

import "gofront"

func main() {
	count := gofront.State(0)
	button := gofront.Query("#increment")
	button.OnClick(func() {
		count.Set(count.Get() + 1)
		gofront.Query("#counter").SetText(count.Get())
	})

	gofront.Wait()
}
