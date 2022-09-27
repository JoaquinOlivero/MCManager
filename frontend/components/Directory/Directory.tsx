import { useEffect, useState } from "react";
import SingleTabDirectory from "../../components/SingleTab/components/SingleTabDirectory/SingleTabDirectory";
import SingleTab from "../../components/SingleTab/SingleTab";
import SingleTabHeader from "../../components/SingleTab/SingleTabHeader";
import { useRouter } from "next/router";
import { useDataContext } from "../../contexts/DataContext"

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
    const [dirData, setDirData] = useState<DirData | null>(null);
    const [currentDir, setCurrentDir] = useState<DirData | null>(null);
    const [selectedFile, setSelectedFile] = useState<string | null>(null);
    const [worldName, setWorldName] = useState<string | null>(null)

    useEffect(() => {
        if (tabType === "world") {
            getDir(tabType, setDirData, setWorldName);
        } else {
            getDir(tabType, setDirData);
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
                    // if directory does not exist redirect to home page.
                    router.push("/")
                }
            }
        }
        return () => { };
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
                    // if directory does not exist redirect to home page.
                    router.push("/")
                }
            }
        }
    }, [dirData]);

    // The function handleEditFile() will set the context editFilepath state to the actual path of the file in the World directory, and then push to page "/edit".
    const handleEditFile = async () => {
        if (router.asPath.includes("world") && worldName) {
            const filepath = router.asPath.replace("world", worldName) + selectedFile;
            await setEditFilepath(filepath)
            return router.push("/edit")
        }
        const filepath = router.asPath + selectedFile;
        await setEditFilepath(filepath)
        return router.push("/edit")
    };

    // single tab layout
    return (
        <>
            <SingleTab header={<SingleTabHeader tabType={tabType} editFile={handleEditFile} selectedFiles={[selectedFile!]} />}>
                <SingleTabDirectory dir={currentDir} selectedFile={selectedFile} setSelectedFile={setSelectedFile} />
            </SingleTab>
        </>
    );
};

export default Directory;

// Gets world directory from api. All of the files and subdirectories are recursively nested inside an array. In this case, data.children.
const getDir = async (tabType: string, setDirData: (value: DirData) => void, setWorldName?: (value: string) => void) => {
    const res = await fetch("/api/" + tabType);

    // for world page
    if (setWorldName) {
        const data: WorldData = await res.json();

        if (res.status === 200 && data) {
            await setDirData(data.dir);
            await setWorldName(data.world_name);
            return
        }
    } else {
        // For all other pages that have directories, e.g: /config and /logs.
        const data: DirData = await res.json();

        if (res.status === 200 && data && data.children) {
            await setDirData(data);
            return
        }
    }
};

// Returns all files and subdirectories of a root node directory. This function is used in this component to access a subdirectory of minecraft's config root directory.
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