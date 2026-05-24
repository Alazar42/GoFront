package components

import "gofront"

func Counter() gofront.Component {
	return gofront.Div(
		gofront.Text("Hello"),
	)
}
