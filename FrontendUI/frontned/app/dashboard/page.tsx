"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Connection, PublicKey, LAMPORTS_PER_SOL, clusterApiUrl } from "@solana/web3.js"
import {
  Search,
  Bell,
  User,
  Settings,
  LogOut,
  ArrowUpRight,
  ArrowDownLeft,
  Copy,
  Check,
  RefreshCw,
  ChevronDown,
  Wallet,
  BarChart3,
  History,
  CreditCard,
  ExternalLink,
  Shield,
  Sparkles,
} from "lucide-react"

import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Separator } from "@/components/ui/separator"
import { Badge } from "@/components/ui/badge"

export default function Dashboard() {
  const router = useRouter()
  const [balance, setBalance] = useState<number>(0)
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [walletAddress, setWalletAddress] = useState<string>("")
  const [copied, setCopied] = useState<boolean>(false)
  const [recipientAddress, setRecipientAddress] = useState<string>("")
  const [amount, setAmount] = useState<string>("")
  const [showBanner, setShowBanner] = useState<boolean>(true)
  const [transactions, setTransactions] = useState<any[]>([
    {
      id: "tx1",
      type: "receive",
      amount: 2.5,
      from: "8Kv...j3M",
      to: "Your wallet",
      date: "Today, 10:23 AM",
      status: "completed",
    },
    {
      id: "tx2",
      type: "send",
      amount: 0.75,
      from: "Your wallet",
      to: "3xR...p7B",
      date: "Yesterday, 3:45 PM",
      status: "completed",
    },
    {
      id: "tx3",
      type: "receive",
      amount: 1.2,
      from: "5Gs...k9P",
      to: "Your wallet",
      date: "Apr 5, 2023",
      status: "completed",
    },
  ])

  useEffect(() => {
    // Get wallet address from localStorage or session
    // const storedWalletData = localStorage.getItem("solanaWallet")
    // if (storedWalletData) {
    //   try {
    //     const walletData = JSON.parse(storedWalletData)
    //     setWalletAddress(walletData.publicKey)
    //     // fetchBalance(walletData.publicKey)
    //   } catch (error) {
    //     console.error("Error parsing wallet data:", error)
    //   }
    // } else {
      const demoAddress = "7dCZzk2XdsGSsXtnwtRjZXXXzR4B4DQ6cWzZeML9SSto"
      setWalletAddress(demoAddress)
      fetchBalance(demoAddress)
    // }
  }, [])

  const fetchBalance = async (address: string) => {
    setIsLoading(true)
    try {
      const connection = new Connection(clusterApiUrl("devnet"), "confirmed")
      const publicKey = new PublicKey(address)
      const balance = await connection.getBalance(publicKey)
      setBalance(balance / LAMPORTS_PER_SOL)
    } catch (error) {
      console.error("Error fetching balance:", error)
      // For demo purposes, set a placeholder balance
      setBalance(4.75)
    } finally {
      setIsLoading(false)
    }
  }

  const refreshBalance = () => {
    if (walletAddress) {
      fetchBalance(walletAddress)
    }
  }

  const copyToClipboard = () => {
    navigator.clipboard.writeText(walletAddress)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const formatAddress = (address: string) => {
    if (!address) return ""
    return `${address.slice(0, 6)}...${address.slice(-4)}`
  }

  const handleSend = () => {
    // In a real app, this would send a transaction
    alert(`Transaction of ${amount} SOL to ${recipientAddress} would be sent here`)

    // Add to transactions for demo
    const newTx = {
      id: `tx${Date.now()}`,
      type: "send",
      amount: Number.parseFloat(amount),
      from: "Your wallet",
      to: formatAddress(recipientAddress),
      date: "Just now",
      status: "pending",
    }

    setTransactions([newTx, ...transactions])
    setRecipientAddress("")
    setAmount("")
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-purple-50 text-gray-800">
      <div className="flex flex-col h-screen">
        {/* Banner */}
        {showBanner && (
          <div className="bg-gradient-to-r from-blue-600 to-purple-600 text-white py-3 px-4 relative">
            <div className="container mx-auto flex items-center justify-center">
              <Shield className="h-5 w-5 mr-2 text-white" />
              <p className="text-sm md:text-base font-medium">
                We rely on two blockchains: Swan Chain (our own chain which use the ethereum Technology to record user action on chain ) and Solana (to make payments)
              </p>
              <Button
                variant="ghost"
                size="sm"
                className="ml-2 text-white hover:bg-white/20 p-1 h-auto"
                onClick={() => setShowBanner(false)}
              >
                <span className="sr-only">Close</span>
                <svg width="15" height="15" viewBox="0 0 15 15" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path
                    d="M11.7816 4.03157C12.0062 3.80702 12.0062 3.44295 11.7816 3.2184C11.5571 2.99385 11.193 2.99385 10.9685 3.2184L7.50005 6.68682L4.03164 3.2184C3.80708 2.99385 3.44301 2.99385 3.21846 3.2184C2.99391 3.44295 2.99391 3.80702 3.21846 4.03157L6.68688 7.49999L3.21846 10.9684C2.99391 11.193 2.99391 11.557 3.21846 11.7816C3.44301 12.0061 3.80708 12.0061 4.03164 11.7816L7.50005 8.31316L10.9685 11.7816C11.193 12.0061 11.5571 12.0061 11.7816 11.7816C12.0062 11.557 12.0062 11.193 11.7816 10.9684L8.31322 7.49999L11.7816 4.03157Z"
                    fill="currentColor"
                    fillRule="evenodd"
                    clipRule="evenodd"
                  ></path>
                </svg>
              </Button>
            </div>
            <div className="absolute inset-0 bg-white/10 rounded-lg"></div>
          </div>
        )}

        {/* Header */}
        <header className="bg-white shadow-sm border-b border-gray-100 py-3">
          <div className="container mx-auto px-4 flex justify-between items-center">
            <div className="flex items-center gap-2">
              <div className="bg-gradient-to-r from-blue-500 to-purple-500 p-2 rounded-lg">
                <Wallet className="h-5 w-5 text-white" />
              </div>
              <span className="font-bold text-xl bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
              BlockChain Service Offering Platform
              </span>
            </div>

            <div className="flex-1 max-w-md mx-4">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="Search transactions..."
                  className="pl-10 bg-white border-gray-200 text-gray-800 placeholder:text-gray-400 focus:border-blue-500 focus:ring-blue-500 rounded-full"
                />
              </div>
            </div>

            <div className="flex items-center gap-4">
              <Button
                variant="ghost"
                size="icon"
                className="text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-full"
              >
                <Bell className="h-5 w-5" />
              </Button>

              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    className="flex items-center gap-2 text-gray-700 hover:text-blue-600 hover:bg-blue-50 rounded-full"
                  >
                    <Avatar className="h-8 w-8 border-2 border-blue-100">
                      <AvatarImage src="/placeholder-user.jpg" />
                      <AvatarFallback className="bg-gradient-to-r from-blue-500 to-purple-500 text-white">
                        JD
                      </AvatarFallback>
                    </Avatar>
                    <span className="hidden md:inline-block font-medium">John Doe</span>
                    <ChevronDown className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56 bg-white border-gray-200 text-gray-800">
                  <DropdownMenuLabel>My Account</DropdownMenuLabel>
                  <DropdownMenuSeparator className="bg-gray-200" />
                  <DropdownMenuItem className="hover:bg-blue-50 cursor-pointer">
                    <User className="mr-2 h-4 w-4 text-blue-500" />
                    <span>Profile</span>
                  </DropdownMenuItem>
                  <DropdownMenuItem className="hover:bg-blue-50 cursor-pointer">
                    <Settings className="mr-2 h-4 w-4 text-blue-500" />
                    <span>Settings</span>
                  </DropdownMenuItem>
                  <DropdownMenuSeparator className="bg-gray-200" />
                  <DropdownMenuItem className="hover:bg-blue-50 cursor-pointer" onClick={() => router.push("/")}>
                    <LogOut className="mr-2 h-4 w-4 text-blue-500" />
                    <span>Log out</span>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </header>

        {/* Main content */}
        <main className="flex-1 overflow-auto">
          <div className="container mx-auto p-4 md:p-6 grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Left column - Wallet and actions */}
            <div className="lg:col-span-2 space-y-6">
              {/* Wallet card */}
              <Card className="border border-gray-100 bg-white shadow-lg rounded-xl overflow-hidden">
                <div className="bg-gradient-to-r from-blue-500 to-purple-600 h-3"></div>
                <CardHeader className="pb-2">
                  <CardTitle className="text-xl flex justify-between items-center">
                    <div className="flex items-center gap-2">
                      <span>Solana Wallet</span>
                      <Badge variant="outline" className="bg-blue-50 text-blue-600 border-blue-200 font-normal">
                        Devnet
                      </Badge>
                    </div>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={refreshBalance}
                      className="text-gray-500 hover:text-blue-600 hover:bg-blue-50 rounded-full"
                      disabled={isLoading}
                    >
                      <RefreshCw className={`h-4 w-4 ${isLoading ? "animate-spin" : ""}`} />
                    </Button>
                  </CardTitle>
                  <CardDescription className="flex items-center gap-2 text-gray-500">
                    <span>{formatAddress(walletAddress)}</span>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={copyToClipboard}
                      className="h-6 w-6 p-0 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-full"
                    >
                      {copied ? <Check className="h-3 w-3 text-green-500" /> : <Copy className="h-3 w-3" />}
                    </Button>
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex flex-col items-center justify-center py-6">
                    <div className="text-4xl font-bold mb-2 bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                      {isLoading ? (
                        <div className="h-10 w-24 bg-gray-100 animate-pulse rounded"></div>
                      ) : (
                        `${balance.toFixed(4)} SOL`
                      )}
                    </div>
                    <div className="text-gray-500 text-sm">
                      {isLoading ? (
                        <div className="h-4 w-16 bg-gray-100 animate-pulse rounded"></div>
                      ) : (
                        `≈ ${(balance * 150).toFixed(2)} USD`
                      )}
                    </div>
                  </div>

                  <div className="flex gap-4 mt-4">
                    <Dialog>
                      <DialogTrigger asChild>
                        <Button className="flex-1 bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white shadow-md hover:shadow-lg transition-all duration-200 rounded-full">
                          <ArrowUpRight className="mr-2 h-4 w-4" />
                          Send
                        </Button>
                      </DialogTrigger>
                      <DialogContent className="bg-white text-gray-800 border-gray-200 rounded-xl">
                        <DialogHeader>
                          <DialogTitle>Send SOL</DialogTitle>
                          <DialogDescription className="text-gray-500">
                            Send SOL to another wallet address.
                          </DialogDescription>
                        </DialogHeader>
                        <div className="space-y-4 py-4">
                          <div className="space-y-2">
                            <Label htmlFor="recipient">Recipient Address</Label>
                            <Input
                              id="recipient"
                              value={recipientAddress}
                              onChange={(e) => setRecipientAddress(e.target.value)}
                              placeholder="Enter Solana address"
                              className="bg-white border-gray-200 text-gray-800 focus:border-blue-500 focus:ring-blue-500"
                            />
                          </div>
                          <div className="space-y-2">
                            <Label htmlFor="amount">Amount (SOL)</Label>
                            <Input
                              id="amount"
                              type="number"
                              value={amount}
                              onChange={(e) => setAmount(e.target.value)}
                              placeholder="0.00"
                              className="bg-white border-gray-200 text-gray-800 focus:border-blue-500 focus:ring-blue-500"
                            />
                            <div className="text-xs text-gray-500 flex justify-between">
                              <span>Available: {balance.toFixed(4)} SOL</span>
                              <Button
                                variant="link"
                                className="h-auto p-0 text-blue-600"
                                onClick={() => setAmount(balance.toString())}
                              >
                                Max
                              </Button>
                            </div>
                          </div>
                        </div>
                        <DialogFooter>
                          <Button
                            onClick={handleSend}
                            disabled={
                              !recipientAddress ||
                              !amount ||
                              Number.parseFloat(amount) <= 0 ||
                              Number.parseFloat(amount) > balance
                            }
                            className="bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white rounded-full"
                          >
                            Send SOL
                          </Button>
                        </DialogFooter>
                      </DialogContent>
                    </Dialog>

                    <Dialog>
                      <DialogTrigger asChild>
                        <Button
                          variant="outline"
                          className="flex-1 border-blue-200 text-blue-600 hover:bg-blue-50 hover:border-blue-300 rounded-full"
                        >
                          <ArrowDownLeft className="mr-2 h-4 w-4" />
                          Receive
                        </Button>
                      </DialogTrigger>
                      <DialogContent className="bg-white text-gray-800 border-gray-200 rounded-xl">
                        <DialogHeader>
                          <DialogTitle>Receive SOL</DialogTitle>
                          <DialogDescription className="text-gray-500">
                            Share your address to receive SOL.
                          </DialogDescription>
                        </DialogHeader>
                        <div className="flex flex-col items-center justify-center py-4">
                          <div className="bg-white p-4 rounded-lg mb-4 border border-gray-200 shadow-md">
                            <div className="h-48 w-48 bg-gradient-to-r from-blue-50 to-purple-50 flex items-center justify-center rounded-lg">
                              <span className="text-gray-400">QR Code</span>
                            </div>
                          </div>
                          <div className="w-full p-3 bg-gray-50 rounded-md font-mono text-sm break-all text-gray-800 mb-2 border border-gray-200">
                            {walletAddress}
                          </div>
                          <Button
                            variant="outline"
                            onClick={copyToClipboard}
                            className="border-blue-200 text-blue-600 hover:bg-blue-50 hover:border-blue-300 rounded-full"
                          >
                            {copied ? <Check className="mr-2 h-4 w-4" /> : <Copy className="mr-2 h-4 w-4" />}
                            Copy Address
                          </Button>
                        </div>
                      </DialogContent>
                    </Dialog>
                  </div>
                </CardContent>
              </Card>

              {/* Blockchain Info Card */}
              <Card className="border border-gray-100 bg-white shadow-md rounded-xl overflow-hidden">
                <div className="bg-gradient-to-r from-purple-500 to-pink-500 h-1"></div>
                <CardHeader className="pb-2">
                  <CardTitle className="text-lg flex items-center gap-2">
                    <Sparkles className="h-5 w-5 text-purple-500" />
                    Blockchain Technology
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="bg-purple-50 rounded-xl p-4 border border-purple-100 relative overflow-hidden group hover:shadow-md transition-all duration-200">
                      <div className="absolute top-0 right-0 w-16 h-16 bg-purple-200 rounded-bl-full opacity-30 group-hover:opacity-50 transition-opacity"></div>
                      <h3 className="font-medium text-purple-700 mb-1 flex items-center gap-2">
                        <Shield className="h-4 w-4" />
                        Swan Chain
                      </h3>
                      <p className="text-sm text-purple-600">
                        Our proprietary blockchain used for secure record-keeping and data integrity.
                      </p>
                      <Button
                        variant="link"
                        size="sm"
                        className="text-purple-700 p-0 h-auto mt-2 text-xs flex items-center"
                      >
                        Learn more <ExternalLink className="h-3 w-3 ml-1" />
                      </Button>
                    </div>

                    <div className="bg-blue-50 rounded-xl p-4 border border-blue-100 relative overflow-hidden group hover:shadow-md transition-all duration-200">
                      <div className="absolute top-0 right-0 w-16 h-16 bg-blue-200 rounded-bl-full opacity-30 group-hover:opacity-50 transition-opacity"></div>
                      <h3 className="font-medium text-blue-700 mb-1 flex items-center gap-2">
                        <svg className="h-4 w-4" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
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
                        Solana
                      </h3>
                      <p className="text-sm text-blue-600">
                        Fast, secure, and low-cost blockchain used for all payment transactions.
                      </p>
                      <Button
                        variant="link"
                        size="sm"
                        className="text-blue-700 p-0 h-auto mt-2 text-xs flex items-center"
                      >
                        Learn more <ExternalLink className="h-3 w-3 ml-1" />
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Transactions */}
              <Card className="border border-gray-100 bg-white shadow-md rounded-xl overflow-hidden">
                <div className="bg-gradient-to-r from-blue-500 to-purple-600 h-1"></div>
                <CardHeader>
                  <CardTitle className="text-lg flex items-center gap-2">
                    <History className="h-5 w-5 text-blue-500" />
                    Recent Transactions
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <Tabs defaultValue="all">
                    <TabsList className="bg-gray-100 mb-4 p-1 rounded-lg">
                      <TabsTrigger
                        value="all"
                        className="data-[state=active]:bg-white data-[state=active]:text-blue-600 data-[state=active]:shadow-sm rounded-md"
                      >
                        All
                      </TabsTrigger>
                      <TabsTrigger
                        value="sent"
                        className="data-[state=active]:bg-white data-[state=active]:text-blue-600 data-[state=active]:shadow-sm rounded-md"
                      >
                        Sent
                      </TabsTrigger>
                      <TabsTrigger
                        value="received"
                        className="data-[state=active]:bg-white data-[state=active]:text-blue-600 data-[state=active]:shadow-sm rounded-md"
                      >
                        Received
                      </TabsTrigger>
                    </TabsList>

                    <TabsContent value="all" className="space-y-4">
                      {transactions.map((tx) => (
                        <div
                          key={tx.id}
                          className="flex items-center justify-between p-3 rounded-xl bg-gray-50 hover:bg-gray-100 transition-colors border border-gray-100 group hover:shadow-sm"
                        >
                          <div className="flex items-center gap-3">
                            <div
                              className={`p-2 rounded-full ${tx.type === "receive" ? "bg-green-100 text-green-600" : "bg-blue-100 text-blue-600"} group-hover:scale-110 transition-transform`}
                            >
                              {tx.type === "receive" ? (
                                <ArrowDownLeft className="h-5 w-5" />
                              ) : (
                                <ArrowUpRight className="h-5 w-5" />
                              )}
                            </div>
                            <div>
                              <div className="font-medium">{tx.type === "receive" ? "Received SOL" : "Sent SOL"}</div>
                              <div className="text-sm text-gray-500">{tx.date}</div>
                            </div>
                          </div>
                          <div className="text-right">
                            <div
                              className={`font-medium ${tx.type === "receive" ? "text-green-600" : "text-blue-600"}`}
                            >
                              {tx.type === "receive" ? "+" : "-"}
                              {tx.amount} SOL
                            </div>
                            <div className="text-xs text-gray-500">
                              {tx.status === "completed" ? (
                                <span className="inline-flex items-center text-green-600">
                                  <Check className="h-3 w-3 mr-1" /> Completed
                                </span>
                              ) : (
                                <span className="inline-flex items-center text-orange-500">
                                  <svg
                                    className="animate-spin h-3 w-3 mr-1"
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
                                  Pending
                                </span>
                              )}
                            </div>
                          </div>
                        </div>
                      ))}
                    </TabsContent>

                    <TabsContent value="sent" className="space-y-4">
                      {transactions
                        .filter((tx) => tx.type === "send")
                        .map((tx) => (
                          <div
                            key={tx.id}
                            className="flex items-center justify-between p-3 rounded-xl bg-gray-50 hover:bg-gray-100 transition-colors border border-gray-100 group hover:shadow-sm"
                          >
                            <div className="flex items-center gap-3">
                              <div className="p-2 rounded-full bg-blue-100 text-blue-600 group-hover:scale-110 transition-transform">
                                <ArrowUpRight className="h-5 w-5" />
                              </div>
                              <div>
                                <div className="font-medium">Sent SOL</div>
                                <div className="text-sm text-gray-500">{tx.date}</div>
                              </div>
                            </div>
                            <div className="text-right">
                              <div className="font-medium text-blue-600">-{tx.amount} SOL</div>
                              <div className="text-xs text-gray-500">
                                {tx.status === "completed" ? (
                                  <span className="inline-flex items-center text-green-600">
                                    <Check className="h-3 w-3 mr-1" /> Completed
                                  </span>
                                ) : (
                                  <span className="inline-flex items-center text-orange-500">
                                    <svg
                                      className="animate-spin h-3 w-3 mr-1"
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
                                    Pending
                                  </span>
                                )}
                              </div>
                            </div>
                          </div>
                        ))}
                    </TabsContent>

                    <TabsContent value="received" className="space-y-4">
                      {transactions
                        .filter((tx) => tx.type === "receive")
                        .map((tx) => (
                          <div
                            key={tx.id}
                            className="flex items-center justify-between p-3 rounded-xl bg-gray-50 hover:bg-gray-100 transition-colors border border-gray-100 group hover:shadow-sm"
                          >
                            <div className="flex items-center gap-3">
                              <div className="p-2 rounded-full bg-green-100 text-green-600 group-hover:scale-110 transition-transform">
                                <ArrowDownLeft className="h-5 w-5" />
                              </div>
                              <div>
                                <div className="font-medium">Received SOL</div>
                                <div className="text-sm text-gray-500">{tx.date}</div>
                              </div>
                            </div>
                            <div className="text-right">
                              <div className="font-medium text-green-600">+{tx.amount} SOL</div>
                              <div className="text-xs text-gray-500">
                                {tx.status === "completed" ? (
                                  <span className="inline-flex items-center text-green-600">
                                    <Check className="h-3 w-3 mr-1" /> Completed
                                  </span>
                                ) : (
                                  <span className="inline-flex items-center text-orange-500">
                                    <svg
                                      className="animate-spin h-3 w-3 mr-1"
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
                                    Pending
                                  </span>
                                )}
                              </div>
                            </div>
                          </div>
                        ))}
                    </TabsContent>
                  </Tabs>
                </CardContent>
                <CardFooter className="border-t border-gray-100 pt-4 flex justify-center">
                  <Button variant="link" className="text-blue-600">
                    View All Transactions
                  </Button>
                </CardFooter>
              </Card>
            </div>

            {/* Right column - User details and stats */}
            <div className="space-y-6">
              {/* User profile card */}
              <Card className="border border-gray-100 bg-white shadow-md rounded-xl overflow-hidden">
                <div className="bg-gradient-to-r from-blue-500 to-purple-600 h-1"></div>
                <CardHeader className="pb-2">
                  <CardTitle className="text-lg flex items-center gap-2">
                    <User className="h-5 w-5 text-blue-500" />
                    User Profile
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex flex-col items-center py-4">
                    <Avatar className="h-20 w-20 mb-4 border-4 border-blue-100 shadow-md">
                      <AvatarImage src="/placeholder-user.jpg" />
                      <AvatarFallback className="bg-gradient-to-r from-blue-500 to-purple-500 text-xl text-white">
                        JD
                      </AvatarFallback>
                    </Avatar>
                    <h3 className="text-xl font-medium text-gray-800">John Doe</h3>
                    <p className="text-gray-500">john.doe@example.com</p>

                    <div className="w-full mt-6 space-y-4">
                      <div className="flex justify-between items-center">
                        <span className="text-gray-500">Account Type</span>
                        <Badge className="bg-blue-100 text-blue-600 hover:bg-blue-200 border-0">Premium</Badge>
                      </div>
                      <Separator className="bg-gray-100" />
                      <div className="flex justify-between items-center">
                        <span className="text-gray-500">Joined</span>
                        <span>April 2023</span>
                      </div>
                      <Separator className="bg-gray-100" />
                      <div className="flex justify-between items-center">
                        <span className="text-gray-500">Verification</span>
                        <Badge className="bg-green-100 text-green-600 hover:bg-green-200 border-0">Verified</Badge>
                      </div>
                    </div>
                  </div>
                </CardContent>
                <CardFooter className="border-t border-gray-100 pt-4">
                  <Button
                    variant="outline"
                    className="w-full border-blue-200 text-blue-600 hover:bg-blue-50 hover:border-blue-300 rounded-full"
                  >
                    <Settings className="mr-2 h-4 w-4" />
                    Edit Profile
                  </Button>
                </CardFooter>
              </Card>

              {/* Quick stats */}
              <Card className="border border-gray-100 bg-white shadow-md rounded-xl overflow-hidden">
                <div className="bg-gradient-to-r from-blue-500 to-purple-600 h-1"></div>
                <CardHeader className="pb-2">
                  <CardTitle className="text-lg flex items-center gap-2">
                    <BarChart3 className="h-5 w-5 text-blue-500" />
                    Quick Stats
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="bg-blue-50 p-4 rounded-xl border border-blue-100 hover:shadow-md transition-shadow group">
                      <div className="flex items-center gap-2 mb-2">
                        <BarChart3 className="h-4 w-4 text-blue-600 group-hover:scale-110 transition-transform" />
                        <span className="text-sm text-blue-600">Total Value</span>
                      </div>
                      <div className="text-xl font-medium text-blue-700">$712.50</div>
                    </div>
                    <div className="bg-purple-50 p-4 rounded-xl border border-purple-100 hover:shadow-md transition-shadow group">
                      <div className="flex items-center gap-2 mb-2">
                        <History className="h-4 w-4 text-purple-600 group-hover:scale-110 transition-transform" />
                        <span className="text-sm text-purple-600">Transactions</span>
                      </div>
                      <div className="text-xl font-medium text-purple-700">24</div>
                    </div>
                    <div className="bg-pink-50 p-4 rounded-xl border border-pink-100 hover:shadow-md transition-shadow group">
                      <div className="flex items-center gap-2 mb-2">
                        <CreditCard className="h-4 w-4 text-pink-600 group-hover:scale-110 transition-transform" />
                        <span className="text-sm text-pink-600">Wallets</span>
                      </div>
                      <div className="text-xl font-medium text-pink-700">2</div>
                    </div>
                    <div className="bg-indigo-50 p-4 rounded-xl border border-indigo-100 hover:shadow-md transition-shadow group">
                      <div className="flex items-center gap-2 mb-2">
                        <Wallet className="h-4 w-4 text-indigo-600 group-hover:scale-110 transition-transform" />
                        <span className="text-sm text-indigo-600">Networks</span>
                      </div>
                      <div className="text-xl font-medium text-indigo-700">2</div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </main>
      </div>
    </div>
  )
}
