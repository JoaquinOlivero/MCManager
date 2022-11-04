import Link from 'next/link'
import { useRouter } from 'next/router'
import styles from '../../styles/components/Layout/Layout.module.scss'
import Variables from '../../styles/Variables.module.scss'
import Gear from '../../svg/icons/Gear'
import { useDataContext } from "../../contexts/DataContext";
import { useEffect } from 'react'
import SignOut from '../../svg/icons/SignOut'

type Props = {
    children: React.ReactNode
}

const Layout = ({ children }: Props) => {
    const { route, push } = useRouter()
    const { editFilepath, setEditFilepath, signedIn, checkSession } = useDataContext()


    useEffect(() => {
        if (signedIn !== true) {
            if (signedIn === false) push("/login")
        }
    }, [signedIn])

    useEffect(() => {
        checkSession()
    }, [])


    const handleClickServerProperties = async () => {
        await setEditFilepath("/server.properties")
        push("/edit")
    }

    const handleSignOut = () => {
        fetch("/api/logout", {
            method: "GET",
            credentials: "include"
        }).then((res) => {
            if (!res.ok) {
                return res.text().then(text => { throw new Error(text) })
            } else {
                push("/login")
            }
        })
            .catch(err => {
                console.log(err);
            });
    }

    return (
        <div className={styles.Layout}>
            {/* Left menu persistent layout */}
            {signedIn && !route.includes("login") ?
                <>
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

                                <div className={styles.Menu_tab} onClick={handleSignOut}>
                                    <SignOut />
                                    <span>Sign Out</span>
                                </div>
                            </div>
                        </div>
                    </div>
                    {/* Page */}
                    <div className={styles.Page}>
                        {children}
                    </div>
                </>

                :

                <>
                    {route.includes("login") ?
                        <div className={styles.Page}>
                            {children}
                        </div>
                        :
                        <></>
                    }
                </>
            }

        </div >
    )
}

export default Layout