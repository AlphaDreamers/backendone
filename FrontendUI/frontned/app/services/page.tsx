'use client';

import React, { useState } from 'react';
import Image from 'next/image';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { 
  BarChart3, 
  Wallet, 
  Users, 
  Settings, 
  Bell, 
  Search,
  ChevronDown,
  Package,
  Star,
  Tag,
  ExternalLink,
  X,
  ShoppingCart,
  User,
  LogOut,
  CreditCard,
  Coins,
  MessageSquare
} from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';

// Updated service data structure
const services = [
  {
    id: "1",
    name: "DeFi Development",
    description: "Custom DeFi protocol development with smart contract integration",
    category: "Blockchain Development",
    price: 2500.00,
    createdAt: "2024-03-15T10:00:00Z",
    updatedAt: "2024-03-15T10:00:00Z",
    userId: "0x1234...5678",
    status: "Active",
    image: "https://picsum.photos/400/200?random=1",
    icon: Wallet,
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
    userId: "0x8765...4321",
    status: "Active",
    image: "https://picsum.photos/400/200?random=2",
    icon: Package,
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
    userId: "0xabcd...efgh",
    status: "Active",
    image: "https://picsum.photos/400/200?random=3",
    icon: BarChart3,
    rating: 4.7,
    reviews: 89,
    offer: "1 month support included"
  },
  {
    id: "4",
    name: "Staking Protocol",
    description: "Custom staking protocol development with rewards distribution",
    category: "DeFi Development",
    price: 3000.00,
    createdAt: "2024-03-12T14:20:00Z",
    updatedAt: "2024-03-12T14:20:00Z",
    userId: "0x2468...1357",
    status: "Active",
    image: "https://picsum.photos/400/200?random=4",
    icon: Settings,
    rating: 4.9,
    reviews: 234,
    offer: "Free audit report"
  },
  {
    id: "5",
    name: "Cross-Chain Bridge",
    description: "Development of cross-chain bridge protocol",
    category: "Blockchain Development",
    price: 4000.00,
    createdAt: "2024-03-11T11:45:00Z",
    updatedAt: "2024-03-11T11:45:00Z",
    userId: "0x1357...2468",
    status: "Active",
    image: "https://picsum.photos/400/200?random=5",
    icon: Users,
    rating: 4.6,
    reviews: 167,
    offer: "3 months maintenance"
  },
  {
    id: "6",
    name: "Analytics Dashboard",
    description: "Custom blockchain analytics dashboard development",
    category: "Web Development",
    price: 2000.00,
    createdAt: "2024-03-10T16:30:00Z",
    updatedAt: "2024-03-10T16:30:00Z",
    userId: "0x9876...5432",
    status: "Active",
    image: "https://picsum.photos/400/200?random=6",
    icon: BarChart3,
    rating: 4.8,
    reviews: 123,
    offer: "Free data API access"
  }
];

// Mock review data
const mockReviews: Record<string, Array<{
  id: number;
  user: string;
  rating: number;
  comment: string;
  date: string;
  transactionId: string;
  verified: boolean;
}>> = {
  "1": [
    {
      id: 1,
      user: "0x1234...5678",
      rating: 5,
      comment: "Excellent DeFi service with great APY rates!",
      date: "2024-03-15",
      transactionId: "0x8a7b...9c4d",
      verified: true
    },
    {
      id: 2,
      user: "0x8765...4321",
      rating: 4,
      comment: "Good service but could use more features",
      date: "2024-03-14",
      transactionId: "0x3f2e...1d0c",
      verified: true
    }
  ],
  "2": [],
  "3": [
    {
      id: 1,
      user: "0xabcd...efgh",
      rating: 5,
      comment: "Fast and efficient token swaps",
      date: "2024-03-13",
      transactionId: "0x7d6c...5b4a",
      verified: true
    }
  ],
  "4": [
    {
      id: 1,
      user: "0x2468...1357",
      rating: 4,
      comment: "Staking platform works well",
      date: "2024-03-12",
      transactionId: "0x9e8d...7c6b",
      verified: true
    }
  ],
  "5": [
    {
      id: 1,
      user: "0x1357...2468",
      rating: 3,
      comment: "Bridge service needs improvement",
      date: "2024-03-11",
      transactionId: "0x2c1d...0e9f",
      verified: true
    }
  ],
  "6": [
    {
      id: 1,
      user: "0x9876...5432",
      rating: 5,
      comment: "Best analytics dashboard I've used",
      date: "2024-03-10",
      transactionId: "0x5a4b...3c2d",
      verified: true
    }
  ]
};

// Mock chat history data
const mockChatHistory = [
  {
    id: 1,
    user: "0x1234...5678",
    lastMessage: "When will the DeFi service be ready?",
    timestamp: "2024-03-15T10:30:00Z",
    unread: true
  },
  {
    id: 2,
    user: "0x8765...4321",
    lastMessage: "Thanks for the NFT contract development!",
    timestamp: "2024-03-14T15:45:00Z",
    unread: false
  },
  {
    id: 3,
    user: "0xabcd...efgh",
    lastMessage: "Can you explain the token swap integration?",
    timestamp: "2024-03-13T09:20:00Z",
    unread: true
  },
  {
    id: 4,
    user: "0x2468...1357",
    lastMessage: "The staking platform is working great!",
    timestamp: "2024-03-12T14:10:00Z",
    unread: false
  }
];

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

export default function ServicesPage() {
  const [selectedService, setSelectedService] = useState<string | null>(null);
  const [isReviewModalOpen, setIsReviewModalOpen] = useState(false);
  const [isChatModalOpen, setIsChatModalOpen] = useState(false);
  const [selectedPaymentMethod, setSelectedPaymentMethod] = useState<string | null>(null);
  const router = useRouter();

  const handleReviewClick = (serviceId: string) => {
    setSelectedService(serviceId);
    setIsReviewModalOpen(true);
  };

  const handleChatClick = () => {
    setIsChatModalOpen(true);
  };

  const handleOrderClick = (service: typeof services[0], paymentMethod: string) => {
    // Handle order placement with selected payment method
    console.log('Ordering service:', service.name, 'with payment method:', paymentMethod);
    
    // Show success message
    toast.success(`Order placed with ${paymentMethod} payment method`);
  };

  const handleLogout = async () => {
    try {
      // Get the token from localStorage or wherever it's stored
      const token = localStorage.getItem('token') || '';
      
      const response = await fetch('http://localhost:8085/api/auth/logout', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: 'swanhtetaungp@gmail.com'
        }),
        credentials: 'include', // This is important for cookies
      });

      if (response.ok) {
        // Clear local storage
        localStorage.removeItem('token');
        
        // Show success message
        toast.success('Logged out successfully');
        
        // Redirect to login page
        router.push('/login');
      } else {
        toast.error('Failed to logout. Please try again.');
      }
    } catch (error) {
      console.error('Logout error:', error);
      toast.error('An error occurred during logout');
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const formatTime = (dateString: string) => {
    return new Date(dateString).toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(price);
  };

  const getReviewsForService = (serviceId: string) => {
    return mockReviews[serviceId as keyof typeof mockReviews] || [];
  };

  return (
    <div className="min-h-screen bg-[#0A0A0A]">
      {/* Dashboard Header */}
      <header className="bg-[#111111] border-b border-[#1E1E1E]">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <h1 className="text-2xl font-bold text-white">Services</h1>
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                <input
                  type="text"
                  placeholder="Search services..."
                  className="pl-10 pr-4 py-2 bg-[#1E1E1E] border border-[#2E2E2E] rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent text-white placeholder-gray-400"
                />
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <Button 
                variant="ghost" 
                size="icon" 
                className="text-gray-400 hover:text-white hover:bg-[#1E1E1E] relative"
                onClick={handleChatClick}
              >
                <MessageSquare className="h-5 w-5" />
                {mockChatHistory.some(chat => chat.unread) && (
                  <span className="absolute top-0 right-0 w-2 h-2 bg-red-500 rounded-full"></span>
                )}
              </Button>
              <Button variant="ghost" size="icon" className="text-gray-400 hover:text-white hover:bg-[#1E1E1E]">
                <Bell className="h-5 w-5" />
              </Button>
              <div className="flex items-center space-x-2">
                <div className="w-8 h-8 rounded-full bg-[#1E1E1E] flex items-center justify-center">
                  <Users className="h-5 w-5 text-gray-400" />
                </div>
                <ChevronDown className="h-4 w-4 text-gray-400" />
              </div>
              <Button 
                variant="ghost" 
                size="sm" 
                className="text-gray-400 hover:text-white hover:bg-[#1E1E1E]"
                onClick={handleLogout}
              >
                <LogOut className="h-4 w-4 mr-2" />
                Logout
              </Button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {services.map((service) => {
            const Icon = service.icon;
            return (
              <Card key={service.id} className="overflow-hidden hover:shadow-lg transition-shadow duration-200 bg-[#111111] border-[#1E1E1E]">
                <div className="relative h-48 w-full">
                  <Image
                    src={service.image}
                    alt={service.name}
                    fill
                    className="object-cover"
                  />
                  <div className="absolute top-4 right-4">
                    <span className={`px-3 py-1 rounded-full text-xs font-medium text-white ${getStatusColor(service.status)}`}>
                      {service.status}
                    </span>
                  </div>
                </div>
                <CardHeader>
                  <div className="flex items-center space-x-3">
                    <div className="p-2 bg-[#1E1E1E] rounded-lg">
                      <Icon className="h-5 w-5 text-blue-500" />
                    </div>
                    <div>
                      <CardTitle className="text-xl text-white">{service.name}</CardTitle>
                      <div className="flex items-center space-x-2 mt-1">
                        <User className="h-4 w-4 text-gray-400" />
                        <span className="text-sm text-gray-400">{service.userId}</span>
                      </div>
                    </div>
                  </div>
                  <CardDescription className="mt-2 text-gray-400">{service.description}</CardDescription>
                  <div className="flex items-center justify-between mt-2">
                    <span className="text-sm text-gray-400">{service.category}</span>
                    <span className="text-lg font-bold text-white">{formatPrice(service.price)}</span>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {/* Reviews Section */}
                    <div className="flex items-center justify-between">
                      <Button 
                        variant="ghost" 
                        className="flex items-center space-x-2 text-gray-400 hover:text-white"
                        onClick={() => handleReviewClick(service.id)}
                      >
                        <div className="flex items-center">
                          <Star className="h-4 w-4 text-yellow-500 fill-yellow-500" />
                          <span className="ml-1 text-sm font-medium text-white">{service.rating}</span>
                        </div>
                        <span className="text-sm text-gray-400">({service.reviews} reviews)</span>
                      </Button>
                      <div className="text-sm text-gray-400">
                        Posted {formatDate(service.createdAt)}
                      </div>
                    </div>
                    
                    {/* Offer Section */}
                    <div className="flex items-center space-x-2 p-2 bg-[#1E1E1E] rounded-lg">
                      <Tag className="h-4 w-4 text-blue-500" />
                      <span className="text-sm text-blue-500 font-medium">{service.offer}</span>
                    </div>
                    
                    {/* Payment Methods */}
                    <div className="space-y-2">
                      <div className="text-sm text-gray-400">Payment Methods:</div>
                      <div className="flex flex-wrap gap-2">
                        <Button 
                          variant="outline" 
                          size="sm" 
                          className="flex items-center space-x-1 bg-[#1E1E1E] border-[#2E2E2E] text-gray-400 hover:text-white"
                          onClick={() => handleOrderClick(service, "Fiat")}
                        >
                          <CreditCard className="h-3 w-3" />
                          <span>Fiat</span>
                        </Button>
                        <Button 
                          variant="outline" 
                          size="sm" 
                          className="flex items-center space-x-1 bg-[#1E1E1E] border-[#2E2E2E] text-gray-400 hover:text-white"
                          onClick={() => handleOrderClick(service, "ETH")}
                        >
                          <Coins className="h-3 w-3" />
                          <span>ETH</span>
                        </Button>
                        <Button 
                          variant="outline" 
                          size="sm" 
                          className="flex items-center space-x-1 bg-[#1E1E1E] border-[#2E2E2E] text-gray-400 hover:text-white"
                          onClick={() => handleOrderClick(service, "SOL")}
                        >
                          <Coins className="h-3 w-3" />
                          <span>SOL</span>
                        </Button>
                      </div>
                    </div>
                    
                    {/* Action Button */}
                    <Button 
                      variant="outline" 
                      size="sm" 
                      className="w-full hover:bg-[#1E1E1E] text-gray-400 border-[#2E2E2E]"
                    >
                      Learn More
                    </Button>
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>
      </main>

      {/* Review Modal */}
      <Dialog open={isReviewModalOpen} onOpenChange={setIsReviewModalOpen}>
        <DialogContent className="bg-[#111111] border-[#1E1E1E] text-white max-w-2xl">
          <DialogHeader>
            <DialogTitle className="text-xl font-bold">
              {selectedService && services.find(s => s.id === selectedService)?.name} Reviews
            </DialogTitle>
          </DialogHeader>
          <div className="mt-4 space-y-4">
            {selectedService && getReviewsForService(selectedService).length > 0 ? (
              getReviewsForService(selectedService).map((review) => (
                <div key={review.id} className="p-4 bg-[#1E1E1E] rounded-lg">
                  <div className="flex justify-between items-start">
                    <div>
                      <div className="flex items-center space-x-2">
                        <span className="font-medium text-white">{review.user}</span>
                        <div className="flex items-center">
                          {[...Array(review.rating)].map((_, i) => (
                            <Star key={i} className="h-4 w-4 text-yellow-500 fill-yellow-500" />
                          ))}
                        </div>
                      </div>
                      <p className="mt-2 text-gray-400">{review.comment}</p>
                      <div className="mt-2 text-sm text-gray-500">{review.date}</div>
                    </div>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="text-blue-500 hover:text-blue-400"
                      onClick={() => window.open(`https://etherscan.io/tx/${review.transactionId}`, '_blank')}
                    >
                      <ExternalLink className="h-4 w-4 mr-1" />
                      View on Chain
                    </Button>
                  </div>
                  <div className="mt-2 text-xs text-gray-500">
                    Transaction ID: {review.transactionId}
                  </div>
                </div>
              ))
            ) : (
              <div className="text-center text-gray-400 py-8">
                No reviews yet
              </div>
            )}
          </div>
        </DialogContent>
      </Dialog>

      {/* Chat History Modal */}
      <Dialog open={isChatModalOpen} onOpenChange={setIsChatModalOpen}>
        <DialogContent className="bg-[#111111] border-[#1E1E1E] text-white max-w-2xl">
          <DialogHeader>
            <DialogTitle className="text-xl font-bold">Chat History</DialogTitle>
          </DialogHeader>
          <div className="mt-4 space-y-4">
            {mockChatHistory.length > 0 ? (
              mockChatHistory.map((chat) => (
                <div 
                  key={chat.id} 
                  className={`p-4 rounded-lg ${chat.unread ? 'bg-[#1E1E1E] border-l-4 border-blue-500' : 'bg-[#0A0A0A]'}`}
                >
                  <div className="flex justify-between items-start">
                    <div>
                      <div className="flex items-center space-x-2">
                        <User className="h-4 w-4 text-gray-400" />
                        <span className="font-medium text-white">{chat.user}</span>
                        {chat.unread && (
                          <span className="px-2 py-0.5 bg-blue-500 text-white text-xs rounded-full">New</span>
                        )}
                      </div>
                      <p className="mt-2 text-gray-400">{chat.lastMessage}</p>
                    </div>
                    <div className="text-xs text-gray-500">
                      {formatTime(chat.timestamp)}
                    </div>
                  </div>
                  <div className="mt-2 text-xs text-gray-500">
                    {formatDate(chat.timestamp)}
                  </div>
                </div>
              ))
            ) : (
              <div className="text-center text-gray-400 py-8">
                No chat history
              </div>
            )}
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
} 