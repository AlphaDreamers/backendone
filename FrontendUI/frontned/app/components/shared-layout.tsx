'use client';

import React, { useState, useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import { 
  Search, 
  Bell, 
  User, 
  Settings, 
  LogOut, 
  Wallet, 
  BarChart3, 
  History, 
  CreditCard, 
  Package,
  MessageSquare,
  Menu,
  X,
  Home,
  Key,
  ShoppingCart
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { getUserProfile, getToken, getUserEmailFromToken } from "@/app/lib/auth-utils";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

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

interface SharedLayoutProps {
  children: React.ReactNode;
}

export default function SharedLayout({ children }: SharedLayoutProps) {
  const router = useRouter();
  const pathname = usePathname();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isChatModalOpen, setIsChatModalOpen] = useState(false);
  const [userProfile, setUserProfile] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [profileError, setProfileError] = useState<string>("");

  useEffect(() => {
    const fetchUserProfile = async () => {
      try {
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

        const response = await getUserProfile(email, token);
        
        if (response.success && response.data) {
          setUserProfile(response.data);
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
          email: userProfile?.email || ''
        }),
        credentials: 'include', // This is important for cookies
      });

      if (response.ok) {
        // Clear local storage
        localStorage.removeItem('token');
        
        // Redirect to login page
        router.push('/login');
      } else {
        console.error('Failed to logout');
      }
    } catch (error) {
      console.error('Logout error:', error);
    }
  };

  const formatTime = (dateString: string) => {
    return new Date(dateString).toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const isActive = (path: string) => {
    return pathname === path;
  };

  const navItems = [
    { name: 'Home', path: '/', icon: Home },
    { name: 'Dashboard', path: '/dashboard', icon: BarChart3 },
    { name: 'Services', path: '/services', icon: Package },
    { name: 'Wallet', path: '/key', icon: Wallet },
    { name: 'Transactions', path: '/transactions', icon: History },
  ];

  return (
    <div className="min-h-screen bg-[#0A0A0A] flex flex-col">
      {/* Header */}
      <header className="bg-[#111111] border-b border-[#1E1E1E] sticky top-0 z-50">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <Link href="/" className="flex items-center space-x-2">
                <div className="bg-gradient-to-r from-blue-600 to-purple-600 p-2 rounded-lg">
                  <Wallet className="h-5 w-5 text-white" />
                </div>
                <span className="font-bold text-xl bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
                  BlockChain Service
                </span>
              </Link>
              
              {/* Desktop Navigation */}
              <nav className="hidden md:flex items-center space-x-1">
                {navItems.map((item) => {
                  const Icon = item.icon;
                  return (
                    <Link 
                      key={item.path} 
                      href={item.path}
                      className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                        isActive(item.path) 
                          ? 'bg-gray-800 text-blue-400' 
                          : 'text-gray-300 hover:bg-gray-700 hover:text-white'
                      }`}
                    >
                      <span className="flex items-center">
                        <Icon className="h-4 w-4 mr-2" />
                        {item.name}
                      </span>
                    </Link>
                  );
                })}
              </nav>
            </div>
            
            <div className="flex items-center space-x-4">
              {/* Search */}
              <div className="hidden md:block relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                <Input
                  placeholder="Search..."
                  className="pl-10 bg-[#1E1E1E] border-[#2E2E2E] rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent text-white placeholder-gray-400 w-64"
                />
              </div>
              
              {/* Chat Icon */}
              <Button 
                variant="ghost" 
                size="icon" 
                className="text-gray-400 hover:text-white hover:bg-[#1E1E1E] relative"
                onClick={() => setIsChatModalOpen(true)}
              >
                <MessageSquare className="h-5 w-5" />
                {mockChatHistory.some(chat => chat.unread) && (
                  <span className="absolute top-0 right-0 w-2 h-2 bg-red-500 rounded-full"></span>
                )}
              </Button>
              
              {/* Notifications */}
              <Button variant="ghost" size="icon" className="text-gray-400 hover:text-white hover:bg-[#1E1E1E]">
                <Bell className="h-5 w-5" />
              </Button>
              
              {/* User Menu */}
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    className="flex items-center gap-2 text-gray-300 hover:text-blue-400 hover:bg-[#1E1E1E] rounded-full"
                  >
                    <Avatar className="h-8 w-8 border-2 border-blue-500">
                      <AvatarFallback className="bg-gradient-to-r from-blue-500 to-purple-500 text-white">
                        {userProfile?.user_name?.charAt(0) || "U"}
                      </AvatarFallback>
                    </Avatar>
                    <span className="hidden md:inline-block font-medium">{userProfile?.user_name || "User"}</span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56 bg-[#111111] border-[#1E1E1E] text-gray-100">
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
                  <DropdownMenuItem className="hover:bg-gray-700 cursor-pointer" onClick={handleLogout}>
                    <LogOut className="mr-2 h-4 w-4 text-blue-400" />
                    <span>Log out</span>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
              
              {/* Mobile Menu Button */}
              <Button 
                variant="ghost" 
                size="icon" 
                className="md:hidden text-gray-400 hover:text-white hover:bg-[#1E1E1E]"
                onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
              >
                {isMobileMenuOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
              </Button>
            </div>
          </div>
        </div>
        
        {/* Mobile Navigation */}
        {isMobileMenuOpen && (
          <div className="md:hidden bg-[#111111] border-t border-[#1E1E1E] py-2">
            <div className="container mx-auto px-4">
              <div className="relative mb-2">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                <Input
                  placeholder="Search..."
                  className="pl-10 bg-[#1E1E1E] border-[#2E2E2E] rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent text-white placeholder-gray-400 w-full"
                />
              </div>
              <nav className="flex flex-col space-y-1">
                {navItems.map((item) => {
                  const Icon = item.icon;
                  return (
                    <Link 
                      key={item.path} 
                      href={item.path}
                      className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                        isActive(item.path) 
                          ? 'bg-gray-800 text-blue-400' 
                          : 'text-gray-300 hover:bg-gray-700 hover:text-white'
                      }`}
                      onClick={() => setIsMobileMenuOpen(false)}
                    >
                      <span className="flex items-center">
                        <Icon className="h-4 w-4 mr-2" />
                        {item.name}
                      </span>
                    </Link>
                  );
                })}
              </nav>
            </div>
          </div>
        )}
      </header>

      {/* Main Content */}
      <main className="flex-1">
        {children}
      </main>

      {/* Footer */}
      <footer className="bg-[#111111] border-t border-[#1E1E1E] py-6">
        <div className="container mx-auto px-4">
          <div className="flex flex-col md:flex-row justify-between items-center">
            <div className="mb-4 md:mb-0">
              <div className="flex items-center space-x-2">
                <div className="bg-gradient-to-r from-blue-600 to-purple-600 p-2 rounded-lg">
                  <Wallet className="h-5 w-5 text-white" />
                </div>
                <span className="font-bold text-xl bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
                  BlockChain Service
                </span>
              </div>
              <p className="text-gray-400 text-sm mt-2">
                A platform for blockchain services and transactions
              </p>
            </div>
            <div className="flex flex-col md:flex-row space-y-4 md:space-y-0 md:space-x-8">
              <div>
                <h3 className="text-white font-medium mb-2">Quick Links</h3>
                <ul className="space-y-1">
                  <li>
                    <Link href="/" className="text-gray-400 hover:text-blue-400 text-sm">
                      Home
                    </Link>
                  </li>
                  <li>
                    <Link href="/dashboard" className="text-gray-400 hover:text-blue-400 text-sm">
                      Dashboard
                    </Link>
                  </li>
                  <li>
                    <Link href="/services" className="text-gray-400 hover:text-blue-400 text-sm">
                      Services
                    </Link>
                  </li>
                </ul>
              </div>
              <div>
                <h3 className="text-white font-medium mb-2">Resources</h3>
                <ul className="space-y-1">
                  <li>
                    <Link href="/key" className="text-gray-400 hover:text-blue-400 text-sm">
                      Wallet
                    </Link>
                  </li>
                  <li>
                    <Link href="/transactions" className="text-gray-400 hover:text-blue-400 text-sm">
                      Transactions
                    </Link>
                  </li>
                  <li>
                    <Link href="#" className="text-gray-400 hover:text-blue-400 text-sm">
                      Documentation
                    </Link>
                  </li>
                </ul>
              </div>
              <div>
                <h3 className="text-white font-medium mb-2">Legal</h3>
                <ul className="space-y-1">
                  <li>
                    <Link href="#" className="text-gray-400 hover:text-blue-400 text-sm">
                      Privacy Policy
                    </Link>
                  </li>
                  <li>
                    <Link href="#" className="text-gray-400 hover:text-blue-400 text-sm">
                      Terms of Service
                    </Link>
                  </li>
                  <li>
                    <Link href="#" className="text-gray-400 hover:text-blue-400 text-sm">
                      Cookie Policy
                    </Link>
                  </li>
                </ul>
              </div>
            </div>
          </div>
          <div className="mt-8 pt-4 border-t border-[#1E1E1E] text-center text-gray-400 text-sm">
            <p>© {new Date().getFullYear()} BlockChain Service. All rights reserved.</p>
          </div>
        </div>
      </footer>

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