import { useDataContext } from '../../../../../contexts/DataContext'
import { useEffect, useState } from 'react'
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
    const [scriptName, setScriptName] = useState<string>(settings.run_method !== "docker" ? settings.start_command : "")
    const [isSaving, setisSaving] = useState<boolean>(false)
    const [responseError, setResponseError] = useState<null | string>(null)
    const [shFiles, setShFiles] = useState<Array<string> | null>(null)
    const { completeSettings, setCompleteSettings } = useDataContext()

    const handleSaveDirAndCommand = () => {
        if (mcDir === "" || scriptName === "") return

        setisSaving(true)
        setResponseError(null)

        fetch("/api/settings/command/save", {
            method: "POST",
            body: JSON.stringify({ "minecraft_directory": mcDir, "script": scriptName })
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
            setScriptName("")
        });
    }

    const handleGetScriptFiles = () => {
        if (mcDir === "") return

        fetch("/api/settings/command/files", {
            method: "POST",
            body: JSON.stringify({ "directory": mcDir })
        }).then(res => {
            if (!res.ok) {
                return res.text().then(text => { throw new Error(text) })
            }

            return res.json().then(json => {
                setShFiles(json.files)

            })
        }).catch(err => {
            console.log(err.message)
        });
    }

    return (
        <div className={styles.SingleTabSettings_option_content}>
            <div className={styles.SingleTabSettings_content_title}>
                Minecraft Directory
                <input type="text" onChange={(e) => setMcDir(e.target.value)} value={mcDir} />
            </div>
            <div className={styles.SingleTabSettings_content_title}>
                Script
                <select name="" id="" onMouseDown={handleGetScriptFiles} onChange={(e) => setScriptName(e.target.value)}>
                    <option value="default" hidden>Select File</option>
                    {shFiles ?
                        shFiles.map((file: string) => {
                            return <option key={file} value={file}>{file}</option>
                        })
                        :
                        <option value="not-found" disabled>No file found</option>
                    }
                </select>
            </div>
            <div className={styles.SingleTabSettings_btn} onClick={handleSaveDirAndCommand} style={mcDir === "" || mcDir === settings.minecraft_directory && scriptName === settings.start_command || scriptName === "" || isSaving ? { opacity: 0.5, pointerEvents: "none" } : {}}>
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