import React, { useEffect, useRef } from "react";
import Spinner from '../../../../svg/icons/Spinner';
import CodeMirror, { ReactCodeMirrorRef } from '@uiw/react-codemirror';
import { langs } from '@uiw/codemirror-extensions-langs';
import styles from '../../../../styles/components/SingleTab/components/SingleTabEditFile/SingleTabEditFile.module.scss'
import { EditorView } from '@codemirror/view';

type Props = {
  file: string | null
  setFile: (value: string | null) => void
}

const SingleTabEditFile = ({ file, setFile }: Props) => {
  const editorContainerRef = useRef<HTMLDivElement>(null)
  const codeMirrorRef = useRef<ReactCodeMirrorRef>(null)
  const onChange = React.useCallback((value: string) => {
    setFile(value)
  }, []);

  const handleOnCreateEditor = (view: EditorView) => {
    const scrollDom = view.scrollDOM
    // scrollDom.setAttributeNS(null, "className", styles.SingleTabEditFile_scrollbar)
  }

  useEffect(() => {
    // console.log(codeMirrorRef.current);
  }, [codeMirrorRef.current])


  return (
    <div className={styles.SingleTabEditFile} ref={editorContainerRef}>
      {file && editorContainerRef.current ?
        <CodeMirror
          value={file}
          height={editorContainerRef.current!.clientHeight.toString() + "px"}
          extensions={[langs.toml()]}
          onChange={onChange}
          theme="dark"
          ref={codeMirrorRef}
          onCreateEditor={(view) => handleOnCreateEditor(view)}
        />
        :
        <Spinner />
      }
    </div>
  )
}

export default SingleTabEditFile