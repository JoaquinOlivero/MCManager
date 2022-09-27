import Link from 'next/link'
import { useRouter } from 'next/router'
import styles from '../../styles/components/Layout/Layout.module.scss'
import Variables from '../../styles/Variables.module.scss'
import Gear from '../../svg/icons/Gear'
import { useDataContext } from "../../contexts/DataContext";

type Props = {
    children: React.ReactNode
}

const Layout = ({ children }: Props) => {
    const { route, push } = useRouter()
    const { editFilepath, setEditFilepath } = useDataContext()

    const handleClickServerProperties = async () => {
        await setEditFilepath("/server.properties")
        push("/edit")
    }

    return (
        <div className={styles.Layout}>
            {/* Left menu persistent layout */}
            <div className={styles.Menu}>
                <Link href='/'>
                    <h2>MCManager</h2>
                </Link>
                <div className={styles.Menu_content}>
                    <div className={styles.Menu_content_tabs}>
                        <div className={styles.Menu_tab}>
                            <Link href="/mods">
                                <span style={route.includes("mods") ? { color: Variables.primaryColor } : {}}>Mods</span>
                            </Link>
                        </div>
                        <div className={styles.Menu_tab}>
                            <Link href="/config">
                                <span style={route.includes("config") ? { color: Variables.primaryColor } : {}}>Config</span>
                            </Link>
                        </div>
                        <div className={styles.Menu_tab}>
                            <span style={editFilepath === "/server.properties" ? { color: Variables.primaryColor } : {}} onClick={handleClickServerProperties}>Server.properties</span>
                        </div>
                        <div className={styles.Menu_tab}>
                            <Link href='/world'>
                                <span style={route.includes("world") ? { color: Variables.primaryColor } : {}}>World</span>
                            </Link>
                        </div>
                        <div className={styles.Menu_tab}>
                            <Link href='/logs'>
                                <span style={route.includes("logs") ? { color: Variables.primaryColor } : {}}>Logs</span>
                            </Link>
                        </div>
                    </div>

                </div>

                <div className={styles.Menu_settings}>
                    <div className={styles.Menu_settings_tabs}>
                        <Link href='/settings'>
                            <div style={route === "/settings" ? { color: Variables.primaryColor, opacity: 1 } : {}} className={styles.Menu_tab}>
                                <Gear fill={route === "/settings" ? Variables.primaryColor : "#ccc"} />
                                <span>Settings</span>
                            </div>
                        </Link>
                    </div>
                </div>
            </div>

            {/* Page */}
            <div className={styles.Page}>
                {children}
            </div>
        </div>
    )
}

export default Layout