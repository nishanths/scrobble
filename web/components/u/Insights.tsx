import React from "react"
import { useSelector } from "react-redux"
import { State } from "../../redux/types/u"
import { Mode, modeFromControlValue, fullPath, InsightType } from "./shared"
import { RouteComponentProps } from "react-router-dom";
import { Header, SegmentedControl, Top } from "./top"
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
}

const insightsOptionData: readonly InsightOption[] = [
	{ type: "artist-discovery", display: "Artist discovered" },
	{ type: "most-played-songs", display: "Most played songs" },
	{ type: "most-listened-artists", display: "Most listened artists" },
	{ type: "longest-songs", display: "Longest songs" },
]

export const Insights: React.FC<InsightsProps> = ({
	profileUsername,
	signedIn,
	insightType,
	history,
}) => {
	const last = useSelector((s: State) => s.last)

	const header = Header(profileUsername, signedIn, true)
	const segmentedControl = SegmentedControl(Mode.Insights, v => {
		const p = fullPath(profileUsername, modeFromControlValue(v), last.color, insightType, undefined)
		history.push(p)
	})
	const top = Top(header, segmentedControl, null, Mode.Insights)

	return <div className="Insights">
		{top}
		<div className="select-container">
			<select value={insightType} onChange={e => { history.push(fullPath(profileUsername, Mode.Insights, last.color, e.target.value as InsightType, undefined)) }}>
				{insightsOptionData.map(d => <option key={d.type} value={d.type}>{d.display}</option>)}
			</select>
		</div>
	</div>
}
