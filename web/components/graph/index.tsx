import React from "react"
import "../../scss/graph/graph.scss"
import { InsightType } from "../u"
import { assertExhaustive } from "../../shared/util"
import * as d3 from "d3"

type GraphProps = {
	data: unknown
	type: InsightType
}

export class Graph extends React.Component<GraphProps> {
	constructor(props: GraphProps) {
		super(props)
	}

	render() {
		let content: React.ReactNode

		switch (this.props.type) {
			case "most-played-songs":
				content = <MostPlayedSongs data={this.props.data} />
			case "most-listened-artists":
				break
			case "artist-discovery":
				break
			case "longest-songs":
				break
			default:
				assertExhaustive(this.props.type)
		}

		return <div className="Graph">
			{content}
		</div>
	}
}

type MostPlayedSongsProps = {
	data: unknown // TODO: make this specific
}

export class MostPlayedSongs extends React.Component<MostPlayedSongsProps> {
	componentDidMount() {
		this.draw()
	}

	componentDidUpdate(newProps: GraphProps) {
		if (this.props.data !== newProps.data) {
			this.draw()
		}
	}

	componentWillUnmount() {
		// svg.remove()
	}

	private draw() {
		// XXX: selection may need to be more specific if there are ever
		// multiple graph elements

		const svg = d3.select(".Graph").append("svg")
			.attr("width", 400)
			.attr("height", 400)

		const data = [12, 5, 6, 6, 9, 10];

		svg.selectAll("rect")
			.data(data)
			.enter()
			.append("rect")
			.attr("x", (d, i) => i * 70)
			.attr("y", (d, i) => 400 - 10 * d)
			.attr("width", 65)
			.attr("height", (d, i) => d * 10)
			.attr("fill", "green")
	}

	render() {
		return null
	}
}
