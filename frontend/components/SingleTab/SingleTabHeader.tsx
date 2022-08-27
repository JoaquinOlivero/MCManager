import Link from 'next/link'
import { useRouter } from 'next/router'
import { useMemo } from 'react'
import styles from '../../styles/components/SingleTab/SingleTabHeader.module.scss'
import Arrow from '../../svg/icons/Arrow'

type Props = {
    tabType: string
}

const SingleTabHeader = ({ tabType }: Props) => {
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
            </div>

            <div className={styles.SingleTabHeader_title}>
                {tabType}
            </div>

            <div className={styles.SingleTabHeader_crud}>

            </div>
        </div>
    )
}

export default SingleTabHeader