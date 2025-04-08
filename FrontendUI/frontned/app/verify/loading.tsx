import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"

export default function VerifyLoading() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-b from-gray-900 to-black px-4 py-12 sm:px-6 lg:px-8">
      <div className="w-full max-w-md">
        <div className="mb-8 text-center">
          <Skeleton className="h-8 w-64 bg-gray-800 mx-auto mb-2" />
          <Skeleton className="h-4 w-48 bg-gray-800 mx-auto" />
        </div>

        <Card className="border-0 bg-gray-800/50 backdrop-blur-sm shadow-2xl">
          <CardHeader className="space-y-1 border-b border-gray-700 pb-6">
            <CardTitle className="text-xl font-medium text-white">
              <Skeleton className="h-6 w-40 bg-gray-700" />
            </CardTitle>
            <CardDescription className="text-gray-400">
              <Skeleton className="h-4 w-full bg-gray-700" />
            </CardDescription>
          </CardHeader>

          <CardContent className="space-y-4 pt-6">
            <div className="flex justify-center space-x-2">
              {Array(6)
                .fill(0)
                .map((_, index) => (
                  <Skeleton key={index} className="w-12 h-14 bg-gray-700" />
                ))}
            </div>

            <Skeleton className="h-4 w-48 bg-gray-700 mx-auto" />

            <div className="pt-4">
              <Skeleton className="h-10 w-full bg-gray-700" />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
