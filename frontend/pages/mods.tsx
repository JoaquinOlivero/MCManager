import type { NextPage } from 'next'
import { useEffect, useState } from 'react'
import SingleTabMods from '../components/SingleTab/components/SingleTabMods/SingleTabMods'
import SingleTab from '../components/SingleTab/SingleTab'
import SingleTabHeader from '../components/SingleTab/SingleTabHeader'

type Mod = {
    "fileName": string,
    "modId": string,
    "version": string
}

const Mods: NextPage = () => {
    const [mods, setMods] = useState<Array<Mod> | null>(null)

    useEffect(() => {
        getMods(setMods)


        return () => {
            setMods(null)
        }
    }, [])

    return (
        // single tab layout
        <SingleTab header={<SingleTabHeader tabType={"mods"} />}>
            <SingleTabMods mods={mods} />
        </SingleTab>
    )
}

export default Mods


const getMods = async (setMods: Function) => {
    const res = await fetch("/api/mods")
    const data = await res.json()
    await setMods(data)
}