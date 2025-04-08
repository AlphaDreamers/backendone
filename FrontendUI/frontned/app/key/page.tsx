"use client"

import type React from "react"
import { useState } from "react"
import { ethers } from "ethers"
import { Card, CardHeader, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Label } from "@/components/ui/label"
import {
  CopyIcon,
  CheckIcon,
  ShieldIcon,
  ServerIcon,
  KeyIcon,
  EyeIcon,
  EyeOffIcon,
  AlertTriangleIcon,
  ArrowRightIcon,
  LockIcon,
  AlertCircle,
} from "lucide-react"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { Progress } from "@/components/ui/progress"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { useRouter } from "next/navigation"
import { Keypair } from "@solana/web3.js"
import * as bip39 from "bip39"

export default function Wallet() {
  const router = useRouter()
  const [mnemonic, setMnemonic] = useState<string>("")
  const [privateKey, setPrivateKey] = useState<string>("")
  const [password, setPassword] = useState<string>("")
  const [confirmPassword, setConfirmPassword] = useState<string>("")
  const [generatedMnemonic, setGeneratedMnemonic] = useState<string>("")
  const [showMnemonic, setShowMnemonic] = useState<boolean>(false)
  const [mnemonicOption, setMnemonicOption] = useState<string>("self")
  const [copied, setCopied] = useState<boolean>(false)
  const [isCreating, setIsCreating] = useState<boolean>(false)
  const [isStoring, setIsStoring] = useState<boolean>(false)
  const [passwordStrength, setPasswordStrength] = useState<number>(0)
  const [activeTab, setActiveTab] = useState<string>("create")
  const [ethereumWallet, setEthereumWallet] = useState<ethers.Wallet | null>(null)
  const [solanaWallet, setSolanaWallet] = useState<{ publicKey: string; secretKey: Uint8Array } | null>(null)
  const [walletDetails, setWalletDetails] = useState<{ ethereum: string; solana: string } | null>(null)

  const handleMnemonicChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setMnemonic(e.target.value)
  }

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value)
    calculatePasswordStrength(e.target.value)
  }

  const handleConfirmPasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setConfirmPassword(e.target.value)
  }

  const calculatePasswordStrength = (password: string) => {
    let strength = 0

    if (password.length >= 8) strength += 20
    if (password.length >= 12) strength += 10
    if (/[A-Z]/.test(password)) strength += 20
    if (/[a-z]/.test(password)) strength += 10
    if (/[0-9]/.test(password)) strength += 20
    if (/[^A-Za-z0-9]/.test(password)) strength += 20

    setPasswordStrength(strength)
  }

  const getPasswordStrengthLabel = () => {
    if (passwordStrength < 40) return { label: "Weak", color: "bg-red-500" }
    if (passwordStrength < 70) return { label: "Medium", color: "bg-yellow-500" }
    return { label: "Strong", color: "bg-green-500" }
  }

  const generatePrivateKey = () => {
    try {
      // Import Ethereum wallet
      const ethWallet = ethers.Wallet.fromPhrase(mnemonic)
      setPrivateKey(ethWallet.privateKey)

      // Import Solana wallet using the same mnemonic
      const importSolanaWallet = async () => {
        try {
          const seed = await bip39.mnemonicToSeed(mnemonic)
          const seedBuffer = Buffer.from(seed).slice(0, 32)
          const solKeypair = Keypair.fromSeed(seedBuffer)

          setWalletDetails({
            ethereum: ethWallet.address,
            solana: solKeypair.publicKey.toString(),
          })
        } catch (error) {
          console.error("Error importing Solana wallet:", error)
        }
      }

      importSolanaWallet()
    } catch (error) {
      console.error("Invalid mnemonic")
      setPrivateKey("Invalid mnemonic")
    }
  }

  const createAccount = async () => {
    if (password.length < 8) {
      alert("Password must be at least 8 characters long")
      return
    }

    if (password !== confirmPassword) {
      alert("Passwords do not match")
      return
    }

    setIsCreating(true)

    try {
      // Create Ethereum wallet
      const ethWallet = ethers.Wallet.createRandom()
      setGeneratedMnemonic(ethWallet.mnemonic?.phrase || "")
      setEthereumWallet(new ethers.Wallet(ethWallet.privateKey))

      const seed = await bip39.mnemonicToSeed(ethWallet.mnemonic?.phrase || "")
      const seedBuffer = Buffer.from(seed).slice(0, 32)
      const solKeypair = Keypair.fromSeed(seedBuffer)
      setSolanaWallet({
        publicKey: solKeypair.publicKey.toString(),
        secretKey: solKeypair.secretKey,
      })

      // Set wallet details for display
      setWalletDetails({
        ethereum: ethWallet.address,
        solana: solKeypair.publicKey.toString(),
      })

      setShowMnemonic(true)

      // Simulate encryption delay
      await new Promise((resolve) => setTimeout(resolve, 1000))

      // In a real app, you'd encrypt the wallet with the password
      // const encryptedWallet = await wallet.encrypt(password);
      // localStorage.setItem("wallet", JSON.stringify(encryptedWallet));
    } catch (error) {
      console.error("Error creating account:", error)
      alert("Failed to create wallet. Please try again.")
    } finally {
      setIsCreating(false)
    }
  }

  const copyToClipboard = () => {
    if (generatedMnemonic) {
      navigator.clipboard.writeText(generatedMnemonic)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }

  // Update the handleMnemonicOptionChange function to navigate to the dashboard
  const handleMnemonicOptionChange = async () => {
    if (!generatedMnemonic) return

    setIsStoring(true)

    try {
      if (mnemonicOption === "rely") {
        // Make the API call to store the mnemonic with the service
        const response = await fetch("http://localhost:8085/api/auth/wallet", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ mnemonic: generatedMnemonic }),
        })

        if (response.ok) {
          // Store wallet data in localStorage for the dashboard to use
          if (solanaWallet) {
            localStorage.setItem(
              "solanaWallet",
              JSON.stringify({
                publicKey: solanaWallet.publicKey,
                // Note: In a real app, you would encrypt the private key
                // We're just storing the public key for demo purposes
              }),
            )
          }
          router.push("/dashboard")
        } else {
          alert("Failed to store recovery phrase.")
        }
      } else {
        // If self-custody, just proceed
        // Store wallet data in localStorage for the dashboard to use
        if (solanaWallet) {
          localStorage.setItem(
            "solanaWallet",
            JSON.stringify({
              publicKey: solanaWallet.publicKey,
              // Note: In a real app, you would encrypt the private key
              // We're just storing the public key for demo purposes
            }),
          )
        }
        router.push("/dashboard")
      }
    } catch (error) {
      console.error("Error storing recovery phrase:", error)
      alert("Failed to store recovery phrase.")
    } finally {
      setIsStoring(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-900 to-black text-white p-4 flex flex-col items-center justify-center">
      <div className="w-full max-w-md">
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold tracking-tight text-white">Wallet Setup</h1>
          <p className="mt-2 text-sm text-gray-400">Create or import your Ethereum wallet</p>
        </div>

        <Card className="border-0 bg-gray-800/50 backdrop-blur-sm shadow-2xl">
          <CardHeader>
            <Tabs defaultValue="create" value={activeTab} onValueChange={setActiveTab} className="w-full">
              <TabsList className="grid w-full grid-cols-2 bg-gray-700/50">
                <TabsTrigger value="create" className="data-[state=active]:bg-blue-600 data-[state=active]:text-white">
                  Create New
                </TabsTrigger>
                <TabsTrigger value="import" className="data-[state=active]:bg-blue-600 data-[state=active]:text-white">
                  Import Existing
                </TabsTrigger>
              </TabsList>

              <TabsContent value="create" className="space-y-6 mt-4">
                {!showMnemonic ? (
                  <div className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="password" className="text-gray-300 flex items-center gap-2">
                        <LockIcon className="h-4 w-4 text-blue-400" />
                        Wallet Password
                      </Label>
                      <Input
                        id="password"
                        type="password"
                        placeholder="Enter secure password"
                        value={password}
                        onChange={handlePasswordChange}
                        className="bg-gray-700/50 border-gray-600 text-white placeholder:text-gray-500 focus:border-blue-500 focus:ring-blue-500"
                      />

                      {password && (
                        <div className="space-y-1 mt-2">
                          <div className="flex justify-between items-center">
                            <span className="text-xs text-gray-400">Password strength</span>
                            <span
                              className={`text-xs ${
                                passwordStrength < 40
                                  ? "text-red-400"
                                  : passwordStrength < 70
                                    ? "text-yellow-400"
                                    : "text-green-400"
                              }`}
                            >
                              {getPasswordStrengthLabel().label}
                            </span>
                          </div>
                          <Progress
                            value={passwordStrength}
                            className="h-1 bg-gray-700"
                            indicatorClassName={getPasswordStrengthLabel().color}
                          />
                        </div>
                      )}
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="confirmPassword" className="text-gray-300">
                        Confirm Password
                      </Label>
                      <Input
                        id="confirmPassword"
                        type="password"
                        placeholder="Confirm your password"
                        value={confirmPassword}
                        onChange={handleConfirmPasswordChange}
                        className="bg-gray-700/50 border-gray-600 text-white placeholder:text-gray-500 focus:border-blue-500 focus:ring-blue-500"
                      />
                      {password && confirmPassword && password !== confirmPassword && (
                        <p className="text-xs text-red-400 mt-1">Passwords do not match</p>
                      )}
                    </div>

                    <Button
                      onClick={createAccount}
                      className="w-full bg-blue-600 hover:bg-blue-500 text-white transition-colors"
                      disabled={isCreating || !password || password !== confirmPassword || passwordStrength < 40}
                    >
                      {isCreating ? (
                        <span className="flex items-center justify-center gap-2">
                          <svg
                            className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                          >
                            <circle
                              className="opacity-25"
                              cx="12"
                              cy="12"
                              r="10"
                              stroke="currentColor"
                              strokeWidth="4"
                            ></circle>
                            <path
                              className="opacity-75"
                              fill="currentColor"
                              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                            ></path>
                          </svg>
                          Creating Wallets...
                        </span>
                      ) : (
                        <span className="flex items-center justify-center gap-2">
                          <KeyIcon className="h-4 w-4" />
                          Create Wallets
                        </span>
                      )}
                    </Button>
                  </div>
                ) : (
                  <div className="space-y-6">
                    <Alert className="bg-yellow-500/20 border-yellow-600/50 text-yellow-200">
                      <AlertTriangleIcon className="h-4 w-4" />
                      <AlertTitle className="text-yellow-200 font-medium">Important Security Notice</AlertTitle>
                      <AlertDescription className="text-yellow-200/80">
                        Your recovery phrase is the only way to restore your wallet. Write it down and keep it in a
                        secure location.
                      </AlertDescription>
                    </Alert>

                    <div className="space-y-2">
                      <div className="flex items-center justify-between">
                        <Label className="text-gray-300 flex items-center gap-2">
                          <ShieldIcon className="h-4 w-4 text-blue-400" />
                          Recovery Phrase
                        </Label>
                        <div className="flex gap-2">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => setShowMnemonic(!showMnemonic)}
                            className="h-8 w-8 p-0 text-gray-400 hover:text-white hover:bg-gray-700/50"
                          >
                            {showMnemonic ? <EyeOffIcon className="h-4 w-4" /> : <EyeIcon className="h-4 w-4" />}
                            <span className="sr-only">{showMnemonic ? "Hide" : "Show"} recovery phrase</span>
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={copyToClipboard}
                            className="h-8 w-8 p-0 text-gray-400 hover:text-white hover:bg-gray-700/50"
                          >
                            {copied ? (
                              <CheckIcon className="h-4 w-4 text-green-400" />
                            ) : (
                              <CopyIcon className="h-4 w-4" />
                            )}
                            <span className="sr-only">Copy recovery phrase</span>
                          </Button>
                        </div>
                      </div>

                      <div className="relative">
                        <div
                          className={`p-4 bg-gray-700/50 rounded-md font-mono text-sm break-all ${showMnemonic ? "text-white" : "blur-sm select-none"}`}
                        >
                          {generatedMnemonic}
                        </div>
                        {!showMnemonic && (
                          <div className="absolute inset-0 flex items-center justify-center">
                            <Button
                              variant="ghost"
                              onClick={() => setShowMnemonic(true)}
                              className="text-blue-400 hover:text-blue-300 hover:bg-gray-700/50"
                            >
                              <EyeIcon className="h-4 w-4 mr-2" />
                              Click to reveal
                            </Button>
                          </div>
                        )}
                      </div>
                    </div>

                    <div className="space-y-4">
                      <Label className="text-gray-300">Storage Preference</Label>
                      <RadioGroup value={mnemonicOption} onValueChange={setMnemonicOption} className="space-y-3">
                        <div className="flex items-start space-x-3 space-y-0 rounded-md border border-gray-700 p-3 hover:bg-gray-700/30 transition-colors">
                          <RadioGroupItem value="self" id="self" className="mt-1 border-gray-600 text-blue-500" />
                          <div className="flex-1">
                            <Label
                              htmlFor="self"
                              className="text-sm font-medium text-white flex items-center cursor-pointer"
                            >
                              <ShieldIcon className="h-4 w-4 mr-2 text-blue-400" />
                              Self-custody
                            </Label>
                            <p className="text-xs text-gray-400 mt-1">
                              I'll manage my own recovery phrase. I understand I'm fully responsible for keeping it
                              safe.
                            </p>
                          </div>
                        </div>

                        <div className="flex items-start space-x-3 space-y-0 rounded-md border border-gray-700 p-3 hover:bg-gray-700/30 transition-colors">
                          <RadioGroupItem value="rely" id="rely" className="mt-1 border-gray-600 text-blue-500" />
                          <div className="flex-1">
                            <Label
                              htmlFor="rely"
                              className="text-sm font-medium text-white flex items-center cursor-pointer"
                            >
                              <ServerIcon className="h-4 w-4 mr-2 text-green-400" />
                              Service-managed
                            </Label>
                            <p className="text-xs text-gray-400 mt-1">
                              Store my recovery phrase with the service. The phrase will be encrypted and stored
                              securely.
                            </p>
                          </div>
                        </div>
                      </RadioGroup>
                    </div>

                    <Button
                      className="w-full bg-blue-600 hover:bg-blue-500 text-white transition-colors"
                      onClick={handleMnemonicOptionChange}
                      disabled={isStoring}
                    >
                      {isStoring ? (
                        <span className="flex items-center justify-center gap-2">
                          <svg
                            className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                          >
                            <circle
                              className="opacity-25"
                              cx="12"
                              cy="12"
                              r="10"
                              stroke="currentColor"
                              strokeWidth="4"
                            ></circle>
                            <path
                              className="opacity-75"
                              fill="currentColor"
                              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                            ></path>
                          </svg>
                          Processing...
                        </span>
                      ) : (
                        <span className="flex items-center justify-center gap-2">
                          Continue
                          <ArrowRightIcon className="h-4 w-4" />
                        </span>
                      )}
                    </Button>
                  </div>
                )}
                {walletDetails && (
                  <div className="space-y-4 mt-4 pt-4 border-t border-gray-700">
                    <h3 className="text-lg font-medium text-white">Your Wallet Details</h3>

                    <div className="space-y-3">
                      <div className="space-y-2">
                        <Label className="text-gray-300 flex items-center gap-2">
                          <svg
                            className="h-4 w-4 text-blue-400"
                            viewBox="0 0 24 24"
                            fill="none"
                            xmlns="http://www.w3.org/2000/svg"
                          >
                            <path
                              d="M11.9975 2L18.7086 5.85V13.55L11.9975 17.4L5.29102 13.55V5.85L11.9975 2Z"
                              stroke="currentColor"
                              strokeWidth="2"
                              strokeLinecap="round"
                              strokeLinejoin="round"
                            />
                            <path
                              d="M12 22V17.4"
                              stroke="currentColor"
                              strokeWidth="2"
                              strokeLinecap="round"
                              strokeLinejoin="round"
                            />
                            <path
                              d="M18.7086 13.55L21.9975 15.4"
                              stroke="currentColor"
                              strokeWidth="2"
                              strokeLinecap="round"
                              strokeLinejoin="round"
                            />
                            <path
                              d="M5.29102 13.55L2.00195 15.4"
                              stroke="currentColor"
                              strokeWidth="2"
                              strokeLinecap="round"
                              strokeLinejoin="round"
                            />
                          </svg>
                          Ethereum Address
                        </Label>
                        <div className="flex items-center">
                          <div className="p-3 bg-gray-700/50 rounded-md font-mono text-sm break-all text-green-300 flex-1">
                            {walletDetails.ethereum}
                          </div>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => {
                              navigator.clipboard.writeText(walletDetails.ethereum)
                              setCopied(true)
                              setTimeout(() => setCopied(false), 2000)
                            }}
                            className="ml-2 h-8 w-8 p-0 text-gray-400 hover:text-white hover:bg-gray-700/50"
                          >
                            {copied ? (
                              <CheckIcon className="h-4 w-4 text-green-400" />
                            ) : (
                              <CopyIcon className="h-4 w-4" />
                            )}
                            <span className="sr-only">Copy Ethereum address</span>
                          </Button>
                        </div>
                      </div>

                      <div className="space-y-2">
                        <Label className="text-gray-300 flex items-center gap-2">
                          <svg
                            className="h-4 w-4 text-purple-400"
                            viewBox="0 0 24 24"
                            fill="none"
                            xmlns="http://www.w3.org/2000/svg"
                          >
                            <path
                              d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z"
                              stroke="currentColor"
                              strokeWidth="2"
                              strokeLinecap="round"
                              strokeLinejoin="round"
                            />
                            <path
                              d="M7.5 12.5L10.5 15.5L16.5 9.5"
                              stroke="currentColor"
                              strokeWidth="2"
                              strokeLinecap="round"
                              strokeLinejoin="round"
                            />
                          </svg>
                          Solana Address
                        </Label>
                        <div className="flex items-center">
                          <div className="p-3 bg-gray-700/50 rounded-md font-mono text-sm break-all text-purple-300 flex-1">
                            {walletDetails.solana}
                          </div>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => {
                              navigator.clipboard.writeText(walletDetails.solana)
                              setCopied(true)
                              setTimeout(() => setCopied(false), 2000)
                            }}
                            className="ml-2 h-8 w-8 p-0 text-gray-400 hover:text-white hover:bg-gray-700/50"
                          >
                            {copied ? (
                              <CheckIcon className="h-4 w-4 text-green-400" />
                            ) : (
                              <CopyIcon className="h-4 w-4" />
                            )}
                            <span className="sr-only">Copy Solana address</span>
                          </Button>
                        </div>
                      </div>
                    </div>
                  </div>
                )}
              </TabsContent>

              <TabsContent value="import" className="space-y-6 mt-4">
                <div className="space-y-2">
                  <Label htmlFor="mnemonic" className="text-gray-300 flex items-center gap-2">
                    <ShieldIcon className="h-4 w-4 text-blue-400" />
                    Recovery Phrase
                  </Label>
                  <Input
                    id="mnemonic"
                    value={mnemonic}
                    onChange={handleMnemonicChange}
                    placeholder="Enter your 12 or 24 word recovery phrase"
                    className="bg-gray-700/50 border-gray-600 text-white placeholder:text-gray-500 focus:border-blue-500 focus:ring-blue-500"
                  />
                  <p className="text-xs text-gray-400">
                    Enter your recovery phrase (12 or 24 words separated by spaces)
                  </p>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="importPassword" className="text-gray-300 flex items-center gap-2">
                    <LockIcon className="h-4 w-4 text-blue-400" />
                    Wallet Password
                  </Label>
                  <Input
                    id="importPassword"
                    type="password"
                    placeholder="Create a password for your wallet"
                    value={password}
                    onChange={handlePasswordChange}
                    className="bg-gray-700/50 border-gray-600 text-white placeholder:text-gray-500 focus:border-blue-500 focus:ring-blue-500"
                  />
                </div>

                <Button
                  onClick={generatePrivateKey}
                  className="w-full bg-blue-600 hover:bg-blue-500 text-white transition-colors"
                  disabled={!mnemonic.trim() || !password.trim()}
                >
                  <KeyIcon className="h-4 w-4 mr-2" />
                  Import Wallet
                </Button>

                {privateKey && privateKey !== "Invalid mnemonic" && (
                  <div className="space-y-4 pt-4 border-t border-gray-700">
                    <h3 className="text-lg font-medium text-white">Your Wallet Details</h3>

                    <div className="space-y-3">
                      <div className="space-y-2">
                        <Label className="text-gray-300 flex items-center gap-2">
                          <svg
                            className="h-4 w-4 text-blue-400"
                            viewBox="0 0 24 24"
                            fill="none"
                            xmlns="http://www.w3.org/2000/svg"
                          >
                            <path
                              d="M11.9975 2L18.7086 5.85V13.55L11.9975 17.4L5.29102 13.55V5.85L11.9975 2Z"
                              stroke="currentColor"
                              strokeWidth="2"
                              strokeLinecap="round"
                              strokeLinejoin="round"
                            />
                            <path
                              d="M12 22V17.4"
                              stroke="currentColor"
                              strokeWidth="2"
                              strokeLinecap="round"
                              strokeLinejoin="round"
                            />
                            <path
                              d="M18.7086 13.55L21.9975 15.4"
                              stroke="currentColor"
                              strokeWidth="2"
                              strokeLinecap="round"
                              strokeLinejoin="round"
                            />
                            <path
                              d="M5.29102 13.55L2.00195 15.4"
                              stroke="currentColor"
                              strokeWidth="2"
                              strokeLinecap="round"
                              strokeLinejoin="round"
                            />
                          </svg>
                          Ethereum Address
                        </Label>
                        <div className="p-3 bg-gray-700/50 rounded-md font-mono text-sm break-all text-green-300">
                          {new ethers.Wallet(privateKey).address}
                        </div>
                      </div>

                      {walletDetails && (
                        <div className="space-y-2">
                          <Label className="text-gray-300 flex items-center gap-2">
                            <svg
                              className="h-4 w-4 text-purple-400"
                              viewBox="0 0 24 24"
                              fill="none"
                              xmlns="http://www.w3.org/2000/svg"
                            >
                              <path
                                d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z"
                                stroke="currentColor"
                                strokeWidth="2"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                              />
                              <path
                                d="M7.5 12.5L10.5 15.5L16.5 9.5"
                                stroke="currentColor"
                                strokeWidth="2"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                              />
                            </svg>
                            Solana Address
                          </Label>
                          <div className="p-3 bg-gray-700/50 rounded-md font-mono text-sm break-all text-purple-300">
                            {walletDetails.solana}
                          </div>
                        </div>
                      )}
                    </div>

                    <div className="flex justify-end mt-4">
                      <Button
                        onClick={() => router.push("/register")}
                        className="bg-blue-600 hover:bg-blue-500 text-white transition-colors"
                      >
                        <span className="flex items-center justify-center gap-2">
                          Continue
                          <ArrowRightIcon className="h-4 w-4" />
                        </span>
                      </Button>
                    </div>
                  </div>
                )}

                {privateKey === "Invalid mnemonic" && (
                  <Alert className="bg-red-500/20 border-red-600/50 text-red-200">
                    <AlertCircle className="h-4 w-4" />
                    <AlertDescription>Invalid recovery phrase. Please check and try again.</AlertDescription>
                  </Alert>
                )}
              </TabsContent>
            </Tabs>
          </CardHeader>

          <CardContent className="pt-0">{/* Card content is now handled by the TabsContent components */}</CardContent>
        </Card>

        <div className="mt-6 flex items-center justify-center">
          <span className="inline-flex items-center text-xs text-gray-500">
            <ShieldIcon className="h-4 w-4 mr-1" />
            All wallet operations are performed locally in your browser
          </span>
        </div>
      </div>
    </div>
  )
}
