import React, { useEffect, useRef } from "react"
import { useSelector, useDispatch } from "react-redux"
import { State } from "../../redux/types/u"
import { Mode, modeFromControlValue, fullPath, InsightType } from "./shared"
import { RouteComponentProps } from "react-router-dom";
import { Header, SegmentedControl, Top } from "./top"
import { fetchInsights } from "../../redux/actions/insights"
import { Graph } from "../graph"
import "../../scss/u/insights.scss"

type History = RouteComponentProps["history"]

type InsightsProps = {
	profileUsername: string
	signedIn: boolean
	insightType: InsightType
	history: History
}

type InsightOption = {
	type: InsightType
	display: string
	disabled?: boolean
}

const insightsOptionData: readonly InsightOption[] = [
	{ type: "artist-discovery", display: "Artist discovery timeline", disabled: true },
	{ type: "most-played-songs", display: "Most played songs", disabled: false },
	{ type: "most-listened-artists", display: "Most listened to artists", disabled: true },
	{ type: "longest-songs", display: "Longest songs", disabled: false },
]

export const Insights: React.FC<InsightsProps> = ({
	profileUsername,
	signedIn,
	insightType,
	history,
}) => {
	const dispatch = useDispatch()
	const last = useSelector((s: State) => s.last)
	const insight = useSelector((s: State) => s.insights.get(insightType)!)

	useEffect(() => {
		if (insight.status === "initial" || insight.status === "error") {
			dispatch(fetchInsights(insightType, profileUsername))
		}
	}, [profileUsername, insightType, insight])

	const header = Header(profileUsername, signedIn, true)
	const segmentedControl = SegmentedControl(Mode.Insights, v => {
		const p = fullPath(profileUsername, modeFromControlValue(v), last.color, insightType, undefined)
		history.push(p)
	})
	const top = Top(header, segmentedControl, null, Mode.Insights)

	console.log(insight)

	return <div className="Insights">
		{top}
		<div className="select-container">
			<select value={insightType} onChange={e => { history.push(fullPath(profileUsername, Mode.Insights, last.color, e.target.value as InsightType, undefined)) }}>
				{insightsOptionData.map(d => <option key={d.type} disabled={!!d.disabled} value={d.type}>{d.display}</option>)}
			</select>
		</div>
		<main>
		</main>
	</div>
}
