import Link from "next/link";
import { useRouter } from "next/router";
import { useMemo } from "react";
import styles from "../../styles/components/SingleTab/SingleTabHeader.module.scss";
import AddFile from "../../svg/icons/AddFile";
import Arrow from "../../svg/icons/Arrow";
import Edit from "../../svg/icons/Edit";
import Save from "../../svg/icons/Save";
import Trash from "../../svg/icons/Trash";
import View from "../../svg/icons/View";
import Zip from "../../svg/icons/Zip";

type UploadStatus = {
  uploading: boolean;
  finished: boolean;
  status: boolean;
};

type Props = {
  tabType: string;
  selectedFiles?: Array<string> | null;
  removeFiles?: Function;
  uploadFiles?: Function;
  editFile?: Function;
  uploadStatus?: UploadStatus;
  saveFile?: () => Promise<void>;
  extractFile?: Function;
};

const SingleTabHeader = ({ tabType, selectedFiles, removeFiles, uploadFiles, editFile, uploadStatus, saveFile, extractFile }: Props) => {
  const router = useRouter();

  const breadcrumbs = useMemo(
    function generateBreadcrumbs() {
      const asPathWithoutQuery = router.asPath.split("?")[0];
      const asPathNestedRoutes = asPathWithoutQuery.split("/").filter((v) => v.length > 0);

      const crumblist = asPathNestedRoutes.map((subpath, idx) => {
        const href = "/" + asPathNestedRoutes.slice(0, idx + 1).join("/");
        return { href, text: subpath };
      });

      return [{ href: "/", text: "Home" }, ...crumblist];
    },
    [router.asPath]
  );

  return (
    <div className={styles.SingleTabHeader}>
      <div className={styles.SingleTabHeader_left_side}>
        {tabType !== "home" && (
          <>
            <Arrow onClick={() => router.back()} />
            <div className={styles.SingleTabHeader_breadcrumbs}>
              {/* <Link href='/'>Home</Link> */}
              {tabType !== "edit" && breadcrumbs.map((crumb: { text: string; href: string }, i: number) => {
                return (
                  <div key={crumb.text} className={styles.SingleTabHeader_breadcrumbs_content}>
                    <Link href={crumb.href}>
                      <div className={styles.SingleTabHeader_breadcrumbs_href}>{crumb.text}</div>
                    </Link>
                    <span className={styles.SingleTabHeader_breadcrumbs_separator}>{breadcrumbs.length - 1 !== i && "/"}</span>
                  </div>
                );
              })}
            </div>
          </>
        )}
      </div>

      {tabType !== "edit" ? <div className={styles.SingleTabHeader_title}>{tabType}</div> : selectedFiles && selectedFiles.length > 0 && <div className={styles.SingleTabHeader_title}>{selectedFiles[selectedFiles.length - 1]}</div>}

      <div className={styles.SingleTabHeader_crud}>
        <div className={styles.SingleTabHeader_crud_msg}>
          {uploadStatus?.uploading && !uploadStatus?.finished && <div>Uploading</div>}
          {uploadStatus?.finished && <div>{uploadStatus?.status ? "Success!" : "Failed!"}</div>}
        </div>
        {tabType === "mods" && uploadFiles && (
          <div className={`${styles.SingleTabHeader_crud_add} ${styles.SingleTabHeader_crud_btn}`}>
            <input type="file" multiple style={{ display: "none" }} id="file-input" accept=".jar" onChange={(e) => uploadFiles(e)} />
            <label htmlFor="file-input">
              <AddFile />
              <span>Add</span>
            </label>
          </div>
        )}

        {tabType === "logs" && extractFile &&
          <div className={`${styles.SingleTabHeader_crud_edit} ${styles.SingleTabHeader_crud_btn}`} style={{ opacity: selectedFiles && selectedFiles.find(file => file.endsWith(".gz")) ? 1 : 0.5, pointerEvents: selectedFiles ? "visible" : "none" }} onClick={() => extractFile()}>
            <Zip />
            <span>Extract</span>
          </div>
        }
        {editFile && (
          <div className={`${styles.SingleTabHeader_crud_edit} ${styles.SingleTabHeader_crud_btn}`} style={{ opacity: selectedFiles && selectedFiles.length === 1 && selectedFiles[0] ? 1 : 0.5, cursor: selectedFiles && selectedFiles.length === 1 && selectedFiles[0] ? "pointer" : "default", pointerEvents: selectedFiles && selectedFiles.length === 1 && selectedFiles[0] ? "visible" : "none" }} onClick={() => editFile()}>
            {tabType === "logs" ?
              <>
                <View />
                <span>View</span>
              </>
              :
              <>
                <Edit />
                <span>Edit</span>
              </>
            }
          </div>
        )}
        {removeFiles && (
          <div className={`${styles.SingleTabHeader_crud_remove} ${styles.SingleTabHeader_crud_btn}`} style={{ opacity: selectedFiles ? 1 : 0.5, cursor: selectedFiles ? "pointer" : "default", pointerEvents: selectedFiles ? "visible" : "none" }} onClick={() => removeFiles()}>
            <Trash />
            <span>Remove</span>
          </div>
        )}

        {/* SingleTabEditFile section */}
        {tabType === "edit" && saveFile &&

          <div className={`${styles.SingleTabHeader_crud_save} ${styles.SingleTabHeader_crud_btn}`} onClick={() => saveFile()}>
            <Save />
            <span>Save File</span>
          </div>
        }
      </div>
    </div>
  );
};

export default SingleTabHeader;
