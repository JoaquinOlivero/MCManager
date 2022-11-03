import { createContext, ReactNode, useState, useContext } from "react";

// type Mod = {
//     "fileName": string,
//     "modId": string,
//     "version": string
// }


type dataContextType = {
    // mods: Array<Mod> | null;
    signedIn: boolean | null;
    setSignedIn: (value: boolean) => void;
    checkSession: () => void;
    editFilepath: string | null
    setEditFilepath: (value: string | null) => void;
}

const dataContextDefaultValue: dataContextType = {
    // mods: null,
    signedIn: null,
    setSignedIn: () => { },
    checkSession: () => { },
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

    const value = {
        // mods
        signedIn,
        setSignedIn,
        checkSession,
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