import { useState } from 'react'
import styles from '../../../../../styles/components/SingleTab/components/SingleTabSettings/SingleTabSettings.module.scss'
import variables from "../../../../../styles/components/SingleTab/components/SingleTabSettings/SingleTabSettings.module.scss"

const ChangePassword = () => {
    const [oldPassword, setOldPassword] = useState<string>("")
    const [newPassword, setNewPassword] = useState<string>("")
    const [isSaving, setIsSaving] = useState<boolean>(false)
    const [responseError, setResponseError] = useState<null | string>(null)

    const handleSaveNewPassword = () => {
        setIsSaving(true)

        const body = { "old_password": oldPassword, "new_password": newPassword }

        fetch("/api/password/change", {
            method: "POST",
            body: JSON.stringify(body)
        }).then(res => {
            if (!res.ok) {
                return res.text().then(text => { throw new Error(text) })
            }
            else {
                setOldPassword("")
                setNewPassword("")
                setResponseError(null)
            }
        }).catch(err => {
            setOldPassword("")
            setNewPassword("")
            setResponseError(err.message)
        });

        setIsSaving(false)
    }

    return (
        <div className={styles.SingleTabSettings_option_content}>
            <div className={styles.SingleTabSettings_content_title}>Change password</div>

            <div className={styles.SingleTabSettings_content_title}>
                <input type="password" onChange={(e) => setOldPassword(e.target.value)} value={oldPassword} placeholder="Current Password" style={responseError && responseError === "Wrong password" ? { outline: "solid 2px " + variables.primaryColor } : {}} />
            </div>

            <div className={styles.SingleTabSettings_content_title}>
                <input type="password" onChange={(e) => setNewPassword(e.target.value)} value={newPassword} placeholder="New Password" />
            </div>

            <div className={styles.SingleTabSettings_btn} onClick={handleSaveNewPassword} style={oldPassword === "" || newPassword === "" || isSaving ? { opacity: 0.5, pointerEvents: "none" } : {}}>
                {isSaving ? "Saving" : "Save"}
            </div>

            {responseError &&
                <div className={styles.SingleTabSettings_error}>
                    {"Not saved. " + responseError + "."}
                </div>
            }
        </div>
    )
}

export default ChangePassword