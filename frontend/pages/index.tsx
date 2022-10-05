import type { NextPage } from 'next'
import { useEffect, useState } from 'react'
import SingleTab from '../components/SingleTab/SingleTab'
import SingleTabHeader from '../components/SingleTab/SingleTabHeader'
import styles from '../styles/Home.module.scss'
import Spinner from '../svg/icons/Spinner'
import { useRouter } from "next/router";

type Data = {
  docker_status: string
  docker_health: string
  rcon_enabled: boolean
  rcon_port: string
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
  const [serverInfo, setServerInfo] = useState<Data | null>(null)
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [isStopping, setIsStopping] = useState<boolean>(false)
  const [isStarting, setIsStarting] = useState<boolean>(false)

  const getHomeData: Function = async () => {
    const res = await fetch("/api")

    if (res.status === 200) {
      const data: Data = await res.json()
      if (data.docker_status === "running" && data.docker_health === "starting") {
        setIsStarting(true)
        setIsLoading(false)
        return setTimeout(async () => await getHomeData(), 5000)
      }

      setServerInfo(data)
      setIsStarting(false)
      return setIsLoading(false)
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
      const res = await fetch("/api?action=start&method=docker", {
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
      const res = await fetch("/api?action=stop&method=docker", {
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

  useEffect(() => {
    const asPathNestedRoutes = router.asPath.split("/").filter((v) => v.length > 0);
    if (asPathNestedRoutes.length > 0) {
      router.push(router.asPath)
    }
    getHomeData()

    return () => {
      setServerInfo(null)
    }
  }, [])


  return (
    <SingleTab header={<SingleTabHeader tabType={"home"} />}>
      {!isLoading ?
        serverInfo && serverInfo.docker_status === "running" ?
          <div className={styles.Home}>
            <h1>{serverInfo.ping_data.description}</h1>
            <div className={styles.Home_status}>
              <div><span className={styles.Home_status_title}>Server Status: </span><span className={styles.Home_status_server} style={{ color: "rgba(96, 230, 18, 1)", textShadow: "rgba(96, 230, 18, 0.5) 0px 0px 4px" }}>Online</span></div>
              <div className={styles.Home_status_control}>
                <div className={styles.Home_control_btn} style={{ borderColor: "#f37d79", pointerEvents: isStopping ? "none" : "auto", opacity: isStopping ? 0.5 : 1 }} onClick={handleServerStop} >
                  {isStopping ? "Stopping" : "Stop"}
                </div>
              </div>
            </div>

            <div className={styles.Home_content}>

              <div className={styles.Home_content_ping_data}>

                <div className={styles.Home_content_ping_data_item}>
                  <span className={styles.Home_content_data_item_title}>Online Players: </span><span className={styles.Home_content_data_item_info}>{serverInfo.ping_data.players.online}/{serverInfo.ping_data.players.max}</span>
                </div>

                {serverInfo.ping_data.players.sample &&
                  <div className={styles.Home_content_ping_data_item}>
                    <span className={styles.Home_content_data_item_title}>Currently Playing: </span>
                    {serverInfo.ping_data.players.sample.map(player => {
                      return <span key={player.name} className={styles.Home_content_data_item_info}>{player.name} {serverInfo.ping_data.players.sample && serverInfo.ping_data.players.sample.length > 1 && "-"}</span>
                    })}
                  </div>
                }

                <div className={styles.Home_content_ping_data_item}>
                  <span className={styles.Home_content_data_item_title}>Server Version: </span> <span className={styles.Home_content_data_item_info}>{serverInfo.ping_data.version.name}</span>
                </div>

              </div>

            </div>
          </div>
          :
          <div className={styles.Home}>
            {/* <h2>Server <span style={{ color: "rgba(255, 60, 60, 1)", textShadow: "rgba(255, 60, 60, 0.5) 0px 0px 4px" }}>Offline</span></h2>
            <div className={styles.Home_content}>
              <div className={styles.Home_content_control}>
                <div className={styles.Home_control_btn} style={{ borderColor: "#79d0bf", pointerEvents: isStarting ? "none" : "auto", opacity: isStarting ? 0.5 : 1 }} onClick={handleServerStart}>
                  {isStarting ? "Starting" : "Start"}
                </div>
              </div>
            </div> */}

            <div className={styles.Home_status}>
              <div><span className={styles.Home_status_title}>Server Status: </span><span className={styles.Home_status_server} style={{ color: "rgba(255, 60, 60, 1)", textShadow: "rgba(255, 60, 60, 0.5) 0px 0px 4px" }}>Offline</span></div>
              <div className={styles.Home_status_control}>
                <div className={styles.Home_control_btn} style={{ borderColor: "#79d0bf", pointerEvents: isStarting ? "none" : "auto", opacity: isStarting ? 0.5 : 1 }} onClick={handleServerStart}>
                  {isStarting ? "Starting" : "Start"}
                </div>
              </div>
            </div>
          </div>
        :
        <Spinner />
      }
    </SingleTab>
  )
}

export default Home


