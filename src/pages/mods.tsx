import type { NextPage } from 'next'
import { ChangeEvent, useEffect, useState } from 'react'
import SingleTabMods from '../components/SingleTab/components/SingleTabMods/SingleTabMods'
import SingleTab from '../components/SingleTab/SingleTab'
import SingleTabHeader from '../components/SingleTab/SingleTabHeader'
import ConfirmPrompt from '../components/Utils/ConfirmPrompt'

type Mod = {
    "fileName": string,
    "modId": string,
    "version": string
}

type UploadStatus = {
    "uploading": boolean
    "finished": boolean
    "status": boolean
}

const Mods: NextPage = () => {
    const [mods, setMods] = useState<Array<Mod> | null>(null)
    const [selectedMods, setSelectedMods] = useState<Array<string> | null>(null)
    const [uploadStatus, setUploadStatus] = useState<UploadStatus>({ "uploading": false, "finished": false, "status": false })
    const [removePrompt, setRemovePrompt] = useState<boolean>(false)

    const uploadMods = async (e: ChangeEvent<HTMLInputElement>) => {
        setUploadStatus({ "uploading": true, "finished": false, "status": false })
        const files = e.target.files
        const formData = new FormData();

        if (files) {
            let i = 0
            while (i < files.length) {
                formData.append("mods", files[i])
                i++
            }
        }

        const res = await fetch("/api/mods/upload", {
            method: "POST",
            body: formData
        })

        if (res.status === 200) {
            getMods(setMods)
            setUploadStatus({ "uploading": false, "finished": true, "status": true })
            setTimeout(() => {
                setUploadStatus({ "uploading": false, "finished": false, "status": false })
            }, 2500);
            return
        }
        getMods(setMods)
        setUploadStatus({ "uploading": false, "finished": true, "status": false })
        return
    }

    const deleteMods = async () => {
        setRemovePrompt(true)
        if (removePrompt && selectedMods) {
            const body = { "mods": selectedMods }
            if (selectedMods) {
                const res = await fetch("/api/mods/remove", {
                    method: "POST",
                    body: JSON.stringify(body)
                })

                if (res.status === 200) {
                    getMods(setMods)
                    setSelectedMods(null)
                    return setRemovePrompt(false)
                }
                return
            }
        }

        return
    }

    useEffect(() => {
        getMods(setMods)


        return () => {
            setMods(null)
        }
    }, [])

    return (
        // single tab layout
        <SingleTab header={<SingleTabHeader tabType={"mods"} selectedFiles={selectedMods} removeFiles={deleteMods} uploadFiles={uploadMods} uploadStatus={uploadStatus} />}>
            <SingleTabMods mods={mods} selectedMods={selectedMods} setSelectedMods={setSelectedMods} />
            {removePrompt && <ConfirmPrompt handleConfirm={deleteMods} handleCancel={() => { setRemovePrompt(false); setSelectedMods(null) }} />}
        </SingleTab>
    )
}

export default Mods

const getMods = (setMods: Function) => {

    fetch("/api/mods", {
        method: "GET",
        credentials: "include"
    }).then(res => {
        if (!res.ok) {
            return res.text().then(text => { throw new Error(text); })
        }
        else {
            if (res.status === 204) {
                return setMods([])
            }

            if (res.status === 200) {
                return res.json().then((json) => {
                    return setMods(json)
                })
            }
        }
    })
        .catch(err => {
            console.log(err.message)
        })
}