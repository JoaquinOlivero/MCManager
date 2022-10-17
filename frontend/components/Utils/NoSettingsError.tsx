import Link from "next/link"
import styles from "../../styles/components/Utils/NoSettingsError.module.scss"

const NoSettingsError = () => {
    return (
        <div className={styles.NoSettingsError}>
            Please configure your MCManager <Link href="/settings" >settings</Link>.
        </div>
    )
}

export default NoSettingsError