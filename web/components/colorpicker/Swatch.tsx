import React from "react"
import { Color } from "./color"
import "../../scss/swatch.scss"

type SwatchProps = {
	color: Color
	selected: boolean
	onSelect?: () => void
}

export const Swatch: React.SFC<SwatchProps> = ({ color, selected, onSelect }) => {
	const className = selected ?
		`Swatch color-${color} selected` :
		`Swatch color-${color}`

	return <div className={className} onClick={() => onSelect?.()}>
		<div className="selection"></div>
	</div>
}
