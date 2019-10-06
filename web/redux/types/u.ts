import reducer from "../reducers/u"

export type State = ReturnType<typeof reducer>
export type PartialState = Partial<State>
