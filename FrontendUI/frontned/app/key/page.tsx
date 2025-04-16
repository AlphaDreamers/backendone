"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Connection, PublicKey, LAMPORTS_PER_SOL, clusterApiUrl, Keypair, SystemProgram, Transaction } from "@solana/web3.js"
import { ethers } from "ethers"
import { QRCodeSVG } from "qrcode.react"
import {
  ArrowUpRight,
  ArrowDownLeft,
  Copy,
  Check,
  RefreshCw,
  Wallet,
  Shield,
  Sparkles,
  ExternalLink,
} from "lucide-react"

import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
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
import { getUserProfile, getToken, getUserEmailFromToken } from "@/app/lib/auth-utils"

export default function WalletPage() {
  const router = useRouter()
  const [solanaBalance, setSolanaBalance] = useState<number>(0)
  const [ethereumBalance, setEthereumBalance] = useState<string>("0")
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [solanaWalletAddress, setSolanaWalletAddress] = useState<string>("")
  const [ethereumWalletAddress, setEthereumWalletAddress] = useState<string>("")
  const [copied, setCopied] = useState<boolean>(false)
  const [recipientAddress, setRecipientAddress] = useState<string>("")
  const [amount, setAmount] = useState<string>("")
  const [activeWallet, setActiveWallet] = useState<"solana" | "ethereum">("solana")
  const [userProfile, setUserProfile] = useState<any>(null)
  const [profileError, setProfileError] = useState<string>("")
  const [walletError, setWalletError] = useState<string>("")

  useEffect(() => {
    // Fetch user profile data
    const fetchUserProfile = async () => {
      try {
        // Debug: Log all localStorage contents
        console.log("All localStorage contents:", Object.keys(localStorage).map(key => ({
          key,
          value: localStorage.getItem(key)
        })));
        
        // Always check localStorage first for wallet addresses
        const storedSolanaWallet = localStorage.getItem("solanaWallet");
        const storedEthereumWallet = localStorage.getItem("ethereumWallet");
        
        console.log("Raw localStorage values:", {
          solanaWallet: storedSolanaWallet,
          ethereumWallet: storedEthereumWallet
        });
        
        if (storedSolanaWallet) {
          console.log("Using Solana wallet from localStorage:", storedSolanaWallet);
          setSolanaWalletAddress(storedSolanaWallet);
          fetchSolanaBalance(storedSolanaWallet);
        } else {
          console.log("No Solana wallet found in localStorage");
        }
        
        if (storedEthereumWallet) {
          console.log("Using Ethereum wallet from localStorage:", storedEthereumWallet);
          setEthereumWalletAddress(storedEthereumWallet);
          fetchEthereumBalance(storedEthereumWallet);
        } else {
          console.log("No Ethereum wallet found in localStorage");
        }
        
        const token = getToken();
        if (!token) {
          router.push('/login');
          return;
        }

        const email = getUserEmailFromToken(token);
        if (!email) {
          setProfileError("Could not retrieve user email from token");
          return;
        }

        console.log("Fetching user profile for email:", email);
        const response = await getUserProfile(email, token);
        console.log("User profile response:", response);
        
        if (response.success && response.data) {
          console.log("Setting user profile data:", response.data);
          setUserProfile(response.data);
          
          // Only set wallet addresses from user profile if they're not already set from localStorage
          if (response.data.solana_address && !storedSolanaWallet) {
            setSolanaWalletAddress(response.data.solana_address);
            fetchSolanaBalance(response.data.solana_address);
          }
          
          if (response.data.ethereum_address && !storedEthereumWallet) {
            setEthereumWalletAddress(response.data.ethereum_address);
            fetchEthereumBalance(response.data.ethereum_address);
          }
        } else {
          setProfileError(response.message || "Failed to fetch user profile");
        }
      } catch (error) {
        setProfileError("An error occurred while fetching user profile");
        console.error(error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchUserProfile();
  }, [router]);

  const fetchSolanaBalance = async (address: string) => {
    try {
      const connection = new Connection(clusterApiUrl("devnet"), "confirmed");
      const publicKey = new PublicKey(address);
      const balance = await connection.getBalance(publicKey);
      setSolanaBalance(balance / LAMPORTS_PER_SOL);
    } catch (error) {
      console.error("Error fetching Solana balance:", error);
      setSolanaBalance(0);
    }
  };

  const fetchEthereumBalance = async (address: string) => {
    try {
      // Use a more reliable public Ethereum endpoint for Sepolia testnet
      const provider = new ethers.JsonRpcProvider("https://eth-sepolia.public.blastapi.io");
      
      // Add retry logic
      let retries = 3;
      let lastError;
      
      while (retries > 0) {
        try {
          const balance = await provider.getBalance(address);
          setEthereumBalance(ethers.formatEther(balance));
          return;
        } catch (error) {
          lastError = error;
          retries--;
          if (retries > 0) {
            // Wait for 1 second before retrying
            await new Promise(resolve => setTimeout(resolve, 1000));
          }
        }
      }
      
      throw lastError;
    } catch (error) {
      console.error("Error fetching Ethereum balance:", error);
      setEthereumBalance("0");
      setWalletError("Failed to fetch Ethereum balance. Please try again later.");
    }
  };

  const refreshBalance = () => {
    if (activeWallet === "solana" && solanaWalletAddress) {
      fetchSolanaBalance(solanaWalletAddress);
    } else if (activeWallet === "ethereum" && ethereumWalletAddress) {
      fetchEthereumBalance(ethereumWalletAddress);
    }
  };

  const copyToClipboard = () => {
    const address = activeWallet === "solana" ? solanaWalletAddress : ethereumWalletAddress;
    navigator.clipboard.writeText(address);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const formatAddress = (address: string) => {
    if (!address) return "";
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
  };

  const handleSend = async () => {
    if (!recipientAddress || !amount || Number.parseFloat(amount) <= 0) {
      alert("Please enter a valid recipient address and amount");
      return;
    }

    if (activeWallet === "solana") {
      try {
        // Get the private key from localStorage
        const storedPrivateKey = localStorage.getItem("solanaPrivateKey");
        if (!storedPrivateKey) {
          alert("Solana wallet private key not found. Please create a wallet in the key session.");
          return;
        }

        // Convert the stored private key string back to Uint8Array
        const privateKeyBytes = new Uint8Array(JSON.parse(storedPrivateKey));
        
        // Create a connection to the Solana devnet
        const connection = new Connection(clusterApiUrl("devnet"), "confirmed");
        
        // Create a keypair from the private key
        const fromKeypair = Keypair.fromSecretKey(privateKeyBytes);
        
        // Create a public key from the recipient address
        const toPublicKey = new PublicKey(recipientAddress);
        
        // Convert amount to lamports (1 SOL = 1,000,000,000 lamports)
        const lamports = Number.parseFloat(amount) * LAMPORTS_PER_SOL;
        
        // Create a transaction
        const transaction = new Transaction().add(
          SystemProgram.transfer({
            fromPubkey: fromKeypair.publicKey,
            toPubkey: toPublicKey,
            lamports: lamports,
          })
        );
        
        // Get the latest blockhash
        const { blockhash } = await connection.getLatestBlockhash();
        transaction.recentBlockhash = blockhash;
        transaction.feePayer = fromKeypair.publicKey;
        
        // Sign the transaction
        transaction.sign(fromKeypair);
        
        // Send the transaction
        const signature = await connection.sendRawTransaction(transaction.serialize());
        
        // Confirm the transaction
        await connection.confirmTransaction(signature);
        
        // Clear the form
        setRecipientAddress("");
        setAmount("");
        
        // Refresh the balance
        fetchSolanaBalance(solanaWalletAddress);
        
        alert(`Transaction sent successfully! Signature: ${signature}`);
      } catch (error: any) {
        console.error("Error sending Solana transaction:", error);
        alert(`Failed to send transaction: ${error.message || 'Unknown error'}`);
      }
    } else if (activeWallet === "ethereum") {
      try {
        // Get the private key from localStorage
        const storedPrivateKey = localStorage.getItem("ethereumPrivateKey");
        if (!storedPrivateKey) {
          alert("Ethereum wallet private key not found. Please create a wallet in the key session.");
          return;
        }
        
        // Create a provider for the Ethereum devnet (Sepolia testnet)
        const provider = new ethers.JsonRpcProvider("https://eth-sepolia.g.alchemy.com/v2/demo");
        
        // Create a wallet from the private key
        const wallet = new ethers.Wallet(storedPrivateKey, provider);
        
        // Convert amount to wei (1 ETH = 1,000,000,000,000,000,000 wei)
        const weiAmount = ethers.parseEther(amount);
        
        // Create a transaction
        const tx = {
          to: recipientAddress,
          value: weiAmount,
        };
        
        // Send the transaction
        const transaction = await wallet.sendTransaction(tx);
        
        // Wait for the transaction to be mined
        const receipt = await transaction.wait();
        
        // Clear the form
        setRecipientAddress("");
        setAmount("");
        
        // Refresh the balance
        fetchEthereumBalance(ethereumWalletAddress);
        
        alert(`Transaction sent successfully! Hash: ${receipt?.hash || 'Unknown hash'}`);
      } catch (error: any) {
        console.error("Error sending Ethereum transaction:", error);
        alert(`Failed to send transaction: ${error.message || 'Unknown error'}`);
      }
    }
  };

  const handleCreateWallet = async () => {
    // Redirect to the key page to create a wallet
    router.push('/key');
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 to-gray-800 text-gray-100">
      <div className="container mx-auto p-4 md:p-6">
        <h1 className="text-3xl font-bold mb-6">My Wallets</h1>
        
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Wallet card */}
          <Card className="border border-gray-700 bg-gray-800 shadow-lg rounded-xl overflow-hidden">
            <div className="bg-gradient-to-r from-blue-600 to-purple-600 h-3"></div>
            <CardHeader className="pb-2">
              <CardTitle className="text-xl flex justify-between items-center text-gray-100">
                <div className="flex items-center gap-2">
                  <span>{activeWallet === "solana" ? "Solana Wallet" : "Ethereum Wallet"}</span>
                  <Badge variant="outline" className="bg-blue-900/30 text-blue-300 border-blue-700 font-normal">
                    Devnet
                  </Badge>
                </div>
                <div className="flex items-center gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setActiveWallet("solana")}
                    className={`text-gray-400 hover:text-blue-400 hover:bg-gray-700 rounded-full ${activeWallet === "solana" ? "bg-gray-700 text-blue-400" : ""}`}
                    disabled={!solanaWalletAddress}
                  >
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
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setActiveWallet("ethereum")}
                    className={`text-gray-400 hover:text-blue-400 hover:bg-gray-700 rounded-full ${activeWallet === "ethereum" ? "bg-gray-700 text-blue-400" : ""}`}
                    disabled={!ethereumWalletAddress}
                  >
                    <svg className="h-4 w-4" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                      <path d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                      <path d="M12 6V18" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                      <path d="M8 10L12 6L16 10" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                      <path d="M8 14L12 18L16 14" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                    </svg>
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={refreshBalance}
                    className="text-gray-400 hover:text-blue-400 hover:bg-gray-700 rounded-full"
                    disabled={isLoading}
                  >
                    <RefreshCw className={`h-4 w-4 ${isLoading ? "animate-spin" : ""}`} />
                  </Button>
                </div>
              </CardTitle>
              {solanaWalletAddress || ethereumWalletAddress ? (
                <CardDescription className="flex items-center gap-2 text-gray-400">
                  <span>{formatAddress(activeWallet === "solana" ? solanaWalletAddress : ethereumWalletAddress)}</span>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={copyToClipboard}
                    className="h-6 w-6 p-0 text-gray-500 hover:text-blue-400 hover:bg-gray-700 rounded-full"
                  >
                    {copied ? <Check className="h-3 w-3 text-green-400" /> : <Copy className="h-3 w-3" />}
                  </Button>
                </CardDescription>
              ) : (
                <CardDescription className="text-yellow-400">
                  <div className="flex items-center gap-2">
                    <svg className="h-4 w-4" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                      <path d="M12 9V11M12 15H12.01M5.07183 19H18.9282C20.4678 19 21.4301 17.3333 20.6603 16L13.7321 4C12.9623 2.66667 11.0378 2.66667 10.268 4L3.33978 16C2.56998 17.3333 3.53223 19 5.07183 19Z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                    </svg>
                    <span>No wallet found. Please create a wallet in the key session.</span>
                  </div>
                </CardDescription>
              )}
            </CardHeader>
            <CardContent>
              {solanaWalletAddress || ethereumWalletAddress ? (
                <>
                  <div className="flex flex-col items-center justify-center py-6">
                    <div className="text-4xl font-bold mb-2 bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
                      {isLoading ? (
                        <div className="h-10 w-24 bg-gray-700 animate-pulse rounded"></div>
                      ) : (
                        `${activeWallet === "solana" ? solanaBalance.toFixed(4) : ethereumBalance} ${activeWallet === "solana" ? "SOL" : "ETH"}`
                      )}
                    </div>
                    <div className="text-gray-400 text-sm">
                      {isLoading ? (
                        <div className="h-4 w-16 bg-gray-700 animate-pulse rounded"></div>
                      ) : (
                        `≈ ${(activeWallet === "solana" ? solanaBalance * 150 : parseFloat(ethereumBalance) * 3000).toFixed(2)} USD`
                      )}
                    </div>
                  </div>

                  <div className="flex gap-4 mt-4">
                    <Dialog>
                      <DialogTrigger asChild>
                        <Button className="flex-1 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white shadow-md hover:shadow-lg transition-all duration-200 rounded-full">
                          <ArrowUpRight className="mr-2 h-4 w-4" />
                          Send
                        </Button>
                      </DialogTrigger>
                      <DialogContent className="bg-gray-800 text-gray-100 border-gray-700 rounded-xl">
                        <DialogHeader>
                          <DialogTitle>Send {activeWallet === "solana" ? "SOL" : "ETH"}</DialogTitle>
                          <DialogDescription className="text-gray-400">
                            Send {activeWallet === "solana" ? "SOL" : "ETH"} to another wallet address.
                          </DialogDescription>
                        </DialogHeader>
                        <div className="space-y-4 py-4">
                          <div className="space-y-2">
                            <Label htmlFor="recipient" className="text-gray-300">Recipient Address</Label>
                            <Input
                              id="recipient"
                              value={recipientAddress}
                              onChange={(e) => setRecipientAddress(e.target.value)}
                              placeholder={`Enter ${activeWallet === "solana" ? "Solana" : "Ethereum"} address`}
                              className="bg-gray-700 border-gray-600 text-gray-100 focus:border-blue-500 focus:ring-blue-500"
                            />
                          </div>
                          <div className="space-y-2">
                            <Label htmlFor="amount" className="text-gray-300">Amount ({activeWallet === "solana" ? "SOL" : "ETH"})</Label>
                            <Input
                              id="amount"
                              type="number"
                              value={amount}
                              onChange={(e) => setAmount(e.target.value)}
                              placeholder="0.00"
                              className="bg-gray-700 border-gray-600 text-gray-100 focus:border-blue-500 focus:ring-blue-500"
                            />
                            <div className="text-xs text-gray-400 flex justify-between">
                              <span>Available: {activeWallet === "solana" ? solanaBalance.toFixed(4) : ethereumBalance} {activeWallet === "solana" ? "SOL" : "ETH"}</span>
                              <Button
                                variant="link"
                                className="h-auto p-0 text-blue-400"
                                onClick={() => setAmount(activeWallet === "solana" ? solanaBalance.toString() : ethereumBalance)}
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
                              Number.parseFloat(amount) > (activeWallet === "solana" ? solanaBalance : parseFloat(ethereumBalance))
                            }
                            className="bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white rounded-full"
                          >
                            Send {activeWallet === "solana" ? "SOL" : "ETH"}
                          </Button>
                        </DialogFooter>
                      </DialogContent>
                    </Dialog>

                    <Dialog>
                      <DialogTrigger asChild>
                        <Button
                          variant="outline"
                          className="flex-1 border-blue-700 text-blue-400 hover:bg-gray-700 hover:border-blue-600 rounded-full"
                        >
                          <ArrowDownLeft className="mr-2 h-4 w-4" />
                          Receive
                        </Button>
                      </DialogTrigger>
                      <DialogContent className="bg-gray-800 text-gray-100 border-gray-700 rounded-xl">
                        <DialogHeader>
                          <DialogTitle>Receive {activeWallet === "solana" ? "SOL" : "ETH"}</DialogTitle>
                          <DialogDescription className="text-gray-400">
                            Share your address to receive {activeWallet === "solana" ? "SOL" : "ETH"}.
                          </DialogDescription>
                        </DialogHeader>
                        <div className="flex flex-col items-center justify-center py-4">
                          <div className="bg-gray-700 p-4 rounded-lg mb-4 border border-gray-600 shadow-md">
                            <QRCodeSVG 
                              value={activeWallet === "solana" ? solanaWalletAddress : ethereumWalletAddress} 
                              size={200}
                              level="H"
                              includeMargin={true}
                              bgColor="#1f2937"
                              fgColor="#ffffff"
                            />
                          </div>
                          <div className="w-full p-3 bg-gray-700 rounded-md font-mono text-sm break-all text-gray-300 mb-2 border border-gray-600">
                            {activeWallet === "solana" ? solanaWalletAddress : ethereumWalletAddress}
                          </div>
                          <Button
                            variant="outline"
                            onClick={copyToClipboard}
                            className="border-blue-700 text-blue-400 hover:bg-gray-700 hover:border-blue-600 rounded-full"
                          >
                            {copied ? <Check className="mr-2 h-4 w-4" /> : <Copy className="mr-2 h-4 w-4" />}
                            Copy Address
                          </Button>
                        </div>
                      </DialogContent>
                    </Dialog>
                  </div>
                </>
              ) : (
                <div className="flex flex-col items-center justify-center py-8 px-4">
                  <div className="bg-gray-700 rounded-full p-4 mb-4">
                    <Wallet className="h-8 w-8 text-gray-400" />
                  </div>
                  <h3 className="text-lg font-medium text-gray-300 mb-2">No Wallet Found</h3>
                  <p className="text-gray-400 text-center mb-4">
                    You need to create a wallet in the key session to use this feature.
                  </p>
                  {walletError && (
                    <div className="mb-4 p-3 bg-red-900/30 border border-red-700 rounded-lg text-red-400 text-sm">
                      {walletError}
                    </div>
                  )}
                  <Button 
                    onClick={handleCreateWallet}
                    className="bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white shadow-md hover:shadow-lg transition-all duration-200 rounded-full"
                  >
                    Go to Key Session
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Blockchain Info Card */}
          <Card className="border border-gray-700 bg-gray-800 shadow-md rounded-xl overflow-hidden">
            <div className="bg-gradient-to-r from-purple-600 to-pink-600 h-1"></div>
            <CardHeader className="pb-2">
              <CardTitle className="text-lg flex items-center gap-2 text-gray-100">
                <Sparkles className="h-5 w-5 text-purple-400" />
                Blockchain Technology
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="bg-gray-700 rounded-xl p-4 border border-gray-600 relative overflow-hidden group hover:shadow-md transition-all duration-200">
                  <div className="absolute top-0 right-0 w-16 h-16 bg-purple-900/30 rounded-bl-full opacity-30 group-hover:opacity-50 transition-opacity"></div>
                  <h3 className="font-medium text-purple-300 mb-1 flex items-center gap-2">
                    <Shield className="h-4 w-4" />
                    Swan Chain
                  </h3>
                  <p className="text-sm text-gray-300">
                    Our proprietary blockchain used for secure record-keeping and data integrity.
                  </p>
                  <Button
                    variant="link"
                    size="sm"
                    className="text-purple-300 p-0 h-auto mt-2 text-xs flex items-center"
                  >
                    Learn more <ExternalLink className="h-3 w-3 ml-1" />
                  </Button>
                </div>

                <div className="bg-gray-700 rounded-xl p-4 border border-gray-600 relative overflow-hidden group hover:shadow-md transition-all duration-200">
                  <div className="absolute top-0 right-0 w-16 h-16 bg-blue-900/30 rounded-bl-full opacity-30 group-hover:opacity-50 transition-opacity"></div>
                  <h3 className="font-medium text-blue-300 mb-1 flex items-center gap-2">
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
                  <p className="text-sm text-gray-300">
                    Fast, secure, and low-cost blockchain used for all payment transactions.
                  </p>
                  <Button
                    variant="link"
                    size="sm"
                    className="text-blue-300 p-0 h-auto mt-2 text-xs flex items-center"
                  >
                    Learn more <ExternalLink className="h-3 w-3 ml-1" />
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
