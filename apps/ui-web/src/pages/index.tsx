import { Product } from "@/services/products"
import { ProductCard } from "@/components/ProductCard"
import { fetchProductsServerSide } from "@/services/products"
import { ProductGrid } from "@/components/ProductGrid"

export default function Home({ products }: { products: Product[] }) {
  return (
    <div className="container mx-auto p-8">
      <h1 className="text-3xl font-bold mb-8 text-white">Products</h1>
      <ProductGrid products={products} />
    </div>
  )
}

export async function getServerSideProps() {
  const products = await fetchProductsServerSide()
  return { props: { products } }
}