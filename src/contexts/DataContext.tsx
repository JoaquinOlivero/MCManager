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
                console.log('a')
                setPasswordExists(true)
                return true
            }
        })
            .catch(err => {
                setPasswordExists(false)
                return false
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