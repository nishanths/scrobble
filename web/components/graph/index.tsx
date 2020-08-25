import React from "react"
import "../../scss/graph/graph.scss"

type GraphProps = {}

type GraphState = {
	data: unknown
}

export class Graph extends React.Component<GraphProps, GraphState> {
	constructor(props: GraphProps) {
		super(props)
	}

	async componentDidMount() {
		try {
			// TODO: move to redux
			const rsp = await fetch("/api/v1/data/most-played-songs")
			switch (rsp.status) {
				case 200:

					break
				default:
					console.error("bad status: %d", rsp.status)
					break
			}
		} catch (e) {
			console.error(e)
		}
	}
	render() {
		return <div className="Graph">
		</div>
	}
}
