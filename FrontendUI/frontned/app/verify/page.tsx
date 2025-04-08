"use client"

import type React from "react"

import { useState, useRef, useEffect } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { AlertCircle, CheckCircle, ArrowRight, RefreshCw } from "lucide-react"
import { Alert, AlertDescription } from "@/components/ui/alert"

export default function VerifyPage() {
  const router = useRouter()
  const [verificationCode, setVerificationCode] = useState<string[]>(Array(6).fill(""))
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState(false)
  const inputRefs = useRef<(HTMLInputElement | null)[]>([])

  // Initialize refs array
  useEffect(() => {
    inputRefs.current = inputRefs.current.slice(0, 6)
  }, [])

  // Handle input change
  const handleChange = (index: number, value: string) => {
    // Only allow numbers
    if (!/^\d*$/.test(value)) return

    const newVerificationCode = [...verificationCode]
    // Take only the last character if multiple are pasted
    newVerificationCode[index] = value.slice(-1)
    setVerificationCode(newVerificationCode)

    // Auto-advance to next field if a digit was entered
    if (value && index < 5) {
      inputRefs.current[index + 1]?.focus()
    }
  }

  // Handle key down events
  const handleKeyDown = (index: number, e: React.KeyboardEvent<HTMLInputElement>) => {
    // Move to previous input on backspace if current input is empty
    if (e.key === "Backspace" && !verificationCode[index] && index > 0) {
      inputRefs.current[index - 1]?.focus()
    }

    // Handle arrow keys
    if (e.key === "ArrowLeft" && index > 0) {
      inputRefs.current[index - 1]?.focus()
    }
    if (e.key === "ArrowRight" && index < 5) {
      inputRefs.current[index + 1]?.focus()
    }
  }

  // Handle paste event
  const handlePaste = (e: React.ClipboardEvent) => {
    e.preventDefault()
    const pastedData = e.clipboardData.getData("text/plain").trim()

    // Check if pasted content is a 6-digit number
    if (/^\d{6}$/.test(pastedData)) {
      const digits = pastedData.split("")
      setVerificationCode(digits)

      // Focus the last input
      inputRefs.current[5]?.focus()
    }
  }

  // Submit verification code
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    const code = verificationCode.join("")

    // Validate code length
    if (code.length !== 6) {
      setError("Please enter all 6 digits of the verification code")
      return
    }

    setIsLoading(true)
    setError("")

    try {
        const payload = {
            email: "swanhtetaungp@gmail.com",
            code: code
        }
      const response = await fetch("http://localhost:8085/api/auth/verify", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(payload),
      })

      const data = await response.json()

    //   if (!response.ok) {
    //     throw new Error(data.message || "Verification failed")
    //   }

      // Show success message before redirecting
      setSuccess(true)

      setTimeout(() => {
        router.push("/login")
      }, 2000)
    } catch (err) {
      console.error("Verification error:", err)
      setError(err instanceof Error ? err.message : "Verification failed")
    } finally {
      setIsLoading(false)
    }
  }

  // Resend verification code
  const handleResend = async () => {
    setIsLoading(true)
    setError("")

    try {
      // This is a placeholder - you would need to implement the actual resend endpoint
      const response = await fetch("http://localhost:8085/api/auth/resend-code", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.message || "Failed to resend code")
      }

      // Clear the current code
      setVerificationCode(Array(6).fill(""))
      // Focus the first input
      inputRefs.current[0]?.focus()
    } catch (err) {
      console.error("Resend error:", err)
      setError(err instanceof Error ? err.message : "Failed to resend verification code")
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-b from-gray-900 to-black px-4 py-12 sm:px-6 lg:px-8">
      <div className="w-full max-w-md">
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold tracking-tight text-white">Verify Your Account</h1>
          <p className="mt-2 text-sm text-gray-400">Enter the 6-digit code sent to your email</p>
        </div>

        <Card className="border-0 bg-gray-800/50 backdrop-blur-sm shadow-2xl">
          <CardHeader className="space-y-1 border-b border-gray-700 pb-6">
            <CardTitle className="text-xl font-medium text-white">Verification Code</CardTitle>
            <CardDescription className="text-gray-400">We've sent a 6-digit code to your email address</CardDescription>
          </CardHeader>

          <form onSubmit={handleSubmit}>
            <CardContent className="space-y-4 pt-6">
              {error && (
                <Alert variant="destructive" className="bg-red-900/50 border-red-800 text-red-200">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>{error}</AlertDescription>
                </Alert>
              )}

              {success && (
                <Alert className="bg-green-900/50 border-green-800 text-green-200">
                  <CheckCircle className="h-4 w-4" />
                  <AlertDescription>Verification successful! Redirecting...</AlertDescription>
                </Alert>
              )}

              <div className="flex justify-center space-x-2">
                {verificationCode.map((digit, index) => (
                  <Input
                    key={index}
                    ref={(el) => {
                      inputRefs.current[index] = el;
                    }}
                    type="text"
                    inputMode="numeric"
                    maxLength={1}
                    value={digit}
                    onChange={(e) => handleChange(index, e.target.value)}
                    onKeyDown={(e) => handleKeyDown(index, e)}
                    onPaste={index === 0 ? handlePaste : undefined}
                    className="w-12 h-14 text-center text-xl font-bold bg-gray-700/50 border-gray-600 text-white focus:border-blue-500 focus:ring-blue-500"
                    autoFocus={index === 0}
                  />
                ))}
              </div>

              <p className="text-center text-xs text-gray-500">
                Didn't receive a code? Check your spam folder or{" "}
                <button
                  type="button"
                  onClick={handleResend}
                  className="text-blue-400 hover:text-blue-300 transition-colors"
                  disabled={isLoading}
                >
                  resend code
                </button>
              </p>
            </CardContent>

            <CardFooter className="flex flex-col gap-4 border-t border-gray-700 pt-6">
              <Button
                type="submit"
                className="w-full bg-blue-600 hover:bg-blue-500 text-white transition-colors"
                disabled={isLoading || success}
              >
                {isLoading ? (
                  <span className="flex items-center justify-center gap-2">
                    <RefreshCw className="h-4 w-4 animate-spin" />
                    Verifying...
                  </span>
                ) : (
                  <span className="flex items-center justify-center gap-2">
                    Verify Account
                    <ArrowRight className="h-4 w-4" />
                  </span>
                )}
              </Button>

              <div className="text-center text-sm text-gray-500">
                <a href="/login" className="font-medium text-blue-400 hover:text-blue-300 transition-colors">
                  Back to login
                </a>
              </div>
            </CardFooter>
          </form>
        </Card>

        <div className="mt-6 flex items-center justify-center">
          <span className="inline-flex items-center text-xs text-gray-500">
            <svg
              className="h-4 w-4 mr-1"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
              ></path>
            </svg>
            Secure verification with end-to-end encryption
          </span>
        </div>
      </div>
    </div>
  )
}
