import React from "react"

export const Header: React.StatelessComponent<{ username: string, signedIn: boolean, showNav: boolean }> = ({
  username,
  signedIn,
  showNav,
}) => {
  return <div className="header">
    <span className="username"><span className="emph">{username}</span><span className="rest">'s scrobbles</span></span>
    {showNav && <span className="nav">
      {signedIn ? <a href="/">Settings</a> : <><a href="/">Home</a> / <a href="/login">Sign In</a></>}
    </span>}
  </div>
}
