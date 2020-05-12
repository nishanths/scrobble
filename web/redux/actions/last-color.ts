import { Color } from "../../components/colorpicker"
import { PartialState } from "../types/u"

export type LastColorAction = ReturnType<typeof setLastColor>

export const setLastColor = (c: Color) => {
  return {
    type: "SET_LAST_COLOR" as const,
    color: c,
  }
}
