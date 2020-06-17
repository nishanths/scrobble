import React, { useState } from "react";
import { Notie } from "../../shared/types"
import { cookieAuthErrorMessage } from "../../shared/util"

interface NewAPIKeyProps {
    apiKeyChange: (s: string) => void
    notie: Notie
}

const genericErr = "Failed to regenerate API key. Try again."

export const NewAPIKey: React.SFC<NewAPIKeyProps> = ({ apiKeyChange, notie }) => {
    const [generating, setGenerating] = useState(false)

    const onClick = (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
        e.preventDefault()
        const ok = window.confirm("Go ahead and regenerate API key? Old keys will no longer function after.")
        if (!ok) {
            return
        }
        if (generating) { return }
        newAPIKey()
    }

    const newAPIKey = () => {
        setGenerating(true)
        let success = false

        fetch(`/newAPIKey`, { method: "POST" })
            .then(res => {
                if (res.status === 200) {
                    success = true
                    return res.json()
                }
                if (res.status === 401) {
                    notie.alert({ type: "error", text: cookieAuthErrorMessage })
                    return res.blob()
                }
                console.log("failed to generate API key: status=%d", res.status)
                notie.alert({ type: "error", text: genericErr })
                return res.blob()
            })
            .then(
                r => {
                    if (!success) { return }
                    apiKeyChange(r as string)
                    notie.alert({ type: "success", text: "Regenerated API key." })
                },
                err => {
                    console.error(err)
                    notie.alert({ type: "error", text: genericErr })
                }
            ).then(() => { // TODO: finally(), if available, would be nicer
                setGenerating(false)
            })
    }

    return <div className="NewAPIKey">
        <div>If you wish you may <a href="" onClick={onClick}>regenerate the API key</a>. (Old keys will no longer function after.)</div>
    </div>
}
