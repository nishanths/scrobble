import React from "react";
import { Notie } from "../../shared/types"
import "../../scss/dashboard/delete-account.scss"

const deletePrompt = `Delete account?`
const failedDelete = `Failed to delete account. Try again or contact nishanths@utexas.edu.`

interface DeleteAccountProps {
    wnd: Window
    logoutURL: string
    notie: Notie
}

export const DeleteAccount: React.SFC<DeleteAccountProps> = ({ wnd, logoutURL, notie }) => {
    const doDelete = (): void => {
        fetch(`/api/v1/account/delete`, { method: "POST" })
            .then(
                r => {
                    if (r.status === 200) {
                        wnd.location.assign(logoutURL)
                        return
                    }
                    console.log("failed to delete: status=%d", r.status)
                    notie.alert({ type: "error", text: failedDelete, stay: true })
                },
                err => {
                    console.error(err)
                    notie.alert({ type: "error", text: failedDelete, stay: true })
                }
            )
    }

    const onClick = (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>): void => {
        e.preventDefault()
        const ok = wnd.confirm(deletePrompt)
        if (!ok) {
            return
        }
        doDelete()
    }

    return <div className="DeleteAccount">
        <div>Click to <a className="danger" href="" onClick={onClick}>delete account</a>.</div>
    </div>
}
