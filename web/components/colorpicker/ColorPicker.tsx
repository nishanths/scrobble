import React, { useState, useEffect } from "react"
import { Swatch } from "./Swatch"
import { colors, Color } from "./color"
import { capitalize } from "../../shared/util"
import "../../scss/color-picker.scss"

type ColorPickerProps = {
	afterSelect?: (c: Color) => void
	prompt?: string;
}

export const ColorPicker: React.FC<ColorPickerProps> = ({ afterSelect, prompt }) => {
	const [selected, setSelected] = useState<Color|undefined>(undefined)

	const elems = colors.map(c => {
		return <div className="elem">
			<Swatch color={c} selected={selected === c} onSelect={() => { setSelected(c); afterSelect?.(c) }}/>
		</div>
	})

	return <div className="ColorPicker">
		<div className="elems">{elems}</div>
		<div className="label">{selected === undefined ? prompt : capitalize(selected)}</div>
	</div>
}
