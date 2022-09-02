import React, { useEffect, useState } from 'react'
import SingleTabSettings from '../components/SingleTab/components/SingleTabSettings/SingleTabSettings'
import SingleTab from '../components/SingleTab/SingleTab'
import SingleTabHeader from '../components/SingleTab/SingleTabHeader'

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


const Settings = () => {
    const [settings, setSettings] = useState<Settings | null>(null)
    const [docker, setDocker] = useState<Array<DockerContainer> | null>(null)

    const getSettings = async () => {
        type Data = {
            settings: Settings
            docker_containers: Array<DockerContainer>
        }
        const res = await fetch("/api/settings")
        const data: Data = await res.json()
        await setSettings(data.settings)
        await setDocker(data.docker_containers)
    }

    useEffect(() => {
        getSettings()

        return () => {
            setSettings(null)
            setDocker(null)
        }
    }, [])


    return (
        <SingleTab header={<SingleTabHeader tabType={"settings"} />}>
            {settings && <SingleTabSettings settings={settings} dockerContainers={docker} getSettings={getSettings} />}
        </SingleTab>
    )
}

export default Settings

