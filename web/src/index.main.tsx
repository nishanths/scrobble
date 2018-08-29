import * as React from "react";
import * as ReactDOM from "react-dom";
import { Index } from "../components/Index";
import { BootstrapArgs } from "./shared";
import "../scss/index.main.scss";

declare var bootstrap: BootstrapArgs;

ReactDOM.render(
  <Index p={bootstrap}/>,
  document.getElementById("app")
);
