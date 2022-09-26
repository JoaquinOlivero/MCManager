import { NextPage } from "next";
import { useEffect, useState } from "react";
import SingleTabDirectory from "../../components/SingleTab/components/SingleTabDirectory/SingleTabDirectory";
import SingleTab from "../../components/SingleTab/SingleTab";
import SingleTabHeader from "../../components/SingleTab/SingleTabHeader";
import { useRouter } from "next/router";
import { useDataContext } from "../../contexts/DataContext"

type ConfigData = {
  name: string;
  type: string;
  children: [ConfigData] | null;
} | null;

const Config: NextPage = () => {
  const { setEditFilepath } = useDataContext()
  const router = useRouter();
  const [configData, setConfigData] = useState<ConfigData | null>(null);
  const [currentConfigDir, setCurrentConfigDir] = useState<ConfigData | null>(null);
  const [selectedFile, setSelectedFile] = useState<string | null>(null);

  useEffect(() => {
    getConfigDir(setConfigData);

    return () => {
      setConfigData(null);
    };
  }, []);

  // this use effect triggers on url path change
  useEffect(() => {
    if (configData) {
      const asPathNestedRoutes = router.asPath.split("/").filter((v) => v.length > 0);
      if (asPathNestedRoutes.length === 1) {
        setCurrentConfigDir(configData);
      } else {
        const currentDir = getCurrentDir(configData, asPathNestedRoutes);
        if (currentDir) {
          setCurrentConfigDir(currentDir!);
        } else {
          // if directory does not exist redirect to home page.
          router.push("/")
        }


      }
    }

    return () => { };
  }, [router.asPath]);

  // when configData changes
  useEffect(() => {
    if (configData) {
      const asPathNestedRoutes = router.asPath.split("/").filter((v) => v.length > 0);
      if (asPathNestedRoutes.length === 1) {
        setCurrentConfigDir(configData);
      } else {
        const currentDir = getCurrentDir(configData, asPathNestedRoutes);
        if (currentDir) {
          setCurrentConfigDir(currentDir!);
        } else {
          // if directory does not exist redirect to home page.
          router.push("/")
        }
      }
    }
  }, [configData]);

  // The function handleEditFile() will set the context editFilepath state to the actual path of the file in the config directory, and then push to page "/edit".
  const handleEditFile = async () => {
    const filepath = router.asPath + selectedFile;
    await setEditFilepath(filepath)
    router.push("/edit")
  };

  // single tab layout
  return (
    <>
      <SingleTab header={<SingleTabHeader tabType={"config"} editFile={handleEditFile} selectedFiles={[selectedFile!]} />}>
        <SingleTabDirectory dir={currentConfigDir} selectedFile={selectedFile} setSelectedFile={setSelectedFile} />
      </SingleTab>
    </>
  );
};

export default Config;

// Gets config directory from api. All of the files and subdirectories are recursively nested inside an array. In this case, data.children.
const getConfigDir = async (setConfigData: Function) => {
  const res = await fetch("/api/config");
  const data: ConfigData = await res.json();
  if (res.status === 200 && data && data.children) {
    await setConfigData(data);
  }
};

// Returns all files and subdirectories of a root node directory. This function is used in this component to access a subdirectory of minecraft's config root directory.
const getCurrentDir: any = (root: ConfigData, nestedPaths: Array<string>) => {
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
