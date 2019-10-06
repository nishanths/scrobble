import * as React from "react";
import * as ReactDOM from "react-dom";
import { U } from "../components/U"
import { UArgs } from "../shared/types";
import { createStore } from "redux";
import { Provider } from "react-redux";
import reducer from "../redux/reducers/u"

declare var uargs: UArgs;

const store = createStore(reducer, { uargs });

ReactDOM.render(
  <Provider store={store}>
    <U {...uargs} />
  </Provider>,
  document.querySelector("#app")
);
