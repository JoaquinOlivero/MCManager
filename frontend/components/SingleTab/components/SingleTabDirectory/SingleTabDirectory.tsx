import { useEffect, useState } from "react";
import styles from "../../../../styles/components/SingleTab/components/SingleTabDirectory/SingleTabDirectory.module.scss";
import variables from "../../../../styles/Variables.module.scss"
import File from "../../../../svg/icons/File";
import Folder from "../../../../svg/icons/Folder";
import Spinner from "../../../../svg/icons/Spinner";
import { useRouter } from 'next/router'

type ConfigData = {
  name: string;
  type: string;
  children: [ConfigData] | null;
} | null;

type Props = {
  dir: ConfigData;
};

const SingleTabDirectory = ({ dir }: Props) => {
  const router = useRouter()
  const [sortedData, setSortedData] = useState<Array<ConfigData> | null>(null);

  // add file info columns
  const headerItems = () => {
    const headerArr = [];

    const element = (
      <div className={styles.SingleTabDirectory_header_info} key={0}>
        <span>Name</span>
      </div>
    );

    if (dir!.children!.length === 1) {
      headerArr.push(element);
    } else {
      var i = 0;

      while (i < 2) {
        const element = (
          <div className={styles.SingleTabDirectory_header_info} key={i}>
            <span>Name</span>
          </div>
        );
        headerArr.push(element);
        i++;
      }
    }
    return headerArr;
  };

  useEffect(() => {
    if (dir) {
      const configSortedData = dir.children?.sort((a: any, b: any) => a.name.localeCompare(b.name));
      setSortedData(configSortedData!);
    }
    return () => { };
  }, [sortedData, dir]);

  const handleClickOnDir = (name: string) => {
    // const asPathNestedRoutes = router.asPath.split("/").filter(v => v.length > 0);
    router.push(`${router.asPath}${name}`, undefined, { shallow: true })
  }

  return (
    <>
      {dir ? (
        <div className={styles.SingleTabDirectory}>
          <div className={styles.SingleTabDirectory_header}>{sortedData && headerItems()}</div>
          <div className={styles.SingleTabDirectory_dir_container}>
            {sortedData &&
              sortedData.map((child) => {
                return (
                  <div key={child?.name} className={styles.SingleTabDirectory_dir_content} onClick={child?.type === "dir" ? () => { handleClickOnDir(child.name) } : undefined}>
                    {
                      child?.type === "file" ? <File fill="white" /> : <Folder fill={variables.primaryColor} />
                    }
                    <div className={styles.SingleTabDirectory_dir_name}>
                      {child?.name}
                    </div>
                  </div>
                );
              })}
          </div>
        </div>
      ) : (
        <Spinner />
      )}
    </>
  );
};

export default SingleTabDirectory;
