import * as React from "react";

export interface PlaybackProps {
  artist: string;
  album: string;
  name: string;
  artwork: string;
  played: Date[];
  count: number;
  appleMusic: string;
  now: Date;
}

function shortMonth(n: number): string {
  if (n < 0 || n > 11) {
    throw "bad month";
  }
  return ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep",
    "Oct", "Nov", "Dec"][n];
}

function dateDisplay(d: Date, now: Date): string {
  const s = now.getFullYear() != d.getFullYear() ?
    `${d.getDate()} ${shortMonth(d.getMonth())}` :
    `${d.getDate()} ${shortMonth(d.getMonth())} ${d.getFullYear()},`;

  const [h, period] = d.getHours() > 12 ?
    [d.getHours()-12, "pm"] :
    [d.getHours()-12, "am"];

  const m = d.getMinutes() < 10 ?
    "0" + d.getMinutes() :
    ""  + d.getMinutes();

  return `${s} ${h}:${m}${period}`;
}

export class Playback extends React.Component<PlaybackProps, {}> {
  render() {
    return (
      <div>
        <picture>
          <source srcSet={this.props.artwork} />
          <img src={this.props.artwork} />
        </picture>
        <span className="artist">{this.props.artist}</span>
        &nbsp;—&nbsp;
        <span className="name">{this.props.name}</span>
        <span className="album">{this.props.album}</span>
        <span className="played">{dateDisplay(this.props.played[0], this.props.now)}</span>
        <a className="appleMusic" href={this.props.appleMusic}></a>
      </div>
    );
  }
}
