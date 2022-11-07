import styles from "../../styles/components/Utils/ToggleButton.module.scss"
import { useState } from 'react'

type Props = {
    label: string
    toggled: boolean
    onClick: Function
}

export const Toggle = ({ label, toggled, onClick }: Props) => {
    const [isToggled, toggle] = useState(toggled)

    const callback = () => {
        toggle(!isToggled)
        onClick(label, !isToggled)
    }

    return (
        <div className={styles.ToggleButton}>
            <label>
                <input type="checkbox" defaultChecked={isToggled} onClick={callback} />
                <span />
                <strong>{label}</strong>
            </label>
        </div>
    )
}