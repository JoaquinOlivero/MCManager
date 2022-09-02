import { useState } from 'react'
import styles from '../../../../styles/components/SingleTab/components/SingleTabSettings/SingleTabSettings.module.scss'
import Spinner from '../../../../svg/icons/Spinner'
import Docker from './components/Docker'

type Settings = {
    minecraft_directory: string
    run_method: string
    docker_container_id: string
    start_script: string
    stop_script: string
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

const SingleTabSettings = ({ settings, dockerContainers, getSettings }: Props) => {


    return (
        <div className={styles.SingleTabSettings}>
            {settings && dockerContainers ?
                <div className={styles.SingleTabSettings_option_section}>
                    <div className={styles.SingleTabSettings_option_header}>
                        <span>Minecraft Server Control</span>
                    </div>
                    <div className={styles.SingleTabSettings_option_container}>
                        <div className={styles.SingleTabSettings_option_content}>
                            <div className={styles.SingleTabSettings_content_title}>
                                Scripts
                            </div>
                        </div>
                        <Docker settings={settings} dockerContainers={dockerContainers} getSettings={getSettings} />
                        {/* <div className={styles.SingleTabSettings_option_content}>
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
                            <div className={styles.SingleTabSettings_btn_connect} onClick={handleConnectDockerContainer} style={settings.run_method === "docker" ? { pointerEvents: "none", opacity: 0.5 } : {}}>
                                <span>{isConnecting ? "Connecting" : settings.run_method !== '' ? "Connected" : "Connect"}</span>
                            </div>
                            {settings.run_method === "docker" &&
                                <div className={styles.SingleTabSettings_btn_connect} onClick={handleDisconnectDockerContainer}>
                                    <span>Disconnect</span>
                                </div>
                            }
                        </div> */}
                    </div>
                </div>
                :
                <Spinner />
            }
        </div>
    )
}

export default SingleTabSettings