import { useDataContext } from '../../../../../contexts/DataContext'
import { useState } from 'react'
import styles from '../../../../../styles/components/SingleTab/components/SingleTabSettings/SingleTabSettings.module.scss'

type Settings = {
    minecraft_directory: string
    run_method: string
    docker_container_id: string
    start_command: string
}

type DockerContainer = {
    container_id: string
    container_name: string
}

type Props = {
    settings: Settings
    dockerContainers: Array<DockerContainer> | null
    getSettings: Function
}


const Docker = ({ settings, dockerContainers, getSettings }: Props) => {
    const [dockerContainerOption, setDockerContainerOption] = useState<string>('default')
    const [isConnecting, setIsConnecting] = useState<boolean>(false)
    const [isDisconnecting, setIsDisconnecting] = useState<boolean>(false)
    const { completeSettings, setCompleteSettings } = useDataContext()

    const handleConnectDockerContainer = async () => {
        setIsConnecting(true)
        if (dockerContainerOption === "default") return
        const res = await fetch("/api/settings/docker/connect", {
            method: "POST",
            body: JSON.stringify({ "container_id": dockerContainerOption })
        })

        if (res.status !== 200) {
            console.log("could not connect to docker container")
        }
        await getSettings()
        setIsConnecting(false)
        if (completeSettings === false) setCompleteSettings(true)
    }

    const handleDisconnectDockerContainer = async () => {
        setIsDisconnecting(true)
        const res = await fetch("/api/settings/docker/disconnect", {
            method: "POST",
        })

        if (res.status !== 200) {
            console.log("could not disconnect from docker container")
        }
        await getSettings()
        setIsDisconnecting(false)
        setCompleteSettings(false)
    }
    return (
        <div className={styles.SingleTabSettings_option_content}>
            <div className={styles.SingleTabSettings_content_title}>
                Docker Container
            </div>
            <select name="" id="" onChange={(e) => setDockerContainerOption(e.target.value)} value={settings.run_method === "docker" ? settings.docker_container_id : dockerContainerOption}>
                <option value="default" disabled hidden>Select Container</option>
                {dockerContainers ?
                    dockerContainers.map((container: DockerContainer) => {
                        return <option key={container.container_id} value={container.container_id}>{container.container_name.substring(1)}</option>
                    })
                    :
                    <option value="not-found">No docker container found</option>
                }
            </select>
            <div className={styles.SingleTabSettings_btn} onClick={handleConnectDockerContainer} style={settings.run_method === "docker" ? { pointerEvents: "none", opacity: 0.5 } : {}}>
                <span>{isConnecting ? "Connecting" : settings.run_method === 'docker' ? "Connected" : "Connect"}</span>
            </div>
            {settings.run_method === "docker" &&
                <div className={styles.SingleTabSettings_btn} onClick={handleDisconnectDockerContainer}>
                    <span>{isDisconnecting ? "Disconnecting" : "Disconnect"}</span>
                </div>
            }
        </div>
    )
}

export default Docker