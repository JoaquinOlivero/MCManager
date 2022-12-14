import React, { useCallback, useEffect, useRef, useState } from 'react'
import styles from '../../../../styles/components/SingleTab/components/SingleTabMods/SingleTabMods.module.scss'
import styleVariables from '../../../../styles/Variables.module.scss'
import JavaIcon from '../../../../svg/icons/JavaIcon'
import Spinner from '../../../../svg/icons/Spinner'
import Error from '../../../Utils/Error'

type Mod = {
    "fileName": string,
    "modId": string,
    "version": string
}

type Props = {
    mods: Array<Mod> | null,
    selectedMods: Array<string> | null
    setSelectedMods: Function
}

const SingleTabMods = ({ mods, selectedMods, setSelectedMods }: Props) => {
    const [isCtrl, setIsCtrl] = useState<boolean>(false)


    // add file info columns
    const headerItems = () => {
        const headerArr = []

        const element =
            <div className={styles.SingleTabMods_header_info} key={0}>
                <span className={styles.SingleTabMods_info_filename}>Mod filename</span>
                {/* <div className={styles.SingleTabMods_info_details}>
                    <span>Mod id</span>
                    <span>Version</span>
                </div> */}
            </div>

        if (mods!.length === 1) {
            headerArr.push(element)
        } else {
            var i = 0

            while (i < 2) {
                const element =
                    <div className={styles.SingleTabMods_header_info} key={i}>
                        <span className={styles.SingleTabMods_info_filename}>Mod filename</span>
                        {/* <div className={styles.SingleTabMods_info_details}>
                    <span>Mod id</span>
                    <span>Version</span>
                </div> */}
                    </div>
                headerArr.push(element)
                i++
            }
        }
        return headerArr
    }

    // handle click on mod.
    const selectModClick = (fileName: string) => {

        if (selectedMods) {
            const modExists = !!~selectedMods.indexOf(fileName)

            // if mod clicked already exists in the array, remove it.
            if (modExists && isCtrl) {
                // remove mod filename from the array
                const filteredMods = selectedMods.filter(m => m !== fileName)
                setSelectedMods(filteredMods)
                return
            }
            if (isCtrl) return setSelectedMods((oldArray: Array<string>) => [...oldArray, fileName])
            if (modExists && selectedMods.length === 1) return setSelectedMods(null)
        }

        // if ctrl is not pressed, only one mod is going to be selected and added to the array of selected mods.
        setSelectedMods([fileName])
    }

    // ctrl key event listener, to select multiple mods from the list.
    useEffect(() => {
        window.addEventListener("keydown", e => {
            if (e.ctrlKey && !isCtrl) setIsCtrl(true)
        })
        window.addEventListener("keyup", e => {
            if (!e.ctrlKey) setIsCtrl(false)
        })

        return () => {
            setSelectedMods(null)
        }
    }, [])

    return (
        <>
            {mods ?
                <div className={styles.SingleTabMods}>
                    {mods.length > 0 ?
                        <>
                            <div className={styles.SingleTabMods_header}>
                                {mods.length !== 0 && headerItems()}
                            </div>
                            {/* backgroundColor: styleVariables.primaryColorLowOpacity  */}
                            <div className={styles.SingleTabMods_mods_container}>
                                {mods.map((mod: Mod, i: number) => {
                                    return <div key={mod.fileName} className={styles.SingleTabMods_mod_container} style={{ borderRight: (i + 1) % 2 === 0 ? "none" : '', backgroundColor: selectedMods && selectedMods.find(m => m === mod.fileName) ? styleVariables.primaryColorLowOpacity : '' }} onClick={() => selectModClick(mod.fileName)}>
                                        <div><JavaIcon style={{ fill: selectedMods && selectedMods.find(m => m === mod.fileName) ? "white" : '', opacity: selectedMods && selectedMods.find(m => m === mod.fileName) ? 1 : 0.8 }} /></div>
                                        <div className={styles.SingleTabMods_mod_filename}>{mod.fileName}</div>
                                        {/* <div className={styles.SingleTabMods_mod_id}>{mod.version}</div> */}
                                    </div>
                                })}
                            </div>
                        </>
                        :
                        <Error message='No mods found.' />
                    }
                </div>
                :
                <Spinner />
            }
        </>
    )
}

export default SingleTabMods