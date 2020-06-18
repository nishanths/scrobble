import { Color } from "../../components/colorpicker"
import { Mode } from "../../components/u/shared"

export type LastState = {
	color: Color | undefined
	scrobblesEndIdx: number | undefined
	scrobblesScrollY: number | undefined
	search: string
}
