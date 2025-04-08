"use client"

import type React from "react"
import { useState, useRef, useCallback } from "react"
import { useRouter } from "next/navigation"
import Webcam from "react-webcam"
import { ethers } from "ethers"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { AlertCircle, Camera, User, Mail, Lock, Globe, Shield, ArrowRight } from "lucide-react"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Separator } from "@/components/ui/separator"

const countries = [
    "United States",
    "Canada",
    "United Kingdom",
    "Australia",
    "Germany",
    "France",
    "Japan",
    "China",
    "India",
    "Brazil",
]

export default function RegisterPage() {
    const router = useRouter()
    const webcamRef = useRef<Webcam>(null)

    const [fullName, setFullName] = useState("")
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [country, setCountry] = useState("")
    const [biometricImage, setBiometricImage] = useState<string | null>(null)
    const [biometricHash, setBiometricHash] = useState<string | null>(null)
    const [isLoading, setIsLoading] = useState(false)
    const [error, setError] = useState("")
    const [showWebcam, setShowWebcam] = useState(false)
    const [cameraError, setCameraError] = useState("")

    const captureImage = useCallback(() => {
        if (webcamRef.current) {
            const imageSrc = webcamRef.current.getScreenshot()
            setBiometricImage(imageSrc)

            if (imageSrc) {
                // Convert base64 image to bytes for proper hashing
                const imageBytes = ethers.toUtf8Bytes(imageSrc)
                const hash = ethers.keccak256(imageBytes)
                setBiometricHash(hash)
            }

            setShowWebcam(false)
        }
    }, [webcamRef])

    const handleCameraError = useCallback((error: Error | string) => {
        console.error("Camera error:", error)
        setCameraError("Failed to access camera. Please check permissions and try again.")
        setShowWebcam(false)
    }, [])

    async function handleSubmit(e: React.FormEvent) {
        e.preventDefault()

        if (fullName.length < 3 || fullName.length > 50) {
            setError("Full name must be between 3 and 50 characters")
            return
        }
        if (country && (country.length < 2 || country.length > 50)) {
            setError("Country must be between 2 and 50 characters")
            return
        }

        if (!biometricHash) {
            setError("Biometric data is required")
            return
        }

        if (password.length < 8) {
            setError("Password must be at least 8 characters")
            return
        }

        setIsLoading(true)
        setError("")

        try {
            const payload = {
                fullname: fullName,
                email: email,
                country: country,
                biometric_hash: biometricHash,
                password: password,
                timestamp: Date.now(),
            }

            const response = await fetch("http://localhost:8085/api/auth/register", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload),
            })

            const data = await response.json()

            if (!response.ok) {
                throw new Error(data.message || "Registration failed")
            }

            router.push("/verify")
        } catch (err) {
            console.error("Registration error:", err)
            setError(err instanceof Error ? err.message : "Something went wrong")
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <div className="flex min-h-screen items-center justify-center bg-gradient-to-b from-gray-900 to-black px-4 py-12 sm:px-6 lg:px-8">
            <div className="w-full max-w-md">
                <div className="mb-8 text-center">
                    <h1 className="text-3xl font-bold tracking-tight text-white">Create Account</h1>
                    <p className="mt-2 text-sm text-gray-400">Complete your profile to get started</p>
                </div>

                <Card className="border-0 bg-gray-800/50 backdrop-blur-sm shadow-2xl">
                    <CardHeader className="space-y-1 border-b border-gray-700 pb-6">
                        <CardTitle className="text-xl font-medium text-white">Registration</CardTitle>
                        <CardDescription className="text-gray-400">Enter your information to create an account</CardDescription>
                    </CardHeader>

                    <form onSubmit={handleSubmit}>
                        <CardContent className="space-y-4 pt-6">
                            {error && (
                                <Alert variant="destructive" className="bg-red-900/50 border-red-800 text-red-200">
                                    <AlertCircle className="h-4 w-4" />
                                    <AlertDescription>{error}</AlertDescription>
                                </Alert>
                            )}

                            <div className="space-y-2">
                                <Label htmlFor="fullName" className="text-gray-300 flex items-center gap-2">
                                    <User className="h-4 w-4 text-blue-400" />
                                    Full Name
                                </Label>
                                <Input
                                    id="fullName"
                                    value={fullName}
                                    onChange={(e) => setFullName(e.target.value)}
                                    required
                                    className="bg-gray-700/50 border-gray-600 text-white placeholder:text-gray-500 focus:border-blue-500 focus:ring-blue-500"
                                    placeholder="Enter your full name"
                                />
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="email" className="text-gray-300 flex items-center gap-2">
                                    <Mail className="h-4 w-4 text-blue-400" />
                                    Email
                                </Label>
                                <Input
                                    id="email"
                                    type="email"
                                    placeholder="name@example.com"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    required
                                    className="bg-gray-700/50 border-gray-600 text-white placeholder:text-gray-500 focus:border-blue-500 focus:ring-blue-500"
                                />
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="password" className="text-gray-300 flex items-center gap-2">
                                    <Lock className="h-4 w-4 text-blue-400" />
                                    Password
                                </Label>
                                <Input
                                    id="password"
                                    type="password"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    required
                                    className="bg-gray-700/50 border-gray-600 text-white placeholder:text-gray-500 focus:border-blue-500 focus:ring-blue-500"
                                    placeholder="Create a secure password"
                                />
                                <p className="text-xs text-gray-500">Password must be at least 8 characters</p>
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="country" className="text-gray-300 flex items-center gap-2">
                                    <Globe className="h-4 w-4 text-blue-400" />
                                    Country
                                </Label>
                                <Select value={country} onValueChange={setCountry} required>
                                    <SelectTrigger className="bg-gray-700/50 border-gray-600 text-white focus:border-blue-500 focus:ring-blue-500">
                                        <SelectValue placeholder="Select your country" />
                                    </SelectTrigger>
                                    <SelectContent className="bg-gray-800 border-gray-700 text-white">
                                        {countries.map((c) => (
                                            <SelectItem key={c} value={c} className="focus:bg-gray-700 focus:text-white">
                                                {c}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            </div>

                            <Separator className="my-2 bg-gray-700" />

                            <div className="space-y-2">
                                <Label className="text-gray-300 flex items-center gap-2">
                                    <Shield className="h-4 w-4 text-blue-400" />
                                    Biometric Verification
                                </Label>

                                {cameraError && (
                                    <Alert variant="destructive" className="bg-red-900/50 border-red-800 text-red-200">
                                        <AlertCircle className="h-4 w-4" />
                                        <AlertDescription>{cameraError}</AlertDescription>
                                    </Alert>
                                )}

                                {showWebcam ? (
                                    <div className="space-y-2">
                                        <div className="relative rounded-md overflow-hidden border border-gray-700">
                                            <Webcam
                                                audio={false}
                                                ref={webcamRef}
                                                screenshotFormat="image/jpeg"
                                                className="w-full rounded-md"
                                                minScreenshotWidth={640}
                                                minScreenshotHeight={480}
                                                videoConstraints={{
                                                    facingMode: "user",
                                                    width: { ideal: 640 },
                                                    height: { ideal: 480 },
                                                }}
                                                onUserMediaError={handleCameraError}
                                            />
                                            <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/70 to-transparent p-3">
                                                <p className="text-xs text-white text-center">Position your face in the frame</p>
                                            </div>
                                        </div>
                                        <Button
                                            type="button"
                                            onClick={captureImage}
                                            className="w-full bg-blue-600 hover:bg-blue-500 text-white transition-colors"
                                        >
                                            Capture Image
                                        </Button>
                                    </div>
                                ) : (
                                    <div className="space-y-2">
                                        {biometricImage ? (
                                            <div className="space-y-2">
                                                <div className="relative rounded-md overflow-hidden border border-gray-700">
                                                    <img
                                                        src={biometricImage || "/placeholder.svg"}
                                                        alt="Captured biometric"
                                                        className="w-full rounded-md aspect-video object-cover"
                                                    />
                                                    <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/70 to-transparent p-3">
                                                        <p className="text-xs text-green-300 text-center">Biometric captured successfully</p>
                                                    </div>
                                                </div>
                                                <Button
                                                    type="button"
                                                    variant="outline"
                                                    onClick={() => setShowWebcam(true)}
                                                    className="w-full border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white"
                                                >
                                                    Retake Image
                                                </Button>
                                            </div>
                                        ) : (
                                            <Button
                                                type="button"
                                                variant="outline"
                                                onClick={() => setShowWebcam(true)}
                                                className="w-full border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white flex items-center justify-center gap-2"
                                            >
                                                <Camera className="h-4 w-4" />
                                                Capture Biometric
                                            </Button>
                                        )}
                                    </div>
                                )}

                                {biometricHash && (
                                    <div className="mt-2 text-xs text-gray-500 break-all">
                                        <p>
                                            Verification hash: {biometricHash.substring(0, 20)}...
                                            {biometricHash.substring(biometricHash.length - 10)}
                                        </p>
                                    </div>
                                )}
                            </div>
                        </CardContent>

                        <CardFooter className="flex flex-col gap-4 border-t border-gray-700 pt-6">
                            <Button
                                type="submit"
                                className="w-full bg-blue-600 hover:bg-blue-500 text-white transition-colors"
                                disabled={isLoading || !biometricHash}
                            >
                                {isLoading ? (
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
                    Registering...
                  </span>
                                ) : (
                                    <span className="flex items-center justify-center gap-2">
                    Complete Registration
                    <ArrowRight className="h-4 w-4" />
                  </span>
                                )}
                            </Button>

                            <div className="text-center text-sm text-gray-500">
                                Already have an account?{" "}
                                <a href="/login" className="font-medium text-blue-400 hover:text-blue-300 transition-colors">
                                    Sign in
                                </a>
                            </div>
                        </CardFooter>
                    </form>
                </Card>
            </div>
        </div>
    )
}
