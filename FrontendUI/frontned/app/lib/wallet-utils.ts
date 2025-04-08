// import { ethers } from "ethers"
// import CryptoJS from "crypto-js"
//
// /**
//  * Generate a new Ethereum wallet with mnemonic
//  * @returns {Object} Wallet information including address, private key, and mnemonic
//  */
// export const generateEthereumWallet = () => {
//     try {
//         // Create a random wallet with mnemonic
//         const wallet = ethers.Wallet.createRandom()
//
//         return {
//             address: wallet.address,
//             privateKey: wallet.privateKey,
//             mnemonic: wallet.mnemonic.phrase,
//             success: true,
//         }
//     } catch (error) {
//         console.error("Error generating wallet:", error)
//         return {
//             address: "",
//             privateKey: "",
//             mnemonic: "",
//             success: false,
//             error: error instanceof Error ? error.message : "Unknown error",
//         }
//     }
// }
//
// /**
//  * Import an Ethereum wallet from a mnemonic phrase
//  * @param {string} mnemonic - The mnemonic phrase (12 or 24 words)
//  * @returns {Object} Wallet information including address and private key
//  */
// export const importWalletFromMnemonic = (mnemonic: string) => {
//     try {
//         // Validate mnemonic
//         if (!ethers.utils.isValidMnemonic(mnemonic.trim())) {
//             throw new Error("Invalid mnemonic phrase")
//         }
//
//         // Create wallet from mnemonic
//         const wallet = ethers.Wallet.fromMnemonic(mnemonic.trim())
//
//         return {
//             address: wallet.address,
//             privateKey: wallet.privateKey,
//             success: true,
//         }
//     } catch (error) {
//         console.error("Error importing wallet:", error)
//         return {
//             address: "",
//             privateKey: "",
//             success: false,
//             error: error instanceof Error ? error.message : "Unknown error",
//         }
//     }
// }
//
// /**
//  * Encrypt a mnemonic phrase with a password
//  * @param {string} mnemonic - The mnemonic phrase to encrypt
//  * @param {string} password - The password to encrypt with
//  * @returns {string} The encrypted mnemonic
//  */
// export const encryptMnemonic = (mnemonic: string, password: string) => {
//     try {
//         return CryptoJS.AES.encrypt(mnemonic, password).toString()
//     } catch (error) {
//         console.error("Error encrypting mnemonic:", error)
//         throw error
//     }
// }
//
// /**
//  * Decrypt an encrypted mnemonic phrase with a password
//  * @param {string} encryptedMnemonic - The encrypted mnemonic phrase
//  * @param {string} password - The password to decrypt with
//  * @returns {string} The decrypted mnemonic
//  */
// export const decryptMnemonic = (encryptedMnemonic: string, password: string) => {
//     try {
//         const bytes = CryptoJS.AES.decrypt(encryptedMnemonic, password)
//         const decryptedMnemonic = bytes.toString(CryptoJS.enc.Utf8)
//
//         if (!decryptedMnemonic) {
//             throw new Error("Incorrect password")
//         }
//
//         return decryptedMnemonic
//     } catch (error) {
//         console.error("Error decrypting mnemonic:", error)
//         throw error
//     }
// }
//
// /**
//  * Sign a message with a private key
//  * @param {string} message - The message to sign
//  * @param {string} privateKey - The private key to sign with
//  * @returns {Object} The signature and other information
//  */
// export const signMessage = async (message: string, privateKey: string) => {
//     try {
//         const wallet = new ethers.Wallet(privateKey)
//         const messageHash = ethers.utils.hashMessage(message)
//         const signature = await wallet.signMessage(message)
//
//         return {
//             messageHash,
//             signature,
//             signer: wallet.address,
//             success: true,
//         }
//     } catch (error) {
//         console.error("Error signing message:", error)
//         return {
//             messageHash: "",
//             signature: "",
//             signer: "",
//             success: false,
//             error: error instanceof Error ? error.message : "Unknown error",
//         }
//     }
// }
//
// /**
//  * Verify a message signature
//  * @param {string} message - The original message
//  * @param {string} signature - The signature to verify
//  * @returns {Object} Verification result including the recovered address
//  */
// export const verifySignature = (message: string, signature: string) => {
//     try {
//         // Recover the address of the signer
//         const recoveredAddress = ethers.utils.verifyMessage(message, signature)
//
//         return {
//             recoveredAddress,
//             success: true,
//         }
//     } catch (error) {
//         console.error("Error verifying signature:", error)
//         return {
//             recoveredAddress: "",
//             success: false,
//             error: error instanceof Error ? error.message : "Unknown error",
//         }
//     }
// }
//
// /**
//  * Calculate password strength score (0-100)
//  * @param {string} password - The password to evaluate
//  * @returns {number} Strength score from 0-100
//  */
// export const calculatePasswordStrength = (password: string) => {
//     let strength = 0
//
//     if (password.length >= 8) strength += 20
//     if (password.length >= 12) strength += 10
//     if (/[A-Z]/.test(password)) strength += 20
//     if (/[a-z]/.test(password)) strength += 10
//     if (/[0-9]/.test(password)) strength += 20
//     if (/[^A-Za-z0-9]/.test(password)) strength += 20
//
//     return strength
// }
