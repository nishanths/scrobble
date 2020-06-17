import React from "react"
import { Color } from "./color"
import "../../scss/colorpicker/swatch.scss"

type SwatchProps = {
    color: Color
    selected: boolean
    onSelect?: () => void
    onEnter?: () => void
    onExit?: () => void
}

export const Swatch: React.SFC<SwatchProps> = ({ color, selected, onSelect, onEnter, onExit }) => {
    const className = selected ?
        `Swatch color-${color} selected` :
        `Swatch color-${color}`

    return <div className={className} onClick={() => onSelect?.()} onMouseEnter={() => onEnter?.()} onMouseLeave={() => onExit?.()}>
        <div className="selection"></div>
    </div>
}
