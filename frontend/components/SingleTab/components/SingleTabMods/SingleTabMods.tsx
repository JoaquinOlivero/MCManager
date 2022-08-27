import styles from '../../../../styles/components/SingleTab/components/SingleTabMods/SingleTabMods.module.scss'
import JavaIcon from '../../../../svg/icons/JavaIcon'

type Mod = {
    "fileName": string,
    "modId": string,
    "version": string
}

type Props = {
    mods: Array<Mod>
}

const SingleTabMods = ({ mods }: Props) => {


    // add file info columns
    const headerItems = () => {
        const headerArr = []
        const element =
            <div className={styles.SingleTabMods_header_info}>
                <span className={styles.SingleTabMods_info_filename}>Mod filename</span>
                {/* <div className={styles.SingleTabMods_info_details}>
                    <span>Mod id</span>
                    <span>Version</span>
                </div> */}
            </div>
        if (mods.length === 1) {
            headerArr.push(element)
        } else {
            var i = 0
            while (i < 2) {
                headerArr.push(element)
                i++
            }
        }
        return headerArr
    }

    return (
        <div className={styles.SingleTabMods} key={"modsTab"}>
            <div className={styles.SingleTabMods_header}>
                {mods.length !== 0 && headerItems()}
            </div>

            <div className={styles.SingleTabMods_mods_container}>
                {mods.map((mod: Mod, i: number) => {
                    return <div key={mod.fileName} className={styles.SingleTabMods_mod_container} style={(i + 1) % 2 === 0 ? { borderRight: "none" } : {}}>
                        <div><JavaIcon /></div>
                        <div className={styles.SingleTabMods_mod_filename}>{mod.fileName}</div>
                        {/* <div className={styles.SingleTabMods_mod_id}>{mod.version}</div> */}
                    </div>
                })}
            </div>

        </div>
    )
}

export default SingleTabMods