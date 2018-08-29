import * as React from "react";
import * as ReactDOM from "react-dom";
import { UsernamePage }  from "../components/UsernamePage"
import { UArgs } from "./shared";

declare var uargs: UArgs;

ReactDOM.render(
  <UsernamePage {...uargs}/>,
  document.querySelector("#app")
);
