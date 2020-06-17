import React from "react"
import { Header as HeaderComponent } from "./Header"
import { ColorPicker as ColorPickerComponent, Color } from "../colorpicker"
import { SegmentedControl } from "../SegmentedControl"
import { controlValues, controlValueForMode, Mode, ControlValue } from "./shared"

export const Header = (profileUsername: string, signedIn: boolean, showNav: boolean) =>
    <HeaderComponent username={profileUsername} signedIn={signedIn} showNav={showNav} />

export const ColorPicker = (color?: Color, onColorChange?: (c: Color) => void) =>
    <div className="colorPicker">
        <ColorPickerComponent initialSelection={color} prompt="Pick a color to see scrobbled artwork of that color." afterSelect={onColorChange} />
    </div>

export const Top = (header: JSX.Element, colorPicker: JSX.Element, mode: Mode, onControlChange: (v: ControlValue) => void) =>
    <>
        {header}
        <div className="control">
            <SegmentedControl
                afterChange={onControlChange}
                values={controlValues}
                initialValue={controlValueForMode(mode)}
            />
        </div>
        {mode === Mode.Color && colorPicker}
    </>
