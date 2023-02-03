import { createContext, ReactNode, useState, useContext } from "react";

// type Mod = {
//     "fileName": string,
//     "modId": string,
//     "version": string
// }


type dataContextType = {
    // mods: Array<Mod> | null;
    signedIn: boolean | null;
    passwordExists: boolean | null;
    setSignedIn: (value: boolean) => void;
    checkSession: () => void;
    checkPasswordExists: () => void;
    checkSettings: () => void;
    completeSettings: boolean | null;
    setCompleteSettings: (value: boolean) => void;
    editFilepath: string | null
    setEditFilepath: (value: string | null) => void;
}

const dataContextDefaultValue: dataContextType = {
    // mods: null,
    signedIn: null,
    passwordExists: null,
    setSignedIn: () => { },
    checkSession: () => { },
    checkPasswordExists: () => { },
    checkSettings: () => { },
    completeSettings: null,
    setCompleteSettings: () => { },
    editFilepath: null,
    setEditFilepath: () => { },
}

export function useDataContext() {
    return useContext(DataContext);
}

type Props = {
    children: ReactNode;
};

const DataContext = createContext<dataContextType>(dataContextDefaultValue)

export function DataProvider({ children }: Props) {
    const [passwordExists, setPasswordExists] = useState<boolean | null>(null)
    const [signedIn, setSignedIn] = useState<boolean | null>(null)
    const [editFilepath, setEditFilepath] = useState<string | null>(null)
    const [completeSettings, setCompleteSettings] = useState<boolean | null>(null)
    // const [mods, setMods] = useState<Array<Mod> | null>(null)

    const checkSession = () => {
        fetch("/api/check", {
            credentials: "include",
            method: "GET",
        }).then(res => {
            if (!res.ok) {
                setSignedIn(false)
                return false
            }
            else {
                setSignedIn(true)
                return true
            }
        })
            .catch(err => {
                setSignedIn(false)
                return false
            });
    }

    const checkPasswordExists = () => {
        fetch("/api/password/check", {
            credentials: "include",
            method: "GET",
        }).then(res => {
            if (!res.ok) {
                setPasswordExists(false)
                return false
            }
            else {
                setPasswordExists(true)
                return true
            }
        })
            .catch(err => {
                setPasswordExists(false)
                return false
            })
    }

    const checkSettings = () => {
        fetch("/api/settings/check", {
            method: "GET",
            credentials: "include"
        }).then(res => {
            if (!res.ok) return

            if (res.status === 204) return setCompleteSettings(false)

            return setCompleteSettings(true)
        })
    }

    const value = {
        // mods
        signedIn,
        setSignedIn,
        passwordExists,
        setPasswordExists,
        checkSession,
        checkPasswordExists,
        checkSettings,
        completeSettings,
        setCompleteSettings,
        editFilepath,
        setEditFilepath
    }
    return (
        <>
            <DataContext.Provider value={value}>
                {children}
            </DataContext.Provider>
        </>
    );
}