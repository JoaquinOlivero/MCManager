import { useDataContext } from '../../../../../contexts/DataContext'
import { useState } from 'react'
import styles from '../../../../../styles/components/SingleTab/components/SingleTabSettings/SingleTabSettings.module.scss'

type Settings = {
    minecraft_directory: string
    run_method: string
    docker_container_id: string
    start_command: string
}

type Props = {
    settings: Settings
    getSettings: Function
}

const Command = ({ settings, getSettings }: Props) => {
    const [mcDir, setMcDir] = useState<string>(settings.run_method !== "docker" ? settings.minecraft_directory : "")
    const [startCommand, setStartCommand] = useState<string>(settings.run_method !== "docker" ? settings.start_command : "")
    const [isSaving, setisSaving] = useState<boolean>(false)
    const [responseError, setResponseError] = useState<null | string>(null)
    const { completeSettings, setCompleteSettings } = useDataContext()

    const handleSaveDirAndCommand = () => {
        if (mcDir === "" || startCommand === "") return

        setisSaving(true)
        setResponseError(null)

        fetch("/api/settings/command/save", {
            method: "POST",
            body: JSON.stringify({ "minecraft_directory": mcDir, "start_command": startCommand })
        }).then(res => {
            if (!res.ok) {
                return res.text().then(text => { throw new Error(text) })
            }
            getSettings()
            setisSaving(false)
            if (completeSettings === false) setCompleteSettings(true)
        }).catch(err => {
            setResponseError(err.message)
            setisSaving(false)
            setMcDir("")
            setStartCommand("")
        });
    }

    return (
        <div className={styles.SingleTabSettings_option_content}>
            <div className={styles.SingleTabSettings_content_title}>
                Minecraft Directory
                <input type="text" onChange={(e) => setMcDir(e.target.value)} value={mcDir} />
            </div>
            <div className={styles.SingleTabSettings_content_title}>
                Start Command
                <input type="text" onChange={(e) => setStartCommand(e.target.value)} value={startCommand} />
            </div>
            <div className={styles.SingleTabSettings_btn} onClick={handleSaveDirAndCommand} style={mcDir === "" || mcDir === settings.minecraft_directory && startCommand === settings.start_command || startCommand === "" || isSaving ? { opacity: 0.5, pointerEvents: "none" } : {}}>
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

export default Command