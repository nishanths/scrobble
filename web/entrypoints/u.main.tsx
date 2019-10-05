import * as React from "react";
import * as ReactDOM from "react-dom";
import { UsernamePage }  from "../components/UsernamePage"
import { UArgs } from "../shared/types";
import { createStore } from "redux";
import { Provider } from "react-redux";
import rootReducer from "../redux/reducers/username"

declare var uargs: UArgs;

const store = createStore(rootReducer, { uargs });

ReactDOM.render(
  <Provider store={store}>
  	<UsernamePage {...uargs}/>
  </Provider>,
  document.querySelector("#app")
);
