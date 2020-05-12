import React from "react";
import * as ReactDOM from "react-dom";
import { Root } from "../components/Root";
import { BootstrapArgs } from "../shared/types";

declare let bootstrap: BootstrapArgs;

ReactDOM.render(
  <Root {...bootstrap} />,
  document.querySelector("#app")
);
