import styles from '../../styles/components/Utils/ConfirmPrompt.module.scss'

type Props = {
    handleConfirm: () => Promise<void>;
    handleCancel: () => void;
}

const ConfirmPrompt = ({ handleConfirm, handleCancel }: Props) => {
    return (
        <>
            <div className={styles.ConfirmPrompt_backdrop}></div>
            <div className={styles.ConfirmPrompt}>
                <div className={styles.ConfirmPrompt_content}>
                    <h3>Do you want to proceed?</h3>
                    <div className={styles.ConfirmPrompt_content_btns}>
                        <span className={styles.ConfirmPrompt_btn_confirm} onClick={handleConfirm}>Confirm</span>
                        <span className={styles.ConfirmPrompt_btn_cancel} onClick={handleCancel}>Cancel</span>
                    </div>
                </div>
            </div>
        </>
    )
}

export default ConfirmPrompt