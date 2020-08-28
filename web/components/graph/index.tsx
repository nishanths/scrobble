import React from "react"
import "../../scss/graph/graph.scss"
import { InsightType } from "../u"
import { assertExhaustive } from "../../shared/util"
import { SongsDataResponse } from "../../shared/types"
import { colors } from "../../shared/const"
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
				content = <MostPlayedSongs data={this.props.data as SongsDataResponse} />
				break
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
	data: SongsDataResponse
}

const maxGraphItems = 125

export class MostPlayedSongs extends React.Component<MostPlayedSongsProps> {
	componentDidMount() {
		this.draw()
	}

	componentDidUpdate(newProps: MostPlayedSongsProps) {
		if (this.props.data !== newProps.data) {
			this.draw()
		}
	}

	componentWillUnmount() {
		// svg.remove()
	}

	private draw() {
		// XXX: d3.select() calls for svg may need to be more specific if there are ever
		// multiple graph elements

		// https://observablehq.com/@d3/zoomable-bar-chart

		if (this.props.data.length === 0) {
			return
		}

		const data = this.props.data.slice(0, maxGraphItems)

		const margin = { top: 20, right: 0, bottom: 30, left: 40 }
		const width = 1100
		const height = 800

		const x = d3.scaleBand()
			.domain(data.map(d => d.ident))
			.range([margin.left, width - margin.right])
			.padding(0.1)

		const y = d3.scaleLinear()
			.domain([0, d3.max(data, d => { return d.playCount })!]).nice()
			.range([height - margin.bottom, margin.top])

		const yAxis = (g: any) => g
			.attr("transform", `translate(${margin.left},0)`)
			.call(d3.axisLeft(y))
			.call((g: any) => g.select(".domain").remove())

		const xAxis = (g: any) => g
			.attr("transform", `translate(0,${height - margin.bottom})`)
			.call(d3.axisBottom(x).tickSizeOuter(0).tickFormat(() => "").tickSize(0))

		function zoom(svg: any) {
			const extent: [[number, number], [number, number]] = [[margin.left, margin.top], [width - margin.right, height - margin.top]];

			svg.call(d3.zoom()
				.scaleExtent([1, 4])
				.translateExtent(extent)
				.extent(extent)
				.on("zoom", zoomed));

			function zoomed() {
				const ns: [number, number] = [margin.left, width - margin.right]
				const mapped = ns.map(d => d3.event.transform.applyX(d)) as [number, number]
				x.range(mapped);
				svg.selectAll(".bars rect").attr("x", (d: any) => x(d.ident)).attr("width", x.bandwidth());
				svg.selectAll(".x-axis").call(xAxis);
			}
		}

		const svg = d3.select(".graph-content").append("svg")
			.attr("viewBox", "0 0 " + width + " " + height)
			.attr("preserveAspectRatio", "xMidYMid meet")
			.call(zoom)

		const tooltip = d3.select(".graph-content").append("div")
			.attr("class", "tooltip")
			.style("opacity", 0);

		svg.append("g")
			.attr("class", "bars")
			.attr("fill", colors.greenAlpha)
			.selectAll("rect")
			.data(data)
			.join("rect")
			.attr("x", (d) => x(d.ident)!)
			.attr("y", (d) => y(d.playCount))
			.attr("height", (d) => y(0) - y(d.playCount))
			.attr("width", x.bandwidth())
			.on("mouseover", function(d) {
				d3.select(this).transition()
					.duration(0)
					.style("fill", "rgb(0,133,210)")

				tooltip.transition()
					.duration(200)
					.style("opacity", .9);
				tooltip.html(d.title + "<br/>" + d.playCount)
					.style("left", (d3.event.pageX) + 20 + "px")
					.style("top", (d3.event.pageY - 80) + "px");
			})
			.on("mouseout", function() {
				d3.select(this).transition()
					.duration(0)
					.style("fill", colors.greenAlpha)

				tooltip.transition()
					.duration(500)
					.style("opacity", 0);
			});

		svg.append("g")
			.attr("class", "x-axis")
			.call(xAxis);

		svg.append("g")
			.attr("class", "y-axis")
			.call(yAxis);

		svg.append("text")
			.attr("class", "x label")
			.attr("text-anchor", "end")
			.attr("x", width)
			.attr("y", height - 6)
			.text("Songs");

		svg.append("text")
			.attr("class", "y label")
			.attr("text-anchor", "end")
			.attr("x", 120)
			.attr("y", 30)
			.attr("dy", ".75em")
			.text("↑ Times played");
	}

	render() {
		if (this.props.data.length === 0) {
			return <div className="info">Not enough data yet. Scrobble away!</div>
		}

		return <div className="MostPlayedSongs">
			<div className="graph-content"></div>
			<div className="instruction">You can hover over, zoom (scroll), and pan this graph.</div>
			<div className="table">
				<table>
					<thead>
						<tr>
							<th>#</th>
							<th>Title</th>
							<th className="artist">Artist</th>
							<th className="album">Album</th>
							<th>Play count</th>
						</tr>
					</thead>
					<tbody>
						{this.props.data.slice(0, maxGraphItems).map((d, i) => {
							return <tr key={d.ident}>
								<td>{i + 1}</td>
								<td>{d.title}</td>
								<td className="artist">{d.artistName}</td>
								<td className="album">{d.albumTitle}</td>
								<td>{d.playCount}</td>
							</tr>
						})}
					</tbody>
				</table>
			</div>
		</div>
	}
}

// tooltip

// copy to other two graphs
