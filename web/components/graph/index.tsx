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

		const data = [
			{name: "E", value: 0.12702},
			{name: "T", value: 0.09056},
			{name: "A", value: 0.08167},
			{name: "O", value: 0.07507},
			{name: "I", value: 0.06966},
			{name: "N", value: 0.06749},
			{name: "S", value: 0.06327},
			{name: "H", value: 0.06094},
		]


		const margin = {top: 20, right: 0, bottom: 30, left: 40}
		const width = 500
		const height = 500

		const x = d3.scaleBand()
			.domain(data.map(d => d.name))
			.range([margin.left, width - margin.right])
			.padding(0.1)

		const y = d3.scaleLinear()
		    .domain([0, d3.max(data, d => d.value)!]).nice()
		    .range([height - margin.bottom, margin.top])

		const yAxis = (g: any) => g
		    .attr("transform", `translate(${margin.left},0)`)
		    .call(d3.axisLeft(y))
		    .call((g: any) => g.select(".domain").remove())

		const xAxis = (g: any) => g
		    .attr("transform", `translate(0,${height - margin.bottom})`)
		    .call(d3.axisBottom(x).tickSizeOuter(0))

		function zoom(svg: any) {
		  const extent: [[number, number], [number, number]] = [[margin.left, margin.top], [width - margin.right, height - margin.top]];

		  svg.call(d3.zoom()
		      .scaleExtent([1, 8])
		      .translateExtent(extent)
		      .extent(extent)
		      .on("zoom", zoomed));

		  function zoomed() {
		  	const ns: [number, number] = [margin.left, width - margin.right]
		  	const mapped = ns.map(d => d3.event.transform.applyX(d)) as [number, number]
		  	x.range(mapped);
		  	svg.selectAll(".bars rect").attr("x", (d: any) => x(d.name)).attr("width", x.bandwidth());
		  	svg.selectAll(".x-axis").call(xAxis);
		  }
		}

		const svg = d3.select(".Graph").append("svg")
			.attr("viewBox", "0, 0, 400, 400")
			.attr("preserveAspectRatio", "xMidYMid meet")
			.call(zoom)

		svg.append("g")
			.attr("class", "bars")
			.attr("fill", "steelblue")
			.selectAll("rect")
			.data(data)
			.join("rect")
			.attr("x", (d: any) => x(d.name)!)
			.attr("y", (d: any) => y(d.value))
			.attr("height", (d: any) => y(0) - y(d.value))
			.attr("width", x.bandwidth());

			svg.append("g")
			.attr("class", "x-axis")
			.call(xAxis);

			svg.append("g")
			.attr("class", "y-axis")
			.call(yAxis);
	}

	render() {
		return null
	}
}
