import "@/styles/globals.css"
import type { AppProps } from "next/app"
import { Toaster } from "react-hot-toast"
import { Navbar } from "@/components/Navbar"

export default function App({ Component, pageProps }: AppProps) {
  return (
    <>
      <Navbar />
      <main className="p-6">
        <Component {...pageProps} />
      </main>
      <Toaster position="top-right" />
    </>
  )
}