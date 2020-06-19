import React from "react"
import { useSelector } from "react-redux"
import { State } from "../../redux/types/u"
import { Mode, modeFromControlValue, fullPath } from "./shared"
import { RouteComponentProps } from "react-router-dom";
import { Header, SegmentedControl, Top } from "./top"

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
	const segmentedControl = SegmentedControl(Mode.Graphs, v => {
		const p = fullPath(profileUsername, modeFromControlValue(v), last.color, undefined)
		history.push(p)
	})
	const top = Top(header, segmentedControl, null, Mode.Graphs)

	return <>
		{top}
		<div className="info">Graphs coming soon!</div>
	</>
}
