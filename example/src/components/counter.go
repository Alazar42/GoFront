package components

import "github.com/Alazar42/GoFront"

func Counter() GoFront.Component {
	return GoFront.Div(
		GoFront.Text("Hello"),
	)
}
