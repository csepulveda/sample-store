"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"

export function Navbar() {
  const pathname = usePathname()

  const linkClass = (href: string) =>
    `px-4 py-2 rounded-md ${
      pathname === href
        ? "bg-blue-600 text-white"
        : "text-zinc-300 hover:bg-zinc-700 hover:text-white"
    }`

  return (
    <nav className="bg-zinc-900 border-b border-zinc-800 p-4 flex gap-6 items-center">
      <Link href="/" className="text-white font-bold text-xl">
        Backoffice
      </Link>
      <div className="flex gap-4">
        <Link href="/products" className={linkClass("/products")}>
          Products
        </Link>
        <Link href="/orders" className={linkClass("/orders")}>
          Orders
        </Link>
      </div>
    </nav>
  )
}