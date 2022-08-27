import styles from "../../styles/components/SingleTab/SingleTab.module.scss"


type Props = {
    header: React.ReactNode
    children: React.ReactNode

}
const SingleTab = ({ header, children }: Props) => {
    return (
        <div className={styles.SingleTab}>
            <div className={styles.SingleTab_header}>
                {header}
            </div>
            <div className={styles.SingleTab_content}>
                {children}
            </div>
        </div>
    )
}

export default SingleTab