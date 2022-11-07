import styles from '../../../../../styles/components/SingleTab/components/SingleTabSettings/SingleTabSettings.module.scss'
import variables from "../../../../../styles/components/SingleTab/components/SingleTabSettings/SingleTabSettings.module.scss"
import { Toggle } from '../../../../Utils/ToggleButton'

type BackupSettings = {
    world: boolean,
    mods: boolean,
    config: boolean,
    server_properties: boolean
}

type Props = {
    backup: BackupSettings
}

const Backup = ({ backup }: Props) => {

    const handleToggleBackup = (option: string, value: boolean) => {
        const body = { "option": option, "value": value }

        fetch("/api/settings/backup", {
            method: "POST",
            credentials: "include",
            body: JSON.stringify(body)
        }).then(res => {
            if (!res.ok) {
                return res.text().then(text => { throw new Error(text) })
            }
            else {

            }
        }).catch(err => {
            console.log(err.message)
        });
    }

    return (
        <div>
            <Toggle label='world' toggled={backup.world} onClick={handleToggleBackup} />

            <Toggle label='mods' toggled={backup.mods} onClick={handleToggleBackup} />

            <Toggle label='config' toggled={backup.config} onClick={handleToggleBackup} />

            <Toggle label='server_properties' toggled={backup.server_properties} onClick={handleToggleBackup} />
        </div>
    )
}

export default Backup