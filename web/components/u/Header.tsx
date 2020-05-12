import React from "react"

export const Header: React.StatelessComponent<{ username: string, signedIn: boolean }> = ({
  username,
  signedIn
}) => {
  return <div className="header">
    <span className="username"><span className="emph">{username}</span><span className="rest">'s scrobbles</span></span>
    <span className="nav">
      {signedIn ? <a href="/">Settings</a> : <a href="/login">Sign In</a>}
    </span>
  </div>
}
