import Link from 'next/link'
import { useRouter } from 'next/router'
import { ChangeEvent, useMemo, useState } from 'react'
import styles from '../../styles/components/SingleTab/SingleTabHeader.module.scss'
import AddFile from '../../svg/icons/AddFile'
import Arrow from '../../svg/icons/Arrow'
import Edit from '../../svg/icons/Edit'
import Trash from '../../svg/icons/Trash'

type UploadStatus = {
    "uploading": boolean
    "finished": boolean
    "status": boolean
}

type Props = {
    tabType: string
    selectedFiles?: Array<string> | null
    removeFiles?: Function
    uploadFiles?: Function
    uploadStatus?: UploadStatus
}

const SingleTabHeader = ({ tabType, selectedFiles, removeFiles, uploadFiles, uploadStatus }: Props) => {

    const router = useRouter()

    const breadcrumbs = useMemo(function generateBreadcrumbs() {
        const asPathWithoutQuery = router.asPath.split("?")[0];
        const asPathNestedRoutes = asPathWithoutQuery.split("/")
            .filter(v => v.length > 0);

        const crumblist = asPathNestedRoutes.map((subpath, idx) => {
            const href = "/" + asPathNestedRoutes.slice(0, idx + 1).join("/");
            return { href, text: subpath };
        })

        return [{ href: "/", text: "Home" }, ...crumblist];
    }, [router.asPath]);

    const goToPreviousCrumb = (breadcrumbs: any) => {
        const pageToGo = breadcrumbs[breadcrumbs.length - 2].href
        router.push(pageToGo)
    }

    return (
        <div className={styles.SingleTabHeader}>
            <div className={styles.SingleTabHeader_left_side}>
                {tabType !== "home" &&
                    <>
                        <Arrow onClick={() => goToPreviousCrumb(breadcrumbs)} />
                        <div className={styles.SingleTabHeader_breadcrumbs}>
                            {/* <Link href='/'>Home</Link> */}
                            {breadcrumbs.map((crumb: any, i: number) => {
                                return <div key={crumb.text} className={styles.SingleTabHeader_breadcrumbs_content}>
                                    <Link href={crumb.href}><div className={styles.SingleTabHeader_breadcrumbs_href}>{crumb.text}</div></Link>
                                    <span className={styles.SingleTabHeader_breadcrumbs_separator}>{breadcrumbs.length - 1 !== i && "/"}</span>
                                </div>
                            })}
                        </div>
                    </>
                }
            </div>

            <div className={styles.SingleTabHeader_title}>
                {tabType}
            </div>

            <div className={styles.SingleTabHeader_crud}>
                <div className={styles.SingleTabHeader_crud_msg}>
                    {uploadStatus?.uploading && !uploadStatus?.finished && <div>Uploading</div>}
                    {uploadStatus?.finished && <div>{uploadStatus?.status ? "Mods uploaded Successfully!" : "Upload Failed"}</div>}
                </div>
                {tabType === "mods" && uploadFiles && <div className={`${styles.SingleTabHeader_crud_add} ${styles.SingleTabHeader_crud_btn}`}>
                    <input type="file" multiple style={{ display: "none" }} id="file-input" accept='.jar' onChange={(e) => uploadFiles(e)} />
                    <label htmlFor='file-input'>
                        <AddFile />
                        <span>Add</span>
                    </label>
                </div>}

                {tabType === "config" && <div className={`${styles.SingleTabHeader_crud_edit} ${styles.SingleTabHeader_crud_btn}`} style={{ opacity: selectedFiles && selectedFiles.length === 1 ? 1 : 0.5, cursor: selectedFiles ? "pointer" : "default" }}>
                    <Edit />
                    <span>Edit</span>
                </div>
                }
                {removeFiles && <div className={`${styles.SingleTabHeader_crud_remove} ${styles.SingleTabHeader_crud_btn}`} style={{ opacity: selectedFiles ? 1 : 0.5, cursor: selectedFiles ? "pointer" : "default", pointerEvents: selectedFiles ? 'visible' : "none" }} onClick={() => removeFiles()}>
                    <Trash />
                    <span>Remove</span>
                </div>
                }

            </div>
        </div>
    )
}

export default SingleTabHeader