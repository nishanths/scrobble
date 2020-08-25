import { InsightsState } from "../types/insights"
import { insightTypes } from "../../components/u"
import { InsightsAction } from "../actions/insights"
import { copyMap } from "../../shared/util"

const defaultInsightsState = (() => {
	const m: InsightsState = new Map()
	insightTypes.forEach(it => {
		defaultInsightsState.set(it, { status: "initial" })
	})
	return m
})()

export function insightsReducer(state = defaultInsightsState, action: InsightsAction): InsightsState {
	const it = action.insightType
	switch (action.type) {
		case "INSIGHTS_START": {
			const s = copyMap(state)
			s.set(it, { status: "fetching" })
			return s
		}
		case "INSIGHTS_SUCCESS": {
			const s = copyMap(state)
			const v = { status: "done" as const, private: action.private, data: action.data }
			s.set(it, v)
			return s
		}
		case "INSIGHTS_FAIL": {
			const s = copyMap(state)
			s.set(it, { status: "error" })
			return s
		}
		default: {
			return state
		}
	}
}
