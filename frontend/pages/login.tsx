import { useRouter } from 'next/router'
import { FormEvent, useEffect, useState } from 'react'
import { useDataContext } from '../contexts/DataContext'
import styles from "../styles/Login.module.scss"


const Login = () => {
    const { push } = useRouter()
    const { signedIn, setSignedIn, checkSession } = useDataContext()
    const [password, setPassword] = useState<string>("")
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        checkSession()

    }, [])

    const handleLogin = (e: FormEvent) => {
        e.preventDefault()
        const body = { "password": password }
        fetch("/api/login", {
            method: "POST",
            body: JSON.stringify(body)
        }).then(res => {
            if (!res.ok) {
                return res.text().then(text => { throw new Error(text) })
            }
            else {
                setSignedIn(true)
                return push("/")
            }
        }).catch(err => {
            setError(err.message)
            setPassword("")
        });
    }

    if (signedIn === false) {

        return (
            <div className={styles.Login}>
                <h2>Login</h2>
                <form action="" onSubmit={(e) => handleLogin(e)}>
                    <label htmlFor="password">Password</label>
                    <input type="password" id='password' onChange={(e) => setPassword(e.target.value)} value={password} />
                    <button>Login</button>
                    {error &&
                        <div className={styles.Login_error}>
                            {error}
                        </div>
                    }
                </form>
            </div>
        )
    }
}

export default Login