import { NextPage } from "next";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import SingleTabEditFile from "../components/SingleTab/components/SingleTabEditFile/SingleTabEditFile";
import SingleTab from "../components/SingleTab/SingleTab";
import SingleTabHeader from "../components/SingleTab/SingleTabHeader";
import { useDataContext } from "../contexts/DataContext";

type UploadStatus = {
    "uploading": boolean
    "finished": boolean
    "status": boolean
}

const Edit: NextPage = () => {
    const router = useRouter()
    const { editFilepath, setEditFilepath } = useDataContext()
    const [file, setFile] = useState<string | null>(null)
    const [fileFormat, setFileFormat] = useState<string | null>(null)
    const [uploadStatus, setUploadStatus] = useState<UploadStatus>({ "uploading": false, "finished": false, "status": false })
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        if (!editFilepath) router.back()

        if (editFilepath) getFileContent(editFilepath, setFile, setFileFormat, setError)

        return () => {
            if (editFilepath === "/server.properties") setEditFilepath(null)
            setFile(null)
            setFileFormat(null)
        }

    }, [editFilepath])

    const saveFile = async () => {
        setUploadStatus({ "uploading": true, "finished": false, "status": false })

        const body = { "filepath": editFilepath, "file_content": file }
        const res = await fetch("/api/edit/save", {
            method: "POST",
            body: JSON.stringify(body),
        })

        if (res.status === 200) {
            setUploadStatus({ "uploading": false, "finished": true, "status": true })
            setTimeout(() => {
                setUploadStatus({ "uploading": false, "finished": false, "status": false })
            }, 2500);
            return
        } else {
            setUploadStatus({ "uploading": false, "finished": true, "status": false })
            setTimeout(() => {
                setUploadStatus({ "uploading": false, "finished": false, "status": false })
            }, 2500);
            return
        }
    }

    return (
        <>
            {editFilepath &&
                <SingleTab header={<SingleTabHeader tabType="edit" selectedFiles={editFilepath.split("/")} saveFile={saveFile} uploadStatus={uploadStatus} />}>
                    <SingleTabEditFile file={file} setFile={setFile} fileFormat={fileFormat} error={error} />
                </SingleTab>
            }
        </>
    )
}

export default Edit

// Gets file content from api.
const getFileContent = async (editFilepath: string, setFile: (value: string | null) => void, setFileFormat: (value: string | null) => void, setError: (value: string | null) => void) => {

    fetch(`/api/edit?filepath=${editFilepath}`)
        .then(res => {
            if (!res.ok) {
                return res.text().then(text => { throw new Error(text) })
            }
            else {
                return res.json().then(data => {
                    setFile(data.file_content)
                    setFileFormat(data.file_format)
                })
            }
        })
        .catch(err => {
            setError(err.message)
        });

}