import { UArgs } from "../../shared/types"

export const uargsReducer = (s: UArgs | undefined): UArgs => {
  return s !== undefined ? s : {
    artworkBaseURL: "",
    host: "",
    self: false,
    profileUsername: "",
    logoutURL: "",
    account: {
      apiKey: "",
      username: "",
      private: false,
    },
    private: false,
  }
}
