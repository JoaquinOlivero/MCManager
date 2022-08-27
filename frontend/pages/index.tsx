import type { NextPage } from 'next'
import styles from '../styles/Home.module.scss'

const Home: NextPage = () => {



  return (
    <div className={styles.Home_container}>
      {/* {mods && <SingleTab mods={mods} />} */}
      {/* <SingleTab /> */}

      HOME PAGE
    </div>
  )
}

export default Home


// const getMods = async (setMods: Function) => {
//   const res = await fetch("/api/mods")
//   const data = await res.json()
//   await setMods(data.mods)
// }
