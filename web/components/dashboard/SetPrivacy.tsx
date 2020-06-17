import React from "react";
import { Notie } from "../../shared/types"
import { cookieAuthErrorMessage } from "../../shared/util"

interface SetPrivacyProps {
    privacy: boolean
    privacyChange: (v: boolean) => void
    notie: Notie
}

export const SetPrivacy: React.SFC<SetPrivacyProps> = ({ privacy, privacyChange, notie }) => {
    const onClick = (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
        e.preventDefault()
        privacy ? setPrivacy(false) : setPrivacy(true)
    }

    const setPrivacy = (priv: boolean): void => {
        const genericErr = "Failed to update privacy. Try again."

        fetch(`/setPrivacy?privacy=${priv.toString()}`, { method: "POST" })
            .then(
                r => {
                    if (r.status === 200) {
                        notie.alert({ type: "success", text: "Updated privacy." })
                        privacyChange(priv)
                    } else {
                        console.log("failed to update privacy: status=%d", r.status)
                        notie.alert({ type: "error", text: (r.status === 401) ? cookieAuthErrorMessage : genericErr })
                    }
                },
                err => {
                    console.error(err)
                    notie.alert({ type: "error", text: genericErr })
                }
            )
    }

    return <div className="SetPrivacy">
        <div>Switch profile <a href="" onClick={onClick}>to {privacy ? "public" : "private"}</a>.</div>
    </div>
}
