import React, { useEffect, useRef, useState } from "react";
import Spinner from '../../../../svg/icons/Spinner';
import { langs } from '@uiw/codemirror-extensions-langs';
import styles from '../../../../styles/components/SingleTab/components/SingleTabEditFile/SingleTabEditFile.module.scss'
import { Extension } from '@codemirror/state';
import dynamic from 'next/dynamic'
import { Suspense } from 'react'
import Error from "../../../Utils/Error";
import NoSettingsError from "../../../Utils/NoSettingsError";

type Props = {
  file: string | null
  setFile: (value: string | null) => void
  fileFormat: string | null
  error: string | null
}

// Lazy load codemirror component
const CodeMirror = dynamic(() => import('@uiw/react-codemirror'), {
  suspense: true,
  ssr: false,
})

const SingleTabEditFile = ({ file, setFile, fileFormat, error }: Props) => {
  const [language, setLanguage] = useState<Extension | null>(null)
  const [editable, setEditable] = useState<boolean>(true)
  const editorContainerRef = useRef<HTMLDivElement>(null)
  const onChange = React.useCallback((value: string) => {
    setFile(value)
  }, []);

  useEffect(() => {
    if (fileFormat) {
      switch (fileFormat) {
        case ".json":
        case ".json5":
          setLanguage(langs.json())
          break;
        case ".toml":
          setLanguage(langs.toml())
          break;
        case ".properties":
          setLanguage(langs.properties())
          break;
        case ".cfg":
        case ".txt":
        case ".log":
          setEditable(false)
          setLanguage(langs.textile())
          break;
        default:
          break;
      }
    }
  }, [fileFormat])

  return (
    <div className={styles.SingleTabEditFile} ref={editorContainerRef}>
      <Suspense fallback={<Spinner />}>
        {file && language && editorContainerRef.current ?
          <CodeMirror
            value={file}
            height={editorContainerRef.current!.clientHeight.toString() + "px"}
            extensions={[language]}
            onChange={onChange}
            theme="dark"
            editable={editable}
          />
          :
          <>
            {error ?
              <>
                {error.includes("no such file or directory") ?
                  <NoSettingsError />
                  :
                  <Error message={error} />
                }
              </>
              :
              <Spinner />
            }
          </>
        }
      </Suspense>
    </div>
  )
}

export default SingleTabEditFile