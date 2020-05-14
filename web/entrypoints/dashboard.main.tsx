import React from "react"
import * as ReactDOM from "react-dom"
import { Dashboard } from "../components/dashboard"
import { BootstrapArgs } from "../shared/types"
import "../scss/dashboard/dashboard.scss"

declare const bootstrap: BootstrapArgs;

ReactDOM.render(
  <Dashboard {...bootstrap} />,
  document.querySelector("#app")
);
