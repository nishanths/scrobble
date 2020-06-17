import React from "react"
import * as ReactDOM from "react-dom"
import { Home } from "../components/home"
import { BootstrapArgs } from "../shared/types"
import "../scss/home/home.scss"

declare const bootstrap: BootstrapArgs;

ReactDOM.render(
	<Home loginURL={bootstrap.loginURL} />,
	document.querySelector("#app")
);
