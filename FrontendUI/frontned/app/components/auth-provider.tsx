// // components/auth-provider.tsx
// "use client"
//
// import { WalletManager } from "../lib/wallet"
// import { createContext, useContext, useEffect, useState } from "react"
//
// type AuthContextType = {
//     initializeWallet: (password: string) => Promise<void>
//     signMessage: (message: string) => Promise<string>
//     logout: () => void
//     isAuthenticated: boolean
// }
//
// const AuthContext = createContext<AuthContextType | null>(null)
//
// export function AuthProvider({ children }: { children: React.ReactNode }) {
//     const [isAuthenticated, setIsAuthenticated] = useState(false)
//
//     const initializeWallet = async (password: string) => {
//         const walletManager = WalletManager.getInstance()
//         await walletManager.initializeWallet(password)
//         setIsAuthenticated(true)
//     }
//
//     const signMessage = async (message: string) => {
//         const walletManager = WalletManager.getInstance()
//         return walletManager.signMessage(message)
//     }
//
//     const logout = () => {
//         const walletManager = WalletManager.getInstance()
//         walletManager.clearSession()
//         setIsAuthenticated(false)
//     }
//
//     return (
//         <AuthContext.Provider
//             value={{ initializeWallet, signMessage, logout, isAuthenticated }}
//         >
//             {children}
//         </AuthContext.Provider>
//     )
// }
//
// export const useAuth = () => {
//     const context = useContext(AuthContext)
//     if (!context) {
//         throw new Error("useAuth must be used within an AuthProvider")
//     }
//     return context
// }