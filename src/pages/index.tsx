import type { NextPage } from 'next'
import { useEffect, useState } from 'react'
import SingleTab from '../components/SingleTab/SingleTab'
import SingleTabHeader from '../components/SingleTab/SingleTabHeader'
import styles from '../styles/Home.module.scss'
import Spinner from '../svg/icons/Spinner'
import { useRouter } from "next/router";
import NoSettingsError from '../components/Utils/NoSettingsError'

type Data = {
  run_method: string
  docker_status: string
  docker_health: string
  start_command: string
  command_status: string
  rcon_enabled: boolean
  rcon_port: string
  rcon_password: string
  ping_data: {
    description: string
    favicon: string
    players: {
      max: number
      online: number
      sample: [{ id: string, name: string }] | null
    }
    version: {
      name: string
      protocol: string
    }
  }
}

const Home: NextPage = () => {
  const router = useRouter();
  const [settings, setSettings] = useState<boolean>(false)
  const [serverInfo, setServerInfo] = useState<Data | null>(null)
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [isStopping, setIsStopping] = useState<boolean>(false)
  const [isStarting, setIsStarting] = useState<boolean>(false)
  const [rconValue, setRconValue] = useState<string>("")
  const [rconResponse, setRconResponse] = useState<string | null>(null)
  const [backupMsg, setBackupMsg] = useState<string | null>(null)

  const getHomeData: Function = async () => {
    const res = await fetch("/api")

    if (res.status === 200) {
      const data: Data = await res.json()

      // Check that config.json "run_method" (response from fetch) is not empty.
      if (data.run_method && data.run_method !== "") {
        setSettings(true)
      }

      // When the minecraft server runs in a docker container.
      if (data.run_method === "docker") {
        if (data.docker_status === "running" && data.docker_health === "starting") {
          setIsStarting(true)
          setIsLoading(false)
          return setTimeout(async () => await getHomeData(), 5000)
        } else if (data.docker_status === "running" && data.docker_health === "healthy") {
          setServerInfo(data)
          setIsStarting(false)
          return setIsLoading(false)
        } else if (data.docker_status === "running" && data.docker_health === "unhealthy") {
          setIsStarting(true)
          setIsLoading(false)
          return setTimeout(async () => await getHomeData(), 5000)
        }
        setIsStarting(false)
        setIsLoading(false)
        setServerInfo(null)
        return
      }

      // When the minecraft server runs bare bones with a command.
      if (data.run_method === "command") {
        if (data.command_status === "starting") {
          setIsStarting(true)
          setIsLoading(false)
          return setTimeout(async () => await getHomeData(), 5000)
        } else if (data.command_status === "online") {
          if (data.ping_data.version.name === "") {
            return setTimeout(async () => await getHomeData(), 5000)
          }
          setServerInfo(data)
          setIsStarting(false)
          return setIsLoading(false)
        }
        setIsStarting(false)
        setIsLoading(false)
        setServerInfo(null)
        return
      }

      setServerInfo(null)
      setIsStarting(false)
      setIsLoading(false)
      return

    } else {
      setServerInfo(null)
      setIsStarting(false)
      setIsLoading(false)
      return
    }

  }

  const handleServerStart = async () => {
    setIsStarting(true)
    try {
      const res = await fetch("/api?action=start", {
        method: "POST"
      })
      if (res.status === 200) {
        await getHomeData()
        return
      }
    } catch (error) {
      console.log(error);
    }
  }

  const handleServerStop = async () => {
    setIsStopping(true)
    try {
      const res = await fetch("/api?action=stop", {
        method: "POST"
      })

      if (res.status === 200) {
        await getHomeData()
        return setIsStopping(false)
      }
    } catch (error) {
      console.log(error)
    }

  }

  const handleSendRcon = () => {
    if (!rconValue || rconValue === "" || !serverInfo) return


    const body = {
      "rcon_command": rconValue
    }

    fetch("/api/rcon", {
      method: "POST",
      body: JSON.stringify(body)
    }).then(res => {
      if (!res.ok) {
        return res.text().then(text => { setRconResponse(text) })
      }
      else {
        return res.text().then(data => {
          setRconResponse(data)
          setRconValue("")
        })
      }
    })
      .catch(err => {
        setRconResponse(err.msg)
      });
  }

  const handleDownloadBackup = () => {
    // Update backup message state.
    setBackupMsg("Preparing backup. This may take a while...")

    fetch("/api/backup", {
      method: "GET",
      credentials: "include"
    })
      .then(res => {
        return res.text()
          .then(filename => {
            router.push("api/backup/download/" + filename)
            setBackupMsg(null)
          })
          .catch(err => {
            setBackupMsg(err.message)
          })
      })
      .catch(err => {
        // console.log(err.message)
        setBackupMsg(err.message)
      })
  }

  useEffect(() => {
    const asPathNestedRoutes = router.asPath.split("/").filter((v) => v.length > 0);
    if (asPathNestedRoutes.length > 0) {
      router.push(router.asPath)
    }
    getHomeData()

    return () => {
      setServerInfo(null)
      setRconValue("")
    }
  }, [])


  return (
    <SingleTab header={<SingleTabHeader tabType={"home"} />}>
      {!isLoading ?
        serverInfo && settings ?
          <div className={styles.Home}>

            {/* MOTD */}
            <h1>{serverInfo.ping_data.description}</h1>

            {/* Server Status  */}
            <div className={styles.Home_status}>
              <div><span className={styles.Home_status_title}>Server Status: </span><span className={styles.Home_status_server} style={{ color: "rgba(96, 230, 18, 1)", textShadow: "rgba(96, 230, 18, 0.5) 0px 0px 4px" }}>Online</span></div>
              <div className={styles.Home_status_control}>
                <div className={styles.Home_control_btn} style={{ borderColor: "#f37d79", pointerEvents: isStopping ? "none" : "auto", opacity: isStopping ? 0.5 : 1 }} onClick={handleServerStop} >
                  {isStopping ? "Stopping" : "Stop"}
                </div>
              </div>
            </div>

            <div className={styles.Home_content}>

              {/* More server information */}
              <div className={styles.Home_content_ping_data}>

                {/* Online Players  */}
                <div className={styles.Home_content_ping_data_item}>
                  <span className={styles.Home_content_data_item_title}>Online Players: </span><span className={styles.Home_content_data_item_info}>{serverInfo.ping_data.players.online}/{serverInfo.ping_data.players.max}</span>
                </div>

                {/* Server Version */}
                <div className={styles.Home_content_ping_data_item}>
                  <span className={styles.Home_content_data_item_title}>Server Version: </span> <span className={styles.Home_content_data_item_info}>{serverInfo.ping_data.version.name}</span>
                </div>

                {/* Current Players */}
                {serverInfo.ping_data.players.sample &&
                  <div className={styles.Home_content_ping_data_item}>
                    <span className={styles.Home_content_data_item_title}>Currently Playing: </span>
                    <ul>
                      {serverInfo.ping_data.players.sample.map(player => {
                        return <li key={player.name}>{player.name}</li>
                      })}
                    </ul>
                  </div>
                }

              </div>

              {/* Server actions  */}
              <div className={styles.Home_content_actions}>

                {/* RCON */}
                {serverInfo.rcon_enabled &&
                  <div className={styles.Home_content_actions_rcon}>
                    <span className={styles.Home_content_rcon_title}>Rcon</span>
                    <input type="text" onChange={(e) => setRconValue(e.target.value)} onSubmit={handleSendRcon} value={rconValue} />
                    <div className={styles.Home_content_rcon_btn} onClick={handleSendRcon}>
                      <span>Send</span>
                    </div>
                    {rconResponse &&
                      <div className={styles.Home_content_rcon_response}>
                        {rconResponse}
                      </div>}
                  </div>
                }

                {/* Backup */}
                <div className={styles.Home_content_actions_backup}>
                  <span className={styles.Home_content_backup_title}>Backup</span>
                  <div className={styles.Home_content_backup_btn} onClick={handleDownloadBackup}>Download</div>
                  {backupMsg && <div className={styles.Home_content_backup_message}>
                    {backupMsg} <div className={styles.Home_content_backup_message_spinner}></div>
                  </div>}
                </div>
              </div>
            </div>
          </div>
          :
          settings ?
            <div className={styles.Home} style={{ display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center" }}>

              <div className={styles.Home_status} style={{ borderBottom: "none" }}>
                <div><span className={styles.Home_status_title}>Server Status: </span><span className={styles.Home_status_server} style={{ color: "rgba(255, 60, 60, 1)", textShadow: "rgba(255, 60, 60, 0.5) 0px 0px 4px" }}>Offline</span></div>
                <div className={styles.Home_status_control}>
                  <div className={styles.Home_control_btn} style={{ borderColor: "#79d0bf", pointerEvents: isStarting ? "none" : "auto", opacity: isStarting ? 0.5 : 1 }} onClick={handleServerStart}>
                    {isStarting ? "Starting" : "Start"}
                  </div>
                </div>
              </div>

            </div>
            :
            <NoSettingsError />
        :
        <Spinner />
      }
    </SingleTab>
  )
}

export default Home


