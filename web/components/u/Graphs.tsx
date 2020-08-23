import React from "react"
import { useSelector } from "react-redux"
import { State } from "../../redux/types/u"
import { Mode, modeFromControlValue, fullPath } from "./shared"
import { RouteComponentProps } from "react-router-dom";
import { Header, SegmentedControl, Top } from "./top"
import "../../scss/u/graphs.scss"

type History = RouteComponentProps["history"]

type GraphsProps = {
	profileUsername: string
	signedIn: boolean
	history: History
}

export const Graphs: React.FC<GraphsProps> = ({
	profileUsername,
	signedIn,
	history,
}) => {
	const last = useSelector((s: State) => s.last)

	const header = Header(profileUsername, signedIn, true)
	const segmentedControl = SegmentedControl(Mode.Insights, v => {
		const p = fullPath(profileUsername, modeFromControlValue(v), last.color, undefined)
		history.push(p)
	})
	const top = Top(header, segmentedControl, null, Mode.Insights)

	return <div className="Graphs">
		{top}
		<div className="select-container">
			<select>
				<option>Longest songs</option>
				<option>Most played songs</option>
				<option>Artist Discovery dates</option>
				<option>Most played artists</option>
			</select>
		</div>
	</div>
}
