import React from "react"
import { guideLink } from "../../shared/const"
import "../../scss/dashboard/base.scss"

type BaseProps = {
  username: string
  nSongs: number
}

export class Base extends React.Component<BaseProps> {
  constructor(props: BaseProps) {
    super(props)
  }

  private count(): JSX.Element {
  	const profile = <>Check out <a href={"/u/" + this.props.username}>your profile</a>.</>

  	if (this.props.nSongs === -1) {
  		return <p>{profile}</p>
  	}
  	if (this.props.nSongs === 0) {
		return <>
	        <p>{profile}</p>
			<p>You do not have any scrobbled songs. Find out <a href={guideLink} target="_blank" rel="">how to scrobble</a> in the guide.</p>
		</>
  	}
	return <p>You have {this.props.nSongs.toLocaleString()} scrobbled songs. {profile}</p>
  }

  render() {
    return <div className="Base">
      <p className="welcome">Welcome, {this.props.username}!</p>
      {this.count()}
    </div>
  }
}
