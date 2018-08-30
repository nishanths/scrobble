import * as React from "react";
import * as ReactDOM from "react-dom";
import { IndexPage } from "../components/IndexPage";
import { BootstrapArgs } from "./shared";

declare var bootstrap: BootstrapArgs;

ReactDOM.render(
  <IndexPage {...bootstrap}/>,
  document.querySelector("#app")
);
