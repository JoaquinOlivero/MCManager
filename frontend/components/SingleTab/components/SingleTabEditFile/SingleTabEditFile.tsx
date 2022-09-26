import React, { useEffect, useRef, useState } from "react";
import Spinner from '../../../../svg/icons/Spinner';
// import CodeMirror from '@uiw/react-codemirror';
import { langs } from '@uiw/codemirror-extensions-langs';
import styles from '../../../../styles/components/SingleTab/components/SingleTabEditFile/SingleTabEditFile.module.scss'
import { Extension } from '@codemirror/state';
import dynamic from 'next/dynamic'
import { Suspense } from 'react'

type Props = {
  file: string | null
  setFile: (value: string | null) => void
  fileFormat: string | null
}

// Lazy load codemirror component
const CodeMirror = dynamic(() => import('@uiw/react-codemirror'), {
  suspense: true,
})

const SingleTabEditFile = ({ file, setFile, fileFormat }: Props) => {
  const [language, setLanguage] = useState<Extension | null>(null)
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
          setLanguage(langs.textile())
          break;
        default:
          break;
      }
    }
  }, [fileFormat])

  return (
    <div className={styles.SingleTabEditFile} ref={editorContainerRef}>
      {/* {file && language && editorContainerRef.current ?
        <CodeMirror
          value={file}
          height={editorContainerRef.current!.clientHeight.toString() + "px"}
          extensions={[language]}
          onChange={onChange}
          theme="dark"
        />
        :
        <Spinner />
      } */}

      <Suspense fallback={<Spinner />}>
        {file && language && editorContainerRef.current ?
          <CodeMirror
            value={file}
            height={editorContainerRef.current!.clientHeight.toString() + "px"}
            extensions={[language]}
            onChange={onChange}
            theme="dark"
          />
          :
          <Spinner />
        }
      </Suspense>
    </div>
  )
}

export default SingleTabEditFile