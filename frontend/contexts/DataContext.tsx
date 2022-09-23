import { createContext, ReactNode, useState, useContext } from "react";

// type Mod = {
//     "fileName": string,
//     "modId": string,
//     "version": string
// }


type dataContextType = {
    // mods: Array<Mod> | null;
    editFilepath: string | null
    setEditFilepath: (value: string | null) => void;
}

const dataContextDefaultValue: dataContextType = {
    // mods: null,
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
    const [editFilepath, setEditFilepath] = useState<string | null>(null)
    // const [mods, setMods] = useState<Array<Mod> | null>(null)



    const value = {
        // mods
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