import "@/styles/globals.css"
import type { AppProps } from "next/app"
import { CartProvider } from "@/context"
import { Navbar } from "@/components/Navbar"

export default function App({ Component, pageProps }: AppProps) {
  return (
    <CartProvider>
      <div className="min-h-screen bg-black">
        <Navbar />
        <Component {...pageProps} />
      </div>
    </CartProvider>
  )
}