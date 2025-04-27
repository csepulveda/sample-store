"use client"

import Link from "next/link"
import { useCart } from "@/context"

export function Navbar() {
  const { cartItems } = useCart()

  const itemCount = cartItems.reduce((total, item) => total + item.quantity, 0)

  return (
    <nav className="bg-zinc-900 p-4 flex justify-between items-center">
      <Link href="/" className="text-white font-bold text-xl">
        Home
      </Link>
      <Link href="/cart" className="text-white font-semibold">
        Cart ({itemCount})
      </Link>
    </nav>
  )
}