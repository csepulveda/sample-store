import Link from "next/link"

export default function Home() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[70vh] text-center gap-8">
      <h1 className="text-4xl font-bold text-white">Welcome to Backoffice</h1>
      <p className="text-zinc-400">Manage your store easily</p>

      <div className="flex gap-6">
        <Link
          href="/products"
          className="bg-blue-600 hover:bg-blue-500 text-white px-6 py-3 rounded-md text-lg"
        >
          Manage Products
        </Link>
        <Link
          href="/orders"
          className="bg-green-600 hover:bg-green-500 text-white px-6 py-3 rounded-md text-lg"
        >
          Manage Orders
        </Link>
      </div>
    </div>
  )
}