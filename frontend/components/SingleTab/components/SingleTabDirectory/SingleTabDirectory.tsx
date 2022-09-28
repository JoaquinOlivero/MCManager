import { useEffect, useState } from "react";
import styles from "../../../../styles/components/SingleTab/components/SingleTabDirectory/SingleTabDirectory.module.scss";
import styleVariables from "../../../../styles/Variables.module.scss";
import File from "../../../../svg/icons/File";
import Folder from "../../../../svg/icons/Folder";
import Spinner from "../../../../svg/icons/Spinner";
import { useRouter } from "next/router";

type ConfigData = {
  name: string;
  type: string;
  children: [ConfigData] | null;
} | null;

type Props = {
  dir: ConfigData;
  selectedFiles: Array<string> | null
  setSelectedFiles: Function;
};

const SingleTabDirectory = ({ dir, selectedFiles, setSelectedFiles }: Props) => {
  const router = useRouter();
  const [sortedData, setSortedData] = useState<Array<ConfigData> | null>(null);
  const [isCtrl, setIsCtrl] = useState<boolean>(false)

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
      const dirSortedData = dir.children?.sort((a: any, b: any) => a.name.localeCompare(b.name));
      setSortedData(dirSortedData!);
    }
    return () => {
      setSelectedFiles(null);
    };
  }, [sortedData, dir]);

  const handleDoubleClickOnDir = async (name: string) => {
    router.push(`${router.asPath}${name}`, undefined, { shallow: true });
  };

  // handle click on file.
  const selectFileClick = (fileName: string) => {

    if (selectedFiles) {
      const fileExists = !!~selectedFiles.indexOf(fileName)

      // if file clicked already exists in the array, remove it.
      if (fileExists && isCtrl) {
        // remove file filename from the array
        const filteredFiles = selectedFiles.filter(m => m !== fileName)
        setSelectedFiles(filteredFiles)
        return
      }
      if (isCtrl) return setSelectedFiles((oldArray: Array<string>) => [...oldArray, fileName])
      if (fileExists && selectedFiles.length === 1) return setSelectedFiles(null)
    }

    // if ctrl is not pressed, only one file is going to be selected and added to the array of selected files.
    setSelectedFiles([fileName])
  }

  // ctrl key event listener, to select multiple files from the list.
  useEffect(() => {
    window.addEventListener("keydown", e => {
      if (e.ctrlKey && !isCtrl) setIsCtrl(true)
    })
    window.addEventListener("keyup", e => {
      if (!e.ctrlKey) setIsCtrl(false)
    })

    return () => {
      setSelectedFiles(null)
    }
  }, [])

  return (
    <>
      {dir ? (
        <div className={styles.SingleTabDirectory}>
          <div className={styles.SingleTabDirectory_header}>{sortedData && headerItems()}</div>
          <div className={styles.SingleTabDirectory_dir_container}>
            {sortedData &&
              sortedData.map((child, i: number) => {
                return (
                  <div key={child?.name} style={{ borderRight: (i + 1) % 2 === 0 ? "none" : '', backgroundColor: selectedFiles && selectedFiles.find(m => m === child?.name) ? styleVariables.primaryColorLowOpacity : '' }} className={styles.SingleTabDirectory_dir_content} onClick={() => selectFileClick(child!.name)} onDoubleClick={child?.type === "dir" ? () => { handleDoubleClickOnDir(child.name) } : undefined}>
                    {
                      child?.type === "file" ? <File fill="white" /> : <Folder fill={styleVariables.primaryColor} />
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
