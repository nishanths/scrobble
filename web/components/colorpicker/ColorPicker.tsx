import React, { useState, useEffect } from "react"
import { Swatch } from "./Swatch"
import { colors, Color } from "./color"
import "../../scss/color-picker.scss"

type ColorPickerProps = {
	afterSelect?: (c: Color) => void
}

export const ColorPicker: React.FC<ColorPickerProps> = ({ afterSelect }) => {
	const [selected, setSelected] = useState<Color|undefined>(undefined)

	return <div className="ColorPicker">
		{
			colors.map(c => {
				const swatch = <Swatch
					color={c}
					selected={selected === c}
					onSelect={() => { setSelected(c); afterSelect?.(c) }}
				/>
				return <div className="space">{swatch}</div>
			})
		}
	</div>
}
