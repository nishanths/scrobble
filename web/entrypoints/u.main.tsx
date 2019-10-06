import React from "react";
import thunk from 'redux-thunk'
import * as ReactDOM from "react-dom";
import { createStore, applyMiddleware } from "redux";
import { Provider } from "react-redux";
import { U } from "../components/Ux"
import { UArgs } from "../shared/types";
import reducer from "../redux/reducers/u"

declare var uargs: UArgs;

const store = createStore(reducer, { uargs }, applyMiddleware(thunk));

ReactDOM.render(
  <Provider store={store}>
    <U {...uargs} wnd={window} />
  </Provider>,
  document.querySelector("#app")
);
