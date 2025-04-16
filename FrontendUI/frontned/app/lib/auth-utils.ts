import { jwtDecode } from "jwt-decode";

const API_BASE_URL = "http://localhost:8085/api/auth";

interface LoginRequest {
  email: string;
  password: string;
}

interface RegisterRequest {
  fullname: string;
  email: string;
  country: string;
  biometric_hash: string;
  password: string;
  timestamp: number;
}

interface VerifyRequest {
  email: string;
  code: string;
}

interface ForgotPasswordRequest {
  email: string;
}

interface ResetPasswordRequest {
  token: string;
  newPassword: string;
}

interface WalletRequest {
  userid: string;
  phrase: string;
}

interface LoginResponse {
  access_token: string;
  user_account_wallet: boolean;
  email: string;
}

interface RegisterResponse {
  email: string;
  full_name: string;
  verification_ttl_minutes: number;
  registered_at: string;
}

interface UserProfileResponse {
  email: string;
  user_name: string;
  verified: boolean;
  created_at: string;
  wallet_created: boolean;
  wallet_created_time?: string;
  solana_address?: string;
  ethereum_address?: string;
}

interface VerificationResponse {
  email: string;
  verified_at: string;
  account_status: string;
}

interface ApiResponse<T> {
  success: boolean;
  message: string;
  data?: T;
  error?: string;
}

// Helper function to handle API responses
async function handleResponse<T>(response: Response): Promise<ApiResponse<T>> {
  const data = await response.json();
  console.log("API Response:", data);
  
  if (!response.ok) {
    return {
      success: false,
      message: data.message || "An error occurred",
      error: data.error || "Unknown error"
    };
  }
  
  return {
    success: true,
    message: data.message || "Operation successful",
    data: data.data
  };
}

// Login user
export async function login(credentials: LoginRequest): Promise<ApiResponse<LoginResponse>> {
  try {
    const response = await fetch(`${API_BASE_URL}/login`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(credentials),
    });
    
    return handleResponse<LoginResponse>(response);
  } catch (error) {
    return {
      success: false,
      message: "Failed to connect to the server",
      error: error instanceof Error ? error.message : "Unknown error"
    };
  }
}

// Register new user
export async function register(userData: RegisterRequest): Promise<ApiResponse<RegisterResponse>> {
  try {
    const response = await fetch(`${API_BASE_URL}/register`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(userData),
    });
    
    return handleResponse<RegisterResponse>(response);
  } catch (error) {
    return {
      success: false,
      message: "Failed to connect to the server",
      error: error instanceof Error ? error.message : "Unknown error"
    };
  }
}

// Verify user email
export async function verifyEmail(verifyData: VerifyRequest): Promise<ApiResponse<VerificationResponse>> {
  try {
    const response = await fetch(`${API_BASE_URL}/verify`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(verifyData),
    });
    
    return handleResponse<VerificationResponse>(response);
  } catch (error) {
    return {
      success: false,
      message: "Failed to connect to the server",
      error: error instanceof Error ? error.message : "Unknown error"
    };
  }
}

// Get user profile
export async function getUserProfile(email: string, token: string): Promise<ApiResponse<UserProfileResponse>> {
  try {
    const response = await fetch(`${API_BASE_URL}/me?email=${encodeURIComponent(email)}`, {
      method: "GET",
      headers: {
        "Authorization": `Bearer ${token}`,
      },
    });
    
    const data = await response.json();
    console.log("User profile raw response:", data);
    
    if (!response.ok) {
      return {
        success: false,
        message: data.message || "Failed to fetch user profile",
        error: data.error || "Unknown error"
      };
    }
    
    return {
      success: true,
      message: data.message || "Profile retrieved successfully",
      data: data.data
    };
  } catch (error) {
    return {
      success: false,
      message: "Failed to connect to the server",
      error: error instanceof Error ? error.message : "Unknown error"
    };
  }
}

// Refresh token
export async function refreshToken(): Promise<ApiResponse<{access_token: string}>> {
  try {
    const response = await fetch(`${API_BASE_URL}/refresh`, {
      method: "GET",
      credentials: "include",
    });
    
    return handleResponse<{access_token: string}>(response);
  } catch (error) {
    return {
      success: false,
      message: "Failed to connect to the server",
      error: error instanceof Error ? error.message : "Unknown error"
    };
  }
}

// Forgot password
export async function forgotPassword(email: string): Promise<ApiResponse<{email: string, expires_at: string}>> {
  try {
    const response = await fetch(`${API_BASE_URL}/forgot-password`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email }),
    });
    
    return handleResponse<{email: string, expires_at: string}>(response);
  } catch (error) {
    return {
      success: false,
      message: "Failed to connect to the server",
      error: error instanceof Error ? error.message : "Unknown error"
    };
  }
}

// Reset password
export async function resetPassword(token: string, newPassword: string): Promise<{ success: boolean; message?: string }> {
  try {
    const response = await fetch(`${API_BASE_URL}/reset-password`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        token,
        newPassword,
      }),
    });

    const data = await response.json();

    if (!response.ok) {
      return {
        success: false,
        message: data.message || 'Failed to reset password',
      };
    }

    return {
      success: true,
      message: 'Password reset successfully',
    };
  } catch (error) {
    return {
      success: false,
      message: error instanceof Error ? error.message : 'An error occurred while resetting password',
    };
  }
}

// Store wallet in vault
export async function storeWallet(walletData: WalletRequest, token: string): Promise<ApiResponse<{user_id: string, wallet_status: string, created_at: string}>> {
  try {
    const response = await fetch(`${API_BASE_URL}/wallet`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`,
      },
      body: JSON.stringify(walletData),
    });
    
    return handleResponse<{user_id: string, wallet_status: string, created_at: string}>(response);
  } catch (error) {
    return {
      success: false,
      message: "Failed to connect to the server",
      error: error instanceof Error ? error.message : "Unknown error"
    };
  }
}

// Logout
export async function logout(token: string): Promise<ApiResponse<{logged_out_at: string, session_ended: boolean}>> {
  try {
    const response = await fetch(`${API_BASE_URL}/logout`, {
      method: "POST",
      headers: {
        "Authorization": `Bearer ${token}`,
      },
      credentials: "include",
    });
    
    return handleResponse<{logged_out_at: string, session_ended: boolean}>(response);
  } catch (error) {
    return {
      success: false,
      message: "Failed to connect to the server",
      error: error instanceof Error ? error.message : "Unknown error"
    };
  }
}

// Check if token is valid
export function isTokenValid(token: string): boolean {
  try {
    const decoded = jwtDecode(token);
    const currentTime = Date.now() / 1000;
    
    return decoded.exp ? decoded.exp > currentTime : false;
  } catch (error) {
    return false;
  }
}

// Get token from localStorage
export function getToken(): string | null {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('access_token');
  }
  return null;
}

// Set token in localStorage
export function setToken(token: string): void {
  if (typeof window !== 'undefined') {
    localStorage.setItem('access_token', token);
  }
}

// Remove token from localStorage
export function removeToken(): void {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('access_token');
  }
}

// Get user email from token
export function getUserEmailFromToken(token: string): string | null {
  try {
    const decoded = jwtDecode<{email: string}>(token);
    return decoded.email || null;
  } catch (error) {
    return null;
  }
}

export async function createWallet(token: string): Promise<{ success: boolean; message: string }> {
  try {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/create-wallet`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      }
    });

    const data = await response.json();
    return {
      success: response.ok,
      message: data.message || (response.ok ? 'Wallet created successfully' : 'Failed to create wallet')
    };
  } catch (error) {
    console.error('Error creating wallet:', error);
    return {
      success: false,
      message: 'An error occurred while creating the wallet'
    };
  }
} 