import React from "react"
import { guideLink } from "../../shared/const"
import { sameDate, shortMonth } from "../../shared/time"
import "../../scss/dashboard/base.scss"

type BaseProps = {
  username: string
  nSongs: number
  lastScrobbleTime: number
}

export class Base extends React.Component<BaseProps> {
  constructor(props: BaseProps) {
    super(props)
  }

  private count(): JSX.Element {
    const profile = <>Check out <a href={"/u/" + this.props.username}>your profile</a>.</>

    if (this.props.nSongs === 0) {
      return <>
        <p>{profile}</p>
        <p>You do not have any scrobbled songs. Find out <a href={guideLink} target="_blank" rel="noopener noreferrer">how to scrobble</a> in the guide.</p>
      </>
    }

    const songs = this.props.nSongs === -1 ?
      null :
      <>You have {this.props.nSongs.toLocaleString()} scrobbled songs.</>

    const time = this.props.lastScrobbleTime === -1 ?
      null :
      <>The last time you scrobbled is <span title={new Date(this.props.lastScrobbleTime * 1000).toString()}>{lastScrobbledDisplay(new Date(this.props.lastScrobbleTime * 1000))}</span>.</>

    return <>
      <p>{profile}</p>
      <p><>{songs} </>{time}</p>
    </>
  }

  render() {
    return <div className="Base">
      <p className="welcome">Welcome, {this.props.username}!</p>
      {this.count()}
    </div>
  }
}

const lastScrobbledDisplay = (d: Date): string => {
  const now = new Date()

  if (sameDate(now, d)) {
    let h = d.getHours()
    let ampm = ""
    if (h > 12) {
      h = h - 12
      ampm = "pm"
    } else {
      ampm = "am"
    }
    console.log(h, d.getMinutes())
    return `${h}:${d.getMinutes()} ${ampm}`
  }

  return d.getFullYear() != now.getFullYear() ?
    `${d.getDate()} ${shortMonth(d)} ${d.getFullYear()}` :
    `${d.getDate()} ${shortMonth(d)}`
}
