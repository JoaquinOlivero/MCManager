import { createContext, ReactNode, useState } from "react";

type Mod = {
    "fileName": string,
    "modId": string,
    "version": string
}


type dataContextType = {
    mods: Array<Mod> | null;
}

const dataContextDefaultValue: dataContextType = {
    mods: null,
}


type Props = {
    children: ReactNode;
};

const DataContext = createContext<dataContextType>(dataContextDefaultValue)

export function DataProvider({ children }: Props) {
    const [mods, setMods] = useState<Array<Mod> | null>(null)



    const value = {
        mods
    }
    return (
        <>
            <DataContext.Provider value={value}>
                {children}
            </DataContext.Provider>
        </>
    );
}