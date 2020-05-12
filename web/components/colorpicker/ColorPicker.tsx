import React, { useState, useEffect } from "react"
import { Swatch } from "./Swatch"
import { colors, Color } from "./color"
import { capitalize } from "../../shared/util"
import "../../scss/color-picker.scss"

type ColorPickerProps = {
  afterSelect?: (c: Color) => void
  prompt?: string
  initialSelection?: Color // if unset no color initially selected
}

export const ColorPicker: React.FC<ColorPickerProps> = ({ afterSelect, prompt, initialSelection }) => {
  const [selected, setSelected] = useState<Color | undefined>(initialSelection)
  useEffect(() => {
    setSelected(initialSelection)
  }, [initialSelection])

  const label = () => {
    if (selected !== undefined) {
      return capitalize(selected)
    }
    return prompt
  }

  const elems = colors.map(c => {
    return <div key={c} className="elem" title={capitalize(c)}>
      <Swatch
        color={c}
        selected={selected === c}
        onSelect={() => { setSelected(c); afterSelect?.(c) }}
      />
    </div>
  })

  return <div className="ColorPicker">
    <div className="elems">{elems}</div>
    <div className="label">{label()}</div>
  </div>
}
