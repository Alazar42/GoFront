package pages

import "gofront"

func HomePage() gofront.Component {
	return gofront.Div(
		gofront.Text("Home"),
	)
}
