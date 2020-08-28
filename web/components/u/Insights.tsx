import React, { useEffect } from "react"
import { useSelector, useDispatch } from "react-redux"
import { State } from "../../redux/types/u"
import { Mode, modeFromControlValue, fullPath, InsightType } from "./shared"
import { RouteComponentProps } from "react-router-dom";
import { Header, SegmentedControl, Top } from "./top"
import { fetchInsights } from "../../redux/actions/insights"
import { Graph } from "../graph"
import { assertExhaustive } from "../../shared/util"
import { NProgress } from "../../shared/types"
import "../../scss/u/insights.scss"

type History = RouteComponentProps["history"]

type InsightsProps = {
	profileUsername: string
	signedIn: boolean
	private: boolean
	self: boolean
	insightType: InsightType
	history: History
	nProgress: NProgress
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
	private: privateProfile,
	self,
	signedIn,
	insightType,
	history,
	nProgress,
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

	let main: React.ReactNode

	switch (insight.status) {
		case "initial":
			// nothing to do
			break
		case "fetching":
			nProgress.start()
			break
		case "done":
			nProgress.done()
			if (insight.private === true) {
				main = <div className="info">(This user's data is private.)</div>
			} else {
				main = <div className="graph">
					<Graph
						data={insight.data}
						type={insightType}
					/>
				</div>
			}
			break
		case "error":
			nProgress.done()
			main = <div className="info">(Failed to fetch insights.)</div>
			break
		default:
			assertExhaustive(insight)
	}

	return <div className="Insights">
		{top}
		<div className="select-container">
			<select value={insightType} onChange={e => { history.push(fullPath(profileUsername, Mode.Insights, last.color, e.target.value as InsightType, undefined)) }}>
				{insightsOptionData.map(d => <option key={d.type} disabled={!!d.disabled} value={d.type}>{d.display}</option>)}
			</select>
		</div>
		<main>
			{main}
		</main>
	</div>
}
