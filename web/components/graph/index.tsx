import React from "react"
import "../../scss/graph/graph.scss"
import { InsightType } from "../u"

type GraphProps = {
	data: unknown
	type: InsightType
}

export class Graph extends React.Component<GraphProps> {
	constructor(props: GraphProps) {
		super(props)
	}

	componentDidMount() {
	}

	render() {
		console.log(this.props.data)

		return <div className="Graph">
		</div>
	}
}
