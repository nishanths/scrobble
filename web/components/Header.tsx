import * as React from "react"

export class Header extends React.Component<{username: string, signedIn: boolean}, {}> {
  render() {
    return <div className="header">
      <span className="username"><span className="emph">{this.props.username}</span><span className="rest">'s scrobbles</span></span>
      {this.props.signedIn && <span className="settings"><a href="/">Settings</a></span>}
    </div>
  }
}
