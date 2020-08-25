import { InsightType } from "../../components/u"

export type InsightsState = Map<InsightType, InsightState>

export type InsightState =
	| { status: "initial" }
	| { status: "fetching" }
	| { status: "error" }
	| { status: "done", data: unknown }
	| { status: "done", private: true }
