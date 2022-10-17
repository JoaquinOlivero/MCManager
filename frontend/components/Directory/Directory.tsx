import { useEffect, useState } from "react";
import SingleTabDirectory from "../../components/SingleTab/components/SingleTabDirectory/SingleTabDirectory";
import SingleTab from "../../components/SingleTab/SingleTab";
import SingleTabHeader from "../../components/SingleTab/SingleTabHeader";
import { useRouter } from "next/router";
import { useDataContext } from "../../contexts/DataContext"
import ConfirmPrompt from "../Utils/ConfirmPrompt";
import NoSettingsError from "../Utils/NoSettingsError";

type Props = {
    tabType: string
}

type DirData = {
    name: string;
    type: string;
    children: [DirData] | null;
} | null;

type WorldData = {
    dir: DirData;
    world_name: string;
}

const Directory = ({ tabType }: Props) => {
    const { setEditFilepath } = useDataContext()
    const router = useRouter();
    const [settings, setSettings] = useState<boolean>(true)
    const [dirData, setDirData] = useState<DirData | null>(null);
    const [currentDir, setCurrentDir] = useState<DirData | null>(null);
    const [selectedFiles, setSelectedFiles] = useState<Array<string> | null>(null)
    const [worldName, setWorldName] = useState<string | null>(null)
    const [removePrompt, setRemovePrompt] = useState<boolean>(false)
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        if (tabType === "world") {
            getDir(tabType, setDirData, setSettings, setWorldName);
        } else {
            getDir(tabType, setDirData, setSettings);
        }

        return () => {
            setDirData(null);
        };
    }, []);

    // this use effect triggers on url path change
    useEffect(() => {
        if (dirData) {
            const asPathNestedRoutes = router.asPath.split("/").filter((v) => v.length > 0);
            if (asPathNestedRoutes.length === 1) {
                setCurrentDir(dirData);
            } else {
                const currentDir = getCurrentDir(dirData, asPathNestedRoutes);
                if (currentDir) {
                    setCurrentDir(currentDir!);
                } else {
                    // if directory is empty set error.
                    setError("Directory is empty.")
                    setCurrentDir(null)
                }
            }
        }
    }, [router.asPath]);

    // when dirData changes set new current directory.
    useEffect(() => {
        if (dirData) {
            const asPathNestedRoutes = router.asPath.split("/").filter((v) => v.length > 0);
            if (asPathNestedRoutes.length === 1) {
                setCurrentDir(dirData);
            } else {
                const currentDir = getCurrentDir(dirData, asPathNestedRoutes);
                if (currentDir) {
                    setCurrentDir(currentDir!);
                } else {
                    // if directory is empty set error.
                    setError("Directory is empty.")
                    setCurrentDir(null)
                }
            }
        }
    }, [dirData]);

    // The function handleEditFile() will set the context editFilepath state to the actual path of the file in the World directory, and then push to page "/edit".
    const handleEditFile = async () => {
        if (router.asPath.includes("world") && worldName && selectedFiles && selectedFiles.length === 1) {
            const filepath = router.asPath.replace("world", worldName) + selectedFiles[0];
            await setEditFilepath(filepath)
            return router.push("/edit")
        }
        const filepath = router.asPath + selectedFiles![0];
        await setEditFilepath(filepath)
        return router.push("/edit")
    };

    const handleRemoveFile = async () => {
        setRemovePrompt(!removePrompt)

        if (removePrompt && selectedFiles) {
            const body = { "files": selectedFiles, "directory": worldName ? router.asPath.replace("world", worldName) : router.asPath }

            const res = await fetch("/api/dir/remove", {
                method: "POST",
                body: JSON.stringify(body)
            })

            if (res.status === 200) {
                if (tabType === "world") {
                    getDir(tabType, setDirData, setSettings, setWorldName);
                } else {
                    getDir(tabType, setDirData, setSettings);
                }
            }
            return setRemovePrompt(false)
        }
    }

    // single tab layout
    return (
        <>
            <SingleTab header={<SingleTabHeader tabType={tabType} editFile={handleEditFile} removeFiles={handleRemoveFile} selectedFiles={selectedFiles} />}>
                {settings ?
                    <>
                        <SingleTabDirectory dir={currentDir} selectedFiles={selectedFiles} setSelectedFiles={setSelectedFiles} error={error} />
                        {removePrompt && <ConfirmPrompt handleConfirm={handleRemoveFile} handleCancel={() => { setRemovePrompt(false); setSelectedFiles(null) }} />}
                    </>
                    :
                    <NoSettingsError />
                }
            </SingleTab>
        </>
    );
};

export default Directory;

// Gets directory from api. All of the files and subdirectories are recursively nested inside an array. In this case, data.children.
const getDir = async (tabType: string, setDirData: (value: DirData) => void, setSettings: (value: boolean) => void, setWorldName?: (value: string) => void) => {

    fetch("/api/dir/" + tabType).then(res => {
        if (!res.ok) {
            return res.text().then(text => { throw new Error(text) })
        }
        else {
            return res.json().then((json) => {
                // for world page
                if (setWorldName) {
                    const data: WorldData = json;

                    setDirData(data.dir);
                    setWorldName(data.world_name);
                    return
                } else {
                    const data: DirData = json;
                    setDirData(data);
                }
            })
        }
    }).catch(err => {
        if (err.message.includes("no such file or directory")) {
            setSettings(false)
        }
    });
};

// Returns all files and subdirectories of a node directory.
const getCurrentDir: any = (root: DirData, nestedPaths: Array<string>) => {
    const totalRoutes = nestedPaths.length;
    const lastSlug = nestedPaths[totalRoutes - 1];
    for (const child in root!.children) {
        if (Object.prototype.hasOwnProperty.call(root!.children, child)) {
            const element = root!.children[child as any];
            if (element?.name === lastSlug) {
                return element;
            } else {
                const parent = nestedPaths.findIndex((path) => path === element?.name);
                if (nestedPaths[parent] === element?.name) {
                    return getCurrentDir(element, nestedPaths);
                }
            }
        }
    }
};