import React from "react"
import "../scss/google-sign-in.scss"

export const GoogleSignIn: React.SFC<{ loginURL: string }> = ({ loginURL }) => {
    return <a href={loginURL}>
        <div className="GoogleSignIn">
            <div className="container">
                <div className="img"></div>
                <div className="text">Sign in with Google</div>
            </div>
        </div>
    </a>
}
