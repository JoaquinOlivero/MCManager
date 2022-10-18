import styles from "../../../../styles/components/SingleTab/components/SingleTabSettings/SingleTabSettings.module.scss";
import Spinner from "../../../../svg/icons/Spinner";
import Command from "./components/Command";
import Docker from "./components/Docker";

type Settings = {
  minecraft_directory: string
  run_method: string
  docker_container_id: string
  start_command: string
}


type DockerContainer = {
  container_id: string;
  container_name: string;
};

type Props = {
  settings: Settings | null;
  dockerContainers: Array<DockerContainer> | null;
  getSettings: Function;
};

const SingleTabSettings = ({ settings, dockerContainers, getSettings }: Props) => {
  return (
    <div className={styles.SingleTabSettings}>
      {settings ? (
        <div className={styles.SingleTabSettings_option_section}>
          <div className={styles.SingleTabSettings_option_header}>
            <span>Minecraft Server Control</span>
          </div>
          <div className={styles.SingleTabSettings_option_container}>
            <Command settings={settings} getSettings={getSettings} />
            <Docker settings={settings} dockerContainers={dockerContainers} getSettings={getSettings} />
          </div>
        </div>
      ) : (
        <Spinner />
      )}
    </div>
  );
};

export default SingleTabSettings;
