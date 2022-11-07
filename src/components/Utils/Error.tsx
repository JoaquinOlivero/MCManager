import styles from "../../styles/components/Utils/Error.module.scss"

type Props = {
    message: string
}

const Error = ({ message }: Props) => {
    return (
        <div className={styles.Error}>{message}</div>
    )
}

export default Error