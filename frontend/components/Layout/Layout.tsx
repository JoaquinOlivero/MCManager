import Link from 'next/link'
import { useRouter } from 'next/router'
import styles from '../../styles/components/Layout/Layout.module.scss'
import Variables from '../../styles/Variables.module.scss'
import Gear from '../../svg/icons/Gear'

type Props = {
    children: React.ReactNode
}

const Layout = ({ children }: Props) => {
    const { route } = useRouter()

    return (
        <div className={styles.Layout}>
            {/* Left menu persistent layout */}
            <div className={styles.Menu}>
                <h2>MCManager</h2>
                <div className={styles.Menu_content}>
                    <div className={styles.Menu_content_tabs}>
                        <div className={styles.Menu_tab}>
                            <Link href="/mods">
                                <span style={route === "/mods" ? { color: Variables.primaryColor } : {}}>Mods</span>
                            </Link>
                        </div>
                        <div className={styles.Menu_tab}>
                            <Link href="/config">
                                <span style={route === "/config" ? { color: Variables.primaryColor } : {}}>Config</span>
                            </Link>
                        </div>
                        {/* <div className={styles.Menu_tab} style={{ pointerEvents: "none" }}>
                            <Link href='/world'>
                                <span style={route === "/world" ? { color: Variables.primaryColor } : {}}>World</span>
                            </Link>
                        </div>
                        <div className={styles.Menu_tab} style={{ pointerEvents: "none" }}>
                            <Link href='/logs'>
                                <span style={route === "/logs" ? { color: Variables.primaryColor } : {}}>Logs</span>
                            </Link>
                        </div> */}
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