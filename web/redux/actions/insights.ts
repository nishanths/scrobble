import { ThunkAction, ThunkDispatch } from "redux-thunk"
import { PartialState } from "../types/u"
import { InsightType } from "../../components/u"
import { assertExhaustive } from "../../shared/util"

export type InsightsAction =
	ReturnType<typeof insightsStart> |
	ReturnType<typeof insightsSuccess> |
	ReturnType<typeof insightsFail>

type InsightsThunkDispatch = ThunkDispatch<PartialState, undefined, InsightsAction>
type InsightsThunkResult<R> = ThunkAction<R, PartialState, undefined, InsightsAction>

export function insightsStart(it: InsightType, username: string) {
	return {
		type: "INSIGHTS_START" as const,
		insightType: it,
		username,
	}
}
export function insightsSuccess(it: InsightType, username: string, data: unknown, priv: boolean) {
	return {
		type: "INSIGHTS_SUCCESS" as const,
		data,
		private: priv,
		username,
	}
}

export function insightsFail(it: InsightType, username: string, err: any) {
	return {
		type: "INSIGHTS_FAIL" as const,
		username,
		err,
	}
}

export const fetchInsights = (it: InsightType, username: string): InsightsThunkResult<void> => {
	return async (dispatch) => {
		dispatch(insightsStart(it, username))

		try {
			const result = await _fetchInsights(it, username)
			dispatch(insightsSuccess(it, username, result.data, result.private))
		} catch (e) {
			dispatch(insightsFail(e, it, username))
		}
	}
}

type FetchInsightsResult = {
	data: unknown
	private: boolean
	err: any | null
}

function insightTypeToAPIPath(it: InsightType): string {
	switch (it) {
		case "most-listened-artists":
			return "/artists/play-count"
		case "most-played-songs":
			return "/songs/play-count"
		case "longest-songs":
			return "/songs/length"
		case "artist-discovery":
			return "/artists/added"
		default:
			assertExhaustive(it)
	}
}

const _fetchInsights = async (it: InsightType, username: string): Promise<FetchInsightsResult> => {
	const url = "/api/v1/data" + insightTypeToAPIPath(it) + "?username=" + username

	const r = await fetch(url)
	switch (r.status) {
		case 200: {
			const data = await r.json()
			return { data, private: false, err: null }
		}
		case 403:
			return { data: null, private: true, err: null }
		// TODO: if we had the ability to display toast notifications, we could show
		// "please sign in again" for 401 status
		default:
			throw "bad status: " + r.status
	}
}
