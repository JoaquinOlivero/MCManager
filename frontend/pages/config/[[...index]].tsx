import { NextPage } from "next";
import { useEffect, useState } from "react";
import SingleTabDirectory from "../../components/SingleTab/components/SingleTabDirectory/SingleTabDirectory";
import SingleTab from "../../components/SingleTab/SingleTab";
import SingleTabHeader from "../../components/SingleTab/SingleTabHeader";
import { useRouter } from "next/router";

type ConfigData = {
  name: string;
  type: string;
  children: [ConfigData] | null;
} | null;

const Config: NextPage = () => {
  const router = useRouter();
  const [configData, setConfigData] = useState<ConfigData | null>(null);
  const [currentConfigDir, setCurrentConfigDir] = useState<ConfigData | null>(null);
  const [selectedFile, setSelectedFile] = useState<string | null>(null);

  useEffect(() => {
    const asPathNestedRoutes = router.asPath.split("/").filter((v) => v.length > 0);
    getConfigDir(setConfigData, setCurrentConfigDir, asPathNestedRoutes);

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
        setCurrentConfigDir(currentDir!);
      }
    }

    return () => {};
  }, [router.asPath]);

  useEffect(() => {
    if (configData) {
      const asPathNestedRoutes = router.asPath.split("/").filter((v) => v.length > 0);
      if (asPathNestedRoutes.length === 1) {
        setCurrentConfigDir(configData);
      } else {
        const currentDir = getCurrentDir(configData, asPathNestedRoutes);
        setCurrentConfigDir(currentDir!);
      }
    }
  }, [configData]);

  // TODO: create edit file page and use https://github.com/uiwjs/react-textarea-code-editor <--- this library to edit files such as .json, .toml and .txt.
  // The function handleEditFile() will transition to page /file/:id. The :id url parameter should be the path to the file in the config folder. For example: /config/byg/byg-biome-dictionary.json.
  const handleEditFile = () => {
    const filepath = router.asPath + selectedFile;
  };

  return (
    // single tab layout
    <SingleTab header={<SingleTabHeader tabType={"config"} editFile={handleEditFile} selectedFiles={[selectedFile!]} />}>
      <SingleTabDirectory dir={currentConfigDir} selectedFile={selectedFile} setSelectedFile={setSelectedFile} />
    </SingleTab>
  );
};

export default Config;

// Gets config directory from api. All of the files and subdirectories are recursively nested inside an array. In this case, data.children.
const getConfigDir = async (setConfigData: Function, setCurrentConfigDir: Function, asPathNestedRoutes: Array<string>) => {
  const res = await fetch("/api/config");
  const data: ConfigData = await res.json();
  if (res.status === 200 && data && data.children) {
    await setConfigData(data);
  }
};

// Returns all files and subdirectories of a root node directory. Used in this component to access a subdirectory of minecraft's config root directory.
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
