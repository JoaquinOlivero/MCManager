import { NextPage } from "next";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import SingleTabEditFile from "../components/SingleTab/components/SingleTabEditFile/SingleTabEditFile";
import SingleTab from "../components/SingleTab/SingleTab";
import SingleTabHeader from "../components/SingleTab/SingleTabHeader";
import { useDataContext } from "../contexts/DataContext";


const Edit: NextPage = () => {
    const router = useRouter()
    const { editFilepath, setEditFilepath } = useDataContext()
    const [file, setFile] = useState<string | null>(null)

    useEffect(() => {
        if (!editFilepath) router.back()

        if (editFilepath) getFileContent(editFilepath, setFile)

        return () => {
            setEditFilepath(null)
        }

    }, [])

    return (
        <>
            {editFilepath &&
                <SingleTab header={<SingleTabHeader tabType="edit" selectedFiles={editFilepath.split("/")} />}>
                    <SingleTabEditFile file={file} setFile={setFile} />
                </SingleTab>
            }
        </>
    )
}

export default Edit

// Gets file content from api.
const getFileContent = async (editFilepath: string, setFile: (value: string | null) => void) => {
    const res = await fetch(`/api/edit?filepath=${editFilepath}`);
    const data = await res.json();
    if (res.status === 200 && data) {
        setFile(data.file_content)
    }
}