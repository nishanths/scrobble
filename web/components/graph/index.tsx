import React from "react"
import "../../scss/graph/graph.scss"
import { InsightType } from "../u"
import { assertExhaustive, pluralize } from "../../shared/util"
import { secondsToHms  } from "../../shared/time"
import { SongsDataResponse, Song, ArtistPlayCountDataResponse, ArtistPlayCountDatum } from "../../shared/types"
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
				const data = this.props.data as ArtistPlayCountDataResponse
				data.sort((a, b) => {
					if (a.totalPlayTime === b.totalPlayTime) {
						return a.artistName.localeCompare(b.artistName)
					}
					if (a.totalPlayTime > b.totalPlayTime) {
						return -1
					}
					return 1
				})
				content = <MostListenedArtists data={data} />
				break
			case "artist-discovery":
				break
			case "longest-songs":
				content = <LongestSongs data={this.props.data as SongsDataResponse} />
				break
			default:
				assertExhaustive(this.props.type)
		}

		return <div className="Graph">
			{content}
		</div>
	}
}


const graphMargin = { top: 20, right: 0, bottom: 30, left: 100 }
const graphHeight = 700
const graphWidth = 1200

type MostPlayedSongsProps = {
	data: SongsDataResponse
}

export class MostPlayedSongs extends React.Component<MostPlayedSongsProps> {
	private svg: any
	private tooltip: any

	private readonly maxDataItems = 125

	componentDidMount() {
		this.draw()
	}

	componentDidUpdate(newProps: MostPlayedSongsProps) {
		if (this.props.data !== newProps.data) {
			this.draw()
		}
	}

	componentWillUnmount() {
		this.svg?.remove()
		this.tooltip?.remove()
	}

	private tooltipHTML(d: Song): string {
		return `${d.artistName} – ${d.title}<br/>${d.playCount.toLocaleString()} ${pluralize("time", d.playCount)}`
	}

	private draw() {
		// https://observablehq.com/@d3/zoomable-bar-chart

		if (this.props.data.length === 0) {
			return
		}

		const data = this.props.data.slice(0, this.maxDataItems)

		const margin = graphMargin
		const width = graphWidth
		const height = graphHeight

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

		const svg = d3.select(".graph-content").append("svg")
			.attr("viewBox", "0 0 " + width + " " + height)
			.attr("width", "100%")
			.attr("preserveAspectRatio", "xMidYMid meet")
			.call(zoom)
		this.svg = svg

		const tooltip = d3.select(".graph-content").append("div")
			.attr("class", "tooltip")
			.style("opacity", 0)
		this.tooltip = tooltip

		function zoom(svg: any) {
			const extent: [[number, number], [number, number]] = [[margin.left, margin.top], [width - margin.right, height - margin.top]];

			const zoomSpec = d3.zoom()
				.scaleExtent([1, 4])
				.translateExtent(extent)
				.extent(extent)
				.on("zoom", zoomed)

			svg.call(zoomSpec)

			function zoomed() {
				const ns: [number, number] = [margin.left, width - margin.right]
				const mapped = ns.map(d => d3.event.transform.applyX(d)) as [number, number]
				x.range(mapped);
				svg.selectAll(".bars rect").attr("x", (d: Song) => x(d.ident)).attr("width", x.bandwidth());
				svg.selectAll(".x-axis").call(xAxis);

				// https://stackoverflow.com/questions/46005546/d3-v4-get-current-zoom-scale/46005996
				if (d3.event.transform.k === 1) {
					// TODO: unused
				}
			}
		}

		const self = this

		svg.append("g")
			.attr("class", "bars")
			.attr("fill", colors.httpPinkAlpha)
			.selectAll("rect")
			.data(data)
			.join("rect")
			.attr("x", (d) => x(d.ident)!)
			.attr("width", x.bandwidth())
			.on("mouseover", function(d) {
				d3.select(this).transition()
					.duration(0)
					.style("fill", colors.httpPink)

				tooltip.transition()
					.duration(200)
					.style("opacity", .9);
				tooltip.html(() => self.tooltipHTML(d))
					.style("left", (d3.event.pageX) + 20 + "px")
					.style("top", (d3.event.pageY - 80) + "px")
			})
			.on("mouseout", function() {
				d3.select(this).transition()
					.duration(0)
					.style("fill", colors.httpPinkAlpha)

				tooltip.transition()
					.duration(500)
					.style("opacity", 0);
			})
			.attr("y", (d) => y(0))
			.attr("height", (d) => 0)
			.transition()
				.duration(500)
				.attr("y", (d) => y(d.playCount))
				.attr("height", (d) => y(0) - y(d.playCount))
				.delay(function(d,i){return(i*20)})

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
			.attr("x", 150)
			.attr("y", 0)
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
						{this.props.data.slice(0, this.maxDataItems).map((d, i) => {
							return <tr key={d.ident}>
								<td>{i + 1}</td>
								<td>{d.title}</td>
								<td className="artist">{d.artistName}</td>
								<td className="album">{d.albumTitle}</td>
								<td>{d.playCount.toLocaleString()}</td>
							</tr>
						})}
					</tbody>
				</table>
			</div>
		</div>
	}
}

function hmsDisplay(h: number, m: number, s: number, long: boolean): [string, string, string] {
	const reth = h === 0 ?
		"" :
		(""+h) + (long ? pluralize(" hour", m) : "h")
	const retm = (long ? ""+m : (""+m).padStart(2, "0")) + (long ? pluralize(" minute", m) : "m")
	const rets = (long ? ""+s : (""+s).padStart(2, "0")) + (long ? pluralize(" second", s) : "s")
	return [reth, retm, rets]
}

type MostListenedArtistsProps = {
	data: ArtistPlayCountDataResponse
}

export class MostListenedArtists extends React.Component<MostListenedArtistsProps> {
	private svg: any
	private tooltip: any

	private readonly maxDataItems = 100

	componentDidMount() {
		this.draw()
	}

	componentDidUpdate(newProps: MostListenedArtistsProps) {
		if (this.props.data !== newProps.data) {
			this.draw()
		}
	}

	componentWillUnmount() {
		this.svg?.remove()
		this.tooltip?.remove()
	}

	private tooltipHTML(d: ArtistPlayCountDatum): string {
		const [h, m, s] = secondsToHms(d.totalPlayTime)
		const [hd, md, sd] = hmsDisplay(h, m, s, true)
		return `${d.artistName}<br/>${hd + " " + md + " " + sd}`
	}

	private draw() {
		// https://observablehq.com/@d3/zoomable-bar-chart

		if (this.props.data.length === 0) {
			return
		}

		const data = this.props.data.slice(0, this.maxDataItems)

		const margin = graphMargin
		const width = graphWidth
		const height = graphHeight

		const x = d3.scaleBand()
			.domain(data.map(d => d.artistName))
			.range([margin.left, width - margin.right])
			.padding(0.1)

		const y = d3.scaleTime()
			.domain([0, d3.max(data, d => { return d.totalPlayTime })!]).nice()
			.range([height - margin.bottom, margin.top])

		const yAxis = (g: any) => g
			.attr("transform", `translate(${margin.left},0)`)
			.call(d3.axisLeft(y).tickFormat((n) => {
				const [h, m, s] = secondsToHms(n as number)
				const [hd, md, ] = hmsDisplay(h, m, s, false)
				return hd + " " + md
			}))
			.call((g: any) => g.select(".domain").remove())

		const xAxis = (g: any) => g
			.attr("transform", `translate(0,${height - margin.bottom})`)
			.call(d3.axisBottom(x).tickSizeOuter(0).tickFormat(() => "").tickSize(0))

		const svg = d3.select(".graph-content").append("svg")
			.attr("viewBox", "0 0 " + width + " " + height)
			.attr("preserveAspectRatio", "xMidYMid meet")
			.call(zoom)
		this.svg = svg

		const tooltip = d3.select(".graph-content").append("div")
			.attr("class", "tooltip")
			.style("opacity", 0)
		this.tooltip = tooltip

		function zoom(svg: any) {
			const extent: [[number, number], [number, number]] = [[margin.left, margin.top], [width - margin.right, height - margin.top]];

			const zoomSpec = d3.zoom()
				.scaleExtent([1, 4])
				.translateExtent(extent)
				.extent(extent)
				.on("zoom", zoomed)

			svg.call(zoomSpec)

			function zoomed() {
				const ns: [number, number] = [margin.left, width - margin.right]
				const mapped = ns.map(d => d3.event.transform.applyX(d)) as [number, number]
				x.range(mapped);
				svg.selectAll(".bars rect").attr("x", (d: Song) => x(d.artistName)).attr("width", x.bandwidth());
				svg.selectAll(".x-axis").call(xAxis);

				// https://stackoverflow.com/questions/46005546/d3-v4-get-current-zoom-scale/46005996
				if (d3.event.transform.k === 1) {
					// TODO: unused
				}
			}
		}

		const self = this

		svg.append("g")
			.attr("class", "bars")
			.attr("fill", colors.purpleAlpha)
			.selectAll("rect")
			.data(data)
			.join("rect")
			.attr("x", (d) => x(d.artistName)!)
			.attr("width", x.bandwidth())
			.on("mouseover", function(d) {
				d3.select(this).transition()
					.duration(0)
					.style("fill", colors.purple)

				tooltip.transition()
					.duration(200)
					.style("opacity", .9);
				tooltip.html(() => self.tooltipHTML(d))
					.style("left", (d3.event.pageX) + 20 + "px")
					.style("top", (d3.event.pageY - 80) + "px")
			})
			.on("mouseout", function() {
				d3.select(this).transition()
					.duration(0)
					.style("fill", colors.purpleAlpha)

				tooltip.transition()
					.duration(500)
					.style("opacity", 0);
			})
			.attr("y", (d) => y(0))
			.attr("height", (d) => 0)
			.transition()
				.duration(500)
				.attr("y", (d) => y(d.totalPlayTime))
				.attr("height", (d) => y(0) - y(d.totalPlayTime))
				.delay(function(d,i){return(i*20)})

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
			.text("Artists");

		svg.append("text")
			.attr("class", "y label")
			.attr("text-anchor", "end")
			.attr("x", 150)
			.attr("y", 0)
			.attr("dy", ".75em")
			.text("↑ Time listened");
	}

	render() {
		if (this.props.data.length === 0) {
			return <div className="info">Not enough data yet. Scrobble away!</div>
		}

		return <div className="MostListenedArtists">
			<div className="graph-content"></div>
			<div className="instruction">You can hover over, zoom (scroll), and pan this graph.</div>
			<div className="table">
				<table>
					<thead>
						<tr>
							<th>#</th>
							<th className="artist">Artist</th>
							<th>Total duration</th>
							<th className="play-count">Play count</th>
						</tr>
					</thead>
					<tbody>
						{this.props.data.slice(0, this.maxDataItems).map((d, i) => {
							const [h, m, s] = secondsToHms(d.totalPlayTime)
							const [hd, md, sd] = hmsDisplay(h, m, s, false)

							return <tr key={d.artistName}>
								<td>{i + 1}</td>
								<td className="artist">{d.artistName}</td>
								<td>{hd + " " + md + " " + sd}</td>
								<td className="play-count">{d.playCount.toLocaleString()}</td>
							</tr>
						})}
					</tbody>
				</table>
			</div>
		</div>
	}
}

type LongestSongsProps = {
	data: SongsDataResponse
}

export class LongestSongs extends React.Component<LongestSongsProps> {
	private svg: any
	private tooltip: any

	private readonly maxDataItems = 125

	componentDidMount() {
		this.draw()
	}

	componentDidUpdate(newProps: LongestSongsProps) {
		if (this.props.data !== newProps.data) {
			this.draw()
		}
	}

	componentWillUnmount() {
		this.svg?.remove()
		this.tooltip?.remove()
	}

	private tooltipHTML(d: Song): string {
		const [h, m, s] = secondsToHms(d.totalTime / 10**9)
		const [hd, md, sd] = hmsDisplay(h, m, s, true)
		return `${d.artistName} – ${d.title}<br/>${(hd + " " + md + " " + sd).trim()}`
	}

	private draw() {
		// https://observablehq.com/@d3/zoomable-bar-chart

		if (this.props.data.length === 0) {
			return
		}

		const data = this.props.data.slice(0, this.maxDataItems)

		const margin = graphMargin
		const width = graphWidth
		const height = graphHeight

		const x = d3.scaleBand()
			.domain(data.map(d => d.ident))
			.range([margin.left, width - margin.right])
			.padding(0.1)

		const y = d3.scaleTime()
			.domain([0, d3.max(data, d => { return d.totalTime })!]).nice()
			.range([height - margin.bottom, margin.top])

		const yAxis = (g: any) => g
			.attr("transform", `translate(${margin.left},0)`)
			.call(d3.axisLeft(y).tickFormat((n) => {
				const [h, m, s] = secondsToHms(n as number / 10**9)
				const [hd, md, sd] = hmsDisplay(h, m, s, false)
				return (hd + " " + md + " " + sd).trim()
			}))
			.call((g: any) => g.select(".domain").remove())

		const xAxis = (g: any) => g
			.attr("transform", `translate(0,${height - margin.bottom})`)
			.call(d3.axisBottom(x).tickSizeOuter(0).tickFormat(() => "").tickSize(0))

		const svg = d3.select(".graph-content").append("svg")
			.attr("viewBox", "0 0 " + width + " " + height)
			.attr("preserveAspectRatio", "xMidYMid meet")
			.call(zoom)
		this.svg = svg

		const tooltip = d3.select(".graph-content").append("div")
			.attr("class", "tooltip")
			.style("opacity", 0)
		this.tooltip = tooltip

		function zoom(svg: any) {
			const extent: [[number, number], [number, number]] = [[margin.left, margin.top], [width - margin.right, height - margin.top]];

			const zoomSpec = d3.zoom()
				.scaleExtent([1, 4])
				.translateExtent(extent)
				.extent(extent)
				.on("zoom", zoomed)

			svg.call(zoomSpec)

			function zoomed() {
				const ns: [number, number] = [margin.left, width - margin.right]
				const mapped = ns.map(d => d3.event.transform.applyX(d)) as [number, number]
				x.range(mapped);
				svg.selectAll(".bars rect").attr("x", (d: Song) => x(d.ident)).attr("width", x.bandwidth());
				svg.selectAll(".x-axis").call(xAxis);

				// https://stackoverflow.com/questions/46005546/d3-v4-get-current-zoom-scale/46005996
				if (d3.event.transform.k === 1) {
					// TODO: unused
				}
			}
		}

		const self = this

		svg.append("g")
			.attr("class", "bars")
			.attr("fill", colors.yellowAlpha)
			.selectAll("rect")
			.data(data)
			.join("rect")
			.attr("x", (d) => x(d.ident)!)
			.attr("width", x.bandwidth())
			.on("mouseover", function(d) {
				d3.select(this).transition()
					.duration(0)
					.style("fill", colors.yellow)

				tooltip.transition()
					.duration(200)
					.style("opacity", .9);
				tooltip.html(() => self.tooltipHTML(d))
					.style("left", (d3.event.pageX) + 20 + "px")
					.style("top", (d3.event.pageY - 80) + "px")
			})
			.on("mouseout", function() {
				d3.select(this).transition()
					.duration(0)
					.style("fill", colors.yellowAlpha)

				tooltip.transition()
					.duration(500)
					.style("opacity", 0);
			})
			.attr("y", (d) => y(0))
			.attr("height", (d) => 0)
			.transition()
				.duration(500)
				.attr("y", (d) => y(d.totalTime))
				.attr("height", (d) => y(0) - y(d.totalTime))
				.delay(function(d,i){return(i*20)})

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
			.attr("x", 150)
			.attr("y", 0)
			.attr("dy", ".75em")
			.text("↑ Song length");
	}

	render() {
		if (this.props.data.length === 0) {
			return <div className="info">Not enough data yet. Scrobble away!</div>
		}

		return <div className="LongestSongs">
			<div className="graph-content"></div>
			<div className="instruction">You can hover over, zoom (scroll), and pan this graph.</div>
			<div className="table">
				<table>
					<thead>
						<tr>
							<th>#</th>
							<th>Title</th>
							<th className="artist">Artist</th>
							<th>Length</th>
						</tr>
					</thead>
					<tbody>
						{this.props.data.slice(0, this.maxDataItems).map((d, i) => {
							const [h, m, s] = secondsToHms(d.totalTime / 10**9)
							const [hd, md, sd] = hmsDisplay(h, m, s, false)

							return <tr key={d.ident}>
								<td>{i + 1}</td>
								<td>{d.title}</td>
								<td className="artist">{d.artistName}</td>
								<td>{(hd + " " + md + " " + sd).trim()}</td>
							</tr>
						})}
					</tbody>
				</table>
			</div>
		</div>
	}
}

// copy to other two graphs



