import { InsightsState } from "../types/insights"
import { insightTypes } from "../../components/u"
import { InsightsAction } from "../actions/insights"

const defaultInsightsState = (() => {
	const m: InsightsState = new Map()
	insightTypes.forEach(it => {
		defaultInsightsState.set(it, { status: "initial" })
	})
	return m
})()

export function insightsReducer(state = defaultInsightsState, action: InsightsAction): InsightsState {
	return new Map()
}
