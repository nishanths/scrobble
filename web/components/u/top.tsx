import React from "react"
import { Header as HeaderComponent } from "./Header"
import { ColorPicker as ColorPickerComponent, Color } from "../colorpicker"
import { SegmentedControl as SegmentedControlComponent } from "../SegmentedControl"
import { controlValues, controlValueForMode, Mode, ControlValue } from "./shared"

export const Header = (profileUsername: string, signedIn: boolean, showNav: boolean) =>
	<HeaderComponent username={profileUsername} signedIn={signedIn} showNav={showNav} />

export const ColorPicker = (color?: Color, onColorChange?: (c: Color) => void) =>
	<div className="colorPicker">
		<ColorPickerComponent initialSelection={color} prompt="Pick a color to see scrobbled artwork of that color." afterSelect={onColorChange} />
	</div>

export const SegmentedControl = (mode: Mode, onControlChange: (v: ControlValue) => void) =>
	<div className="control">
		<SegmentedControlComponent
			afterChange={onControlChange}
			values={controlValues}
			initialValue={controlValueForMode(mode)}
			newBadges={new Set(["Insights" as const])}
		/>
	</div>

export const Top = (
	header: JSX.Element,
	segmentedControl: JSX.Element,
	colorPicker: JSX.Element | null, // necessary only when mode === Mode.Color
	mode: Mode,
) =>
	<>
		{header}
		{segmentedControl}
		{mode === Mode.Color && colorPicker}
	</>
