import React from "react"
import "../../scss/graph/graph.scss"

type GraphProps = {
	data: unknown
}

export class Graph extends React.Component<GraphProps> {
	constructor(props: GraphProps) {
		super(props)
	}

	componentDidMount() {
	}

	render() {
		return <div className="Graph">
		</div>
	}
}
