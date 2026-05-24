package pages

import "github.com/Alazar42/GoFront"

func HomePage() GoFront.Component {
	return GoFront.Div(
		GoFront.Text("Home"),
	)
}
