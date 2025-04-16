"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Connection, PublicKey, LAMPORTS_PER_SOL, clusterApiUrl, Keypair, SystemProgram, Transaction } from "@solana/web3.js"
import { ethers } from "ethers"
import { QRCodeSVG } from "qrcode.react"
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
  QrCode,
  Package,
  Star,
  Tag,
  TrendingUp,
  TrendingDown,
  DollarSign,
  Percent,
  Users,
  Activity,
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
import { getUserProfile, getToken, getUserEmailFromToken, createWallet } from "@/app/lib/auth-utils"

// Mock user services data
const userServices = [
  {
    id: "1",
    name: "DeFi Development",
    description: "Custom DeFi protocol development with smart contract integration",
    category: "Blockchain Development",
    price: 2500.00,
    createdAt: "2024-03-15T10:00:00Z",
    updatedAt: "2024-03-15T10:00:00Z",
    status: "Active",
    rating: 4.8,
    reviews: 342,
    offer: "20% discount on first project"
  },
  {
    id: "2",
    name: "NFT Smart Contract",
    description: "ERC-721/ERC-1155 smart contract development with metadata",
    category: "Smart Contracts",
    price: 1800.00,
    createdAt: "2024-03-14T15:30:00Z",
    updatedAt: "2024-03-14T15:30:00Z",
    status: "Active",
    rating: 4.5,
    reviews: 156,
    offer: "Free deployment on testnet"
  },
  {
    id: "3",
    name: "Token Swap Integration",
    description: "Integration of DEX protocols for token swapping",
    category: "DeFi Integration",
    price: 1200.00,
    createdAt: "2024-03-13T09:15:00Z",
    updatedAt: "2024-03-13T09:15:00Z",
    status: "Active",
    rating: 4.7,
    reviews: 89,
    offer: "1 month support included"
  }
];

// Mock transactions data
const transactions = [
  {
    id: "tx1",
    type: "receive",
    amount: 2.5,
    from: "8Kv...j3M",
    to: "Your wallet",
    date: "Today, 10:23 AM",
    status: "completed",
    network: "solana"
  },
  {
    id: "tx2",
    type: "send",
    amount: 0.75,
    from: "Your wallet",
    to: "3xR...p7B",
    date: "Yesterday, 3:45 PM",
    status: "completed",
    network: "solana"
  },
  {
    id: "tx3",
    type: "receive",
    amount: 1.2,
    from: "5Gs...k9P",
    to: "Your wallet",
    date: "Apr 5, 2023",
    status: "completed",
    network: "ethereum"
  },
];

export default function Dashboard() {
  const router = useRouter()
  const [solanaBalance, setSolanaBalance] = useState<number>(0)
  const [ethereumBalance, setEthereumBalance] = useState<string>("0")
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [solanaWalletAddress, setSolanaWalletAddress] = useState<string>("")
  const [ethereumWalletAddress, setEthereumWalletAddress] = useState<string>("")
  const [copied, setCopied] = useState<boolean>(false)
  const [recipientAddress, setRecipientAddress] = useState<string>("")
  const [amount, setAmount] = useState<string>("")
  const [showBanner, setShowBanner] = useState<boolean>(true)
  const [userProfile, setUserProfile] = useState<any>(null)
  const [profileError, setProfileError] = useState<string>("")
  const [activeWallet, setActiveWallet] = useState<"solana" | "ethereum">("solana")
  const [transactions, setTransactions] = useState<any[]>([
    {
      id: "tx1",
      type: "receive",
      amount: 2.5,
      from: "8Kv...j3M",
      to: "Your wallet",
      date: "Today, 10:23 AM",
      status: "completed",
      network: "solana"
    },
    {
      id: "tx2",
      type: "send",
      amount: 0.75,
      from: "Your wallet",
      to: "3xR...p7B",
      date: "Yesterday, 3:45 PM",
      status: "completed",
      network: "solana"
    },
    {
      id: "tx3",
      type: "receive",
      amount: 1.2,
      from: "5Gs...k9P",
      to: "Your wallet",
      date: "Apr 5, 2023",
      status: "completed",
      network: "ethereum"
    },
  ])
  const [isCreatingWallet, setIsCreatingWallet] = useState<boolean>(false)
  const [walletError, setWalletError] = useState<string>("")
  const [activeTab, setActiveTab] = useState("all")
  
  // Blockchain data states
  const [solanaPrice, setSolanaPrice] = useState<number>(0)
  const [ethereumPrice, setEthereumPrice] = useState<number>(0)
  const [solanaChange24h, setSolanaChange24h] = useState<number>(0)
  const [ethereumChange24h, setEthereumChange24h] = useState<number>(0)
  const [solanaMarketCap, setSolanaMarketCap] = useState<number>(0)
  const [ethereumMarketCap, setEthereumMarketCap] = useState<number>(0)
  const [solanaVolume24h, setSolanaVolume24h] = useState<number>(0)
  const [ethereumVolume24h, setEthereumVolume24h] = useState<number>(0)
  const [solanaTransactions24h, setSolanaTransactions24h] = useState<number>(0)
  const [ethereumTransactions24h, setEthereumTransactions24h] = useState<number>(0)
  const [solanaActiveAddresses, setSolanaActiveAddresses] = useState<number>(0)
  const [ethereumActiveAddresses, setEthereumActiveAddresses] = useState<number>(0)
  const [solanaNetworkTps, setSolanaNetworkTps] = useState<number>(0)
  const [ethereumNetworkTps, setEthereumNetworkTps] = useState<number>(0)

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
    
    // Fetch blockchain data
    fetchBlockchainData();
    
    // Set up interval to refresh blockchain data every 30 seconds
    const interval = setInterval(fetchBlockchainData, 30000);
    
    return () => clearInterval(interval);
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
      const provider = new ethers.JsonRpcProvider("https://eth-sepolia.g.alchemy.com/v2/demo");
      const balance = await provider.getBalance(address);
      setEthereumBalance(ethers.formatEther(balance));
    } catch (error) {
      console.error("Error fetching Ethereum balance:", error);
      setEthereumBalance("0");
    }
  };

  const fetchBlockchainData = async () => {
    try {
      // Fetch Solana data
      const solanaConnection = new Connection(clusterApiUrl("mainnet-beta"), "confirmed");
      
      // Get recent performance samples to calculate TPS
      const performanceSamples = await solanaConnection.getRecentPerformanceSamples(1);
      if (performanceSamples.length > 0) {
        setSolanaNetworkTps(performanceSamples[0].numTransactions / performanceSamples[0].samplePeriodSecs);
      }
      
      // Mock data for Solana price and market data (in a real app, you would use an API)
      setSolanaPrice(150.25);
      setSolanaChange24h(5.2);
      setSolanaMarketCap(65000000000);
      setSolanaVolume24h(2500000000);
      setSolanaTransactions24h(35000000);
      setSolanaActiveAddresses(1200000);
      
      // Mock data for Ethereum price and market data
      setEthereumPrice(3000.75);
      setEthereumChange24h(-2.1);
      setEthereumMarketCap(350000000000);
      setEthereumVolume24h(15000000000);
      setEthereumTransactions24h(1200000);
      setEthereumActiveAddresses(800000);
      setEthereumNetworkTps(15);
      
    } catch (error) {
      console.error("Error fetching blockchain data:", error);
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
        
        // Add to transactions for display
        const newTx = {
          id: `tx${Date.now()}`,
          type: "send",
          amount: Number.parseFloat(amount),
          from: formatAddress(fromKeypair.publicKey.toString()),
          to: formatAddress(recipientAddress),
          date: "Just now",
          status: "completed",
          network: "solana"
        };
        
        setTransactions([newTx, ...transactions]);
        
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
        
        // Add to transactions for display
        const newTx = {
          id: `tx${Date.now()}`,
          type: "send",
          amount: Number.parseFloat(amount),
          from: formatAddress(wallet.address),
          to: formatAddress(recipientAddress),
          date: "Just now",
          status: "completed",
          network: "ethereum"
        };
        
        setTransactions([newTx, ...transactions]);
        
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

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(value);
  };

  const formatLargeNumber = (value: number) => {
    if (value >= 1000000000) {
      return (value / 1000000000).toFixed(2) + 'B';
    } else if (value >= 1000000) {
      return (value / 1000000).toFixed(2) + 'M';
    } else if (value >= 1000) {
      return (value / 1000).toFixed(2) + 'K';
    } else {
      return value.toString();
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(price);
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'active':
        return 'bg-green-500';
      case 'coming soon':
        return 'bg-yellow-500';
      case 'beta':
        return 'bg-blue-500';
      case 'maintenance':
        return 'bg-red-500';
      default:
        return 'bg-gray-500';
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 to-gray-800 text-gray-100">
      <div className="flex flex-col h-screen">
        {/* Banner */}
        {showBanner && (
          <div className="bg-gradient-to-r from-blue-900 to-purple-900 text-white py-3 px-4 relative">
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
        <header className="bg-gray-800 shadow-md border-b border-gray-700 py-3">
          <div className="container mx-auto px-4 flex justify-between items-center">
            <div className="flex items-center gap-2">
              <div className="bg-gradient-to-r from-blue-600 to-purple-600 p-2 rounded-lg">
                <Wallet className="h-5 w-5 text-white" />
              </div>
              <span className="font-bold text-xl bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
                BlockChain Service Offering Platform
              </span>
            </div>

            <div className="flex-1 max-w-md mx-4">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="Search transactions..."
                  className="pl-10 bg-gray-700 border-gray-600 text-gray-100 placeholder:text-gray-400 focus:border-blue-500 focus:ring-blue-500 rounded-full"
                />
              </div>
            </div>

            <div className="flex items-center gap-4">
              <Button
                variant="ghost"
                size="icon"
                className="text-gray-300 hover:text-blue-400 hover:bg-gray-700 rounded-full"
              >
                <Bell className="h-5 w-5" />
              </Button>

              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    className="flex items-center gap-2 text-gray-300 hover:text-blue-400 hover:bg-gray-700 rounded-full"
                  >
                    <Avatar className="h-8 w-8 border-2 border-blue-500">
                      <AvatarImage src="/placeholder-user.jpg" />
                      <AvatarFallback className="bg-gradient-to-r from-blue-500 to-purple-500 text-white">
                        {userProfile?.user_name?.charAt(0) || "U"}
                      </AvatarFallback>
                    </Avatar>
                    <span className="hidden md:inline-block font-medium">{userProfile?.user_name || "User"}</span>
                    <ChevronDown className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56 bg-gray-800 border-gray-700 text-gray-100">
                  <DropdownMenuLabel>My Account</DropdownMenuLabel>
                  <DropdownMenuSeparator className="bg-gray-700" />
                  <DropdownMenuItem className="hover:bg-gray-700 cursor-pointer">
                    <User className="mr-2 h-4 w-4 text-blue-400" />
                    <span>Profile</span>
                  </DropdownMenuItem>
                  <DropdownMenuItem className="hover:bg-gray-700 cursor-pointer">
                    <Settings className="mr-2 h-4 w-4 text-blue-400" />
                    <span>Settings</span>
                  </DropdownMenuItem>
                  <DropdownMenuSeparator className="bg-gray-700" />
                  <DropdownMenuItem className="hover:bg-gray-700 cursor-pointer" onClick={() => router.push("/")}>
                    <LogOut className="mr-2 h-4 w-4 text-blue-400" />
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

              {/* Blockchain Overview Card */}
              <Card className="border border-gray-700 bg-gray-800 shadow-lg rounded-xl overflow-hidden">
                <div className="bg-gradient-to-r from-blue-600 to-purple-600 h-3"></div>
                <CardHeader className="pb-2">
                  <CardTitle className="text-xl flex items-center gap-2 text-gray-100">
                    <Activity className="h-5 w-5 text-blue-400" />
                    Blockchain Overview
                  </CardTitle>
                  <CardDescription className="text-gray-400">
                    Real-time data from Solana and Ethereum networks
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {/* Solana Data */}
                    <div className="bg-gray-700 rounded-xl p-4 border border-gray-600">
                      <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center gap-2">
                          <svg className="h-5 w-5 text-blue-400" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
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
                          <h3 className="font-medium text-white">Solana</h3>
                        </div>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={fetchBlockchainData}
                          className="text-gray-400 hover:text-blue-400 hover:bg-gray-600 rounded-full"
                        >
                          <RefreshCw className="h-4 w-4" />
                        </Button>
                      </div>
                      
                      <div className="space-y-4">
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">Price</span>
                          <span className="font-medium text-white">${solanaPrice.toFixed(2)}</span>
                        </div>
                        
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">24h Change</span>
                          <div className={`flex items-center ${solanaChange24h >= 0 ? 'text-green-400' : 'text-red-400'}`}>
                            {solanaChange24h >= 0 ? <TrendingUp className="h-4 w-4 mr-1" /> : <TrendingDown className="h-4 w-4 mr-1" />}
                            <span>{Math.abs(solanaChange24h).toFixed(2)}%</span>
                          </div>
                        </div>
                        
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">Market Cap</span>
                          <span className="font-medium text-white">{formatCurrency(solanaMarketCap)}</span>
                        </div>
                        
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">24h Volume</span>
                          <span className="font-medium text-white">{formatCurrency(solanaVolume24h)}</span>
                        </div>
                        
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">24h Transactions</span>
                          <span className="font-medium text-white">{formatLargeNumber(solanaTransactions24h)}</span>
                        </div>
                        
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">Active Addresses</span>
                          <span className="font-medium text-white">{formatLargeNumber(solanaActiveAddresses)}</span>
                        </div>
                        
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">Network TPS</span>
                          <span className="font-medium text-white">{solanaNetworkTps.toFixed(2)}</span>
                        </div>
                      </div>
                    </div>
                    
                    {/* Ethereum Data */}
                    <div className="bg-gray-700 rounded-xl p-4 border border-gray-600">
                      <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center gap-2">
                          <svg className="h-5 w-5 text-blue-400" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                            <path d="M12 6V18" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                            <path d="M8 10L12 6L16 10" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                            <path d="M8 14L12 18L16 14" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                          </svg>
                          <h3 className="font-medium text-white">Ethereum</h3>
                        </div>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={fetchBlockchainData}
                          className="text-gray-400 hover:text-blue-400 hover:bg-gray-600 rounded-full"
                        >
                          <RefreshCw className="h-4 w-4" />
                        </Button>
                      </div>
                      
                      <div className="space-y-4">
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">Price</span>
                          <span className="font-medium text-white">${ethereumPrice.toFixed(2)}</span>
                        </div>
                        
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">24h Change</span>
                          <div className={`flex items-center ${ethereumChange24h >= 0 ? 'text-green-400' : 'text-red-400'}`}>
                            {ethereumChange24h >= 0 ? <TrendingUp className="h-4 w-4 mr-1" /> : <TrendingDown className="h-4 w-4 mr-1" />}
                            <span>{Math.abs(ethereumChange24h).toFixed(2)}%</span>
                          </div>
                        </div>
                        
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">Market Cap</span>
                          <span className="font-medium text-white">{formatCurrency(ethereumMarketCap)}</span>
                        </div>
                        
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">24h Volume</span>
                          <span className="font-medium text-white">{formatCurrency(ethereumVolume24h)}</span>
                        </div>
                        
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">24h Transactions</span>
                          <span className="font-medium text-white">{formatLargeNumber(ethereumTransactions24h)}</span>
                        </div>
                        
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">Active Addresses</span>
                          <span className="font-medium text-white">{formatLargeNumber(ethereumActiveAddresses)}</span>
                        </div>
                        
                        <div className="flex justify-between items-center">
                          <span className="text-gray-400">Network TPS</span>
                          <span className="font-medium text-white">{ethereumNetworkTps.toFixed(2)}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Transactions */}
              <Card className="border border-gray-700 bg-gray-800 shadow-md rounded-xl overflow-hidden">
                <div className="bg-gradient-to-r from-blue-600 to-purple-600 h-1"></div>
                <CardHeader>
                  <CardTitle className="text-lg flex items-center gap-2 text-gray-100">
                    <History className="h-5 w-5 text-blue-400" />
                    Recent Transactions
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <Tabs defaultValue="all">
                    <TabsList className="bg-gray-700 mb-4 p-1 rounded-lg">
                      <TabsTrigger
                        value="all"
                        className="data-[state=active]:bg-gray-800 data-[state=active]:text-blue-400 data-[state=active]:shadow-sm rounded-md text-gray-300"
                      >
                        All
                      </TabsTrigger>
                      <TabsTrigger
                        value="sent"
                        className="data-[state=active]:bg-gray-800 data-[state=active]:text-blue-400 data-[state=active]:shadow-sm rounded-md text-gray-300"
                      >
                        Sent
                      </TabsTrigger>
                      <TabsTrigger
                        value="received"
                        className="data-[state=active]:bg-gray-800 data-[state=active]:text-blue-400 data-[state=active]:shadow-sm rounded-md text-gray-300"
                      >
                        Received
                      </TabsTrigger>
                      <TabsTrigger
                        value="services"
                        className="data-[state=active]:bg-gray-800 data-[state=active]:text-blue-400 data-[state=active]:shadow-sm rounded-md text-gray-300"
                      >
                        My Services
                      </TabsTrigger>
                    </TabsList>

                    <TabsContent value="all" className="space-y-4">
                      {transactions.map((tx) => (
                        <div
                          key={tx.id}
                          className="flex items-center justify-between p-3 rounded-xl bg-gray-700 hover:bg-gray-600 transition-colors border border-gray-600 group hover:shadow-sm"
                        >
                          <div className="flex items-center gap-3">
                            <div
                              className={`p-2 rounded-full ${tx.type === "receive" ? "bg-green-900/30 text-green-400" : "bg-blue-900/30 text-blue-400"} group-hover:scale-110 transition-transform`}
                            >
                              {tx.type === "receive" ? (
                                <ArrowDownLeft className="h-5 w-5" />
                              ) : (
                                <ArrowUpRight className="h-5 w-5" />
                              )}
                            </div>
                            <div>
                              <div className="font-medium text-gray-100">{tx.type === "receive" ? "Received" : "Sent"} {tx.network === "solana" ? "SOL" : "ETH"}</div>
                              <div className="text-sm text-gray-400">{tx.date}</div>
                            </div>
                          </div>
                          <div className="text-right">
                            <div
                              className={`font-medium ${tx.type === "receive" ? "text-green-400" : "text-blue-400"}`}
                            >
                              {tx.type === "receive" ? "+" : "-"}
                              {tx.amount} {tx.network === "solana" ? "SOL" : "ETH"}
                            </div>
                            <div className="text-xs text-gray-400">
                              {tx.status === "completed" ? (
                                <span className="inline-flex items-center text-green-400">
                                  <Check className="h-3 w-3 mr-1" /> Completed
                                </span>
                              ) : (
                                <span className="inline-flex items-center text-orange-400">
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
                            className="flex items-center justify-between p-3 rounded-xl bg-gray-700 hover:bg-gray-600 transition-colors border border-gray-600 group hover:shadow-sm"
                          >
                            <div className="flex items-center gap-3">
                              <div className="p-2 rounded-full bg-blue-900/30 text-blue-400 group-hover:scale-110 transition-transform">
                                <ArrowUpRight className="h-5 w-5" />
                              </div>
                              <div>
                                <div className="font-medium text-gray-100">Sent {tx.network === "solana" ? "SOL" : "ETH"}</div>
                                <div className="text-sm text-gray-400">{tx.date}</div>
                              </div>
                            </div>
                            <div className="text-right">
                              <div className="font-medium text-blue-400">-{tx.amount} {tx.network === "solana" ? "SOL" : "ETH"}</div>
                              <div className="text-xs text-gray-400">
                                {tx.status === "completed" ? (
                                  <span className="inline-flex items-center text-green-400">
                                    <Check className="h-3 w-3 mr-1" /> Completed
                                  </span>
                                ) : (
                                  <span className="inline-flex items-center text-orange-400">
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
                            className="flex items-center justify-between p-3 rounded-xl bg-gray-700 hover:bg-gray-600 transition-colors border border-gray-600 group hover:shadow-sm"
                          >
                            <div className="flex items-center gap-3">
                              <div className="p-2 rounded-full bg-green-900/30 text-green-400 group-hover:scale-110 transition-transform">
                                <ArrowDownLeft className="h-5 w-5" />
                              </div>
                              <div>
                                <div className="font-medium text-gray-100">Received {tx.network === "solana" ? "SOL" : "ETH"}</div>
                                <div className="text-sm text-gray-400">{tx.date}</div>
                              </div>
                            </div>
                            <div className="text-right">
                              <div className="font-medium text-green-400">+{tx.amount} {tx.network === "solana" ? "SOL" : "ETH"}</div>
                              <div className="text-xs text-gray-400">
                                {tx.status === "completed" ? (
                                  <span className="inline-flex items-center text-green-400">
                                    <Check className="h-3 w-3 mr-1" /> Completed
                                  </span>
                                ) : (
                                  <span className="inline-flex items-center text-orange-400">
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

                    <TabsContent value="services" className="space-y-4">
                      <div className="flex justify-between items-center mb-4">
                        <h3 className="text-lg font-medium text-white">My Services</h3>
                        <Button 
                          variant="outline" 
                          size="sm" 
                          className="bg-blue-600 hover:bg-blue-700 text-white border-none"
                          onClick={() => router.push('/services')}
                        >
                          <Package className="h-4 w-4 mr-2" />
                          Add New Service
                        </Button>
                      </div>
                      
                      {userServices.length > 0 ? (
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                          {userServices.map((service) => (
                            <div 
                              key={service.id} 
                              className="p-4 rounded-lg bg-gray-700 hover:bg-gray-600 transition-colors border border-gray-600"
                            >
                              <div className="flex justify-between items-start">
                                <div>
                                  <div className="flex items-center space-x-2">
                                    <h4 className="font-medium text-white">{service.name}</h4>
                                    <span className={`px-2 py-0.5 rounded-full text-xs font-medium text-white ${getStatusColor(service.status)}`}>
                                      {service.status}
                                    </span>
                                  </div>
                                  <p className="mt-1 text-sm text-gray-300">{service.description}</p>
                                </div>
                                <span className="text-lg font-bold text-white">{formatPrice(service.price)}</span>
                              </div>
                              
                              <div className="mt-3 flex items-center justify-between">
                                <div className="flex items-center space-x-2">
                                  <Star className="h-4 w-4 text-yellow-500 fill-yellow-500" />
                                  <span className="text-sm text-gray-300">{service.rating}</span>
                                  <span className="text-sm text-gray-400">({service.reviews} reviews)</span>
                                </div>
                                <div className="text-xs text-gray-400">
                                  Posted {formatDate(service.createdAt)}
                                </div>
                              </div>
                              
                              <div className="mt-3 flex items-center space-x-2 p-2 bg-gray-800 rounded-lg">
                                <Tag className="h-4 w-4 text-blue-400" />
                                <span className="text-sm text-blue-400 font-medium">{service.offer}</span>
                              </div>
                              
                              <div className="mt-3 flex justify-end space-x-2">
                                <Button 
                                  variant="outline" 
                                  size="sm" 
                                  className="text-gray-300 hover:text-white"
                                >
                                  Edit
                                </Button>
                                <Button 
                                  variant="outline" 
                                  size="sm" 
                                  className="text-gray-300 hover:text-white"
                                >
                                  View Orders
                                </Button>
                              </div>
                            </div>
                          ))}
                        </div>
                      ) : (
                        <div className="text-center py-8 text-gray-400">
                          <Package className="h-12 w-12 mx-auto mb-4 text-gray-500" />
                          <p>You haven't created any services yet.</p>
                          <Button 
                            variant="outline" 
                            size="sm" 
                            className="mt-4 bg-blue-600 hover:bg-blue-700 text-white border-none"
                            onClick={() => router.push('/services')}
                          >
                            Create Your First Service
                          </Button>
                        </div>
                      )}
                    </TabsContent>
                  </Tabs>
                </CardContent>
                <CardFooter className="border-t border-gray-700 pt-4 flex justify-center">
                  <Button variant="link" className="text-blue-400">
                    View All Transactions
                  </Button>
                </CardFooter>
              </Card>
            </div>

            {/* Right column - User details and stats */}
            <div className="space-y-6">
              {/* User profile card */}
              <Card className="border border-gray-700 bg-gray-800 shadow-md rounded-xl overflow-hidden">
                <div className="bg-gradient-to-r from-blue-600 to-purple-600 h-1"></div>
                <CardHeader className="pb-2">
                  <CardTitle className="text-lg flex items-center gap-2 text-gray-100">
                    <BarChart3 className="h-5 w-5 text-blue-400" />
                    User Profile
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {isLoading ? (
                    <div className="flex items-center justify-center p-6">
                      <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
                    </div>
                  ) : profileError ? (
                    <div className="p-4 bg-red-100 text-red-700 rounded-md">
                      <p>{profileError}</p>
                    </div>
                  ) : !userProfile ? (
                    <div className="p-4 bg-yellow-100 text-yellow-700 rounded-md">
                      <p>No user profile data available</p>
                    </div>
                  ) : (
                    <div className="space-y-4">
                      <div className="flex items-center space-x-4">
                        <div className="h-16 w-16 rounded-full bg-gradient-to-r from-blue-500 to-purple-500 flex items-center justify-center text-white text-xl font-bold">
                          {userProfile.user_name?.charAt(0) || userProfile.email?.charAt(0) || "U"}
                        </div>
                        <div>
                          <h2 className="text-xl font-bold text-white">{userProfile.user_name || "User"}</h2>
                          <p className="text-gray-400">{userProfile.email}</p>
                        </div>
                      </div>
                      
                      <div className="grid grid-cols-1 gap-4">
                        <div className="p-4 bg-gray-700 rounded-md">
                          <p className="text-sm text-gray-400">Account Status</p>
                          <p className="font-medium text-white">{userProfile.verified ? "Verified" : "Unverified"}</p>
                        </div>
                        <div className="p-4 bg-gray-700 rounded-md">
                          <p className="text-sm text-gray-400">Wallet Status</p>
                          <p className="font-medium text-white">{userProfile.wallet_created ? "Created" : "Not Created"}</p>
                        </div>
                        <div className="p-4 bg-gray-700 rounded-md">
                          <p className="text-sm text-gray-400">Member Since</p>
                          <p className="font-medium text-white">{new Date(userProfile.created_at).toLocaleDateString()}</p>
                        </div>
                        {userProfile.wallet_created_time && (
                          <div className="p-4 bg-gray-700 rounded-md">
                            <p className="text-sm text-gray-400">Wallet Created</p>
                            <p className="font-medium text-white">{new Date(userProfile.wallet_created_time).toLocaleDateString()}</p>
                          </div>
                        )}
                      </div>
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Quick stats */}
              <Card className="border border-gray-700 bg-gray-800 shadow-md rounded-xl overflow-hidden">
                <div className="bg-gradient-to-r from-blue-600 to-purple-600 h-1"></div>
                <CardHeader className="pb-2">
                  <CardTitle className="text-lg flex items-center gap-2 text-gray-100">
                    <BarChart3 className="h-5 w-5 text-blue-400" />
                    Quick Stats
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="bg-gray-700 p-4 rounded-xl border border-gray-600 hover:shadow-md transition-shadow group">
                      <div className="flex items-center gap-2 mb-2">
                        <DollarSign className="h-4 w-4 text-blue-400 group-hover:scale-110 transition-transform" />
                        <span className="text-sm text-blue-400">Total Value</span>
                      </div>
                      <div className="text-xl font-medium text-blue-300">$712.50</div>
                    </div>
                    <div className="bg-gray-700 p-4 rounded-xl border border-gray-600 hover:shadow-md transition-shadow group">
                      <div className="flex items-center gap-2 mb-2">
                        <History className="h-4 w-4 text-purple-400 group-hover:scale-110 transition-transform" />
                        <span className="text-sm text-purple-400">Transactions</span>
                      </div>
                      <div className="text-xl font-medium text-purple-300">24</div>
                    </div>
                    <div className="bg-gray-700 p-4 rounded-xl border border-gray-600 hover:shadow-md transition-shadow group">
                      <div className="flex items-center gap-2 mb-2">
                        <Package className="h-4 w-4 text-pink-400 group-hover:scale-110 transition-transform" />
                        <span className="text-sm text-pink-400">Services</span>
                      </div>
                      <div className="text-xl font-medium text-pink-300">{userServices.length}</div>
                    </div>
                    <div className="bg-gray-700 p-4 rounded-xl border border-gray-600 hover:shadow-md transition-shadow group">
                      <div className="flex items-center gap-2 mb-2">
                        <Users className="h-4 w-4 text-indigo-400 group-hover:scale-110 transition-transform" />
                        <span className="text-sm text-indigo-400">Customers</span>
                      </div>
                      <div className="text-xl font-medium text-indigo-300">12</div>
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
