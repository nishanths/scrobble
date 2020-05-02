import React from "react"
import { Color } from "./color"
import "../../scss/swatch.scss"

type SwatchProps = {
	color: Color
	selected: boolean
}

export const Swatch: React.SFC<SwatchProps> = ({ color, selected }) => {
	const className = selected ?
		`Swatch color-${color} selected` :
		`Swatch color-${color}`

	return <div className={className}>
		<div className="selection"></div>
	</div>
}
