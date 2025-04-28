import { useState } from "react"
import { fetchProductsServerSide, Product, createProduct } from "@/services/products"
import { ProductTable } from "@/components/ProductTable"

export default function ProductsPage({ initialProducts }: { initialProducts: Product[] }) {
  const [products, setProducts] = useState<Product[]>(initialProducts)
  const [newProduct, setNewProduct] = useState<Partial<Product>>({
    name: "",
    description: "",
    price: 0,
    stock: 0,
  })
  const [showForm, setShowForm] = useState(false)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  const reload = async () => {
    const updatedProducts = await fetchProductsServerSide()
    setProducts(updatedProducts)
  }

  const handleAddProduct = async () => {
    setSaving(true)
    setErrorMessage(null)

    try {
      await createProduct(newProduct)
      setShowForm(false)
      setNewProduct({ name: "", description: "", price: 0, stock: 0 })
      reload()
    } catch (err: any) {
      const message = err.message || "Failed to create product"
      setErrorMessage(message)
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="container mx-auto p-8">
      <h1 className="text-3xl font-bold mb-8">Manage Products</h1>

      <button
        onClick={() => {
          setShowForm(!showForm)
          setErrorMessage(null)
        }}
        className="mb-6 bg-green-600 hover:bg-green-500 text-white px-4 py-2 rounded"
      >
        {showForm ? "Cancel" : "Add Product"}
      </button>

      {showForm && (
        <div className="bg-zinc-800 p-4 rounded-lg mb-6">
          <div className="flex flex-col gap-4">
            {errorMessage && (
              <div className="text-red-500 text-sm">{errorMessage}</div>
            )}
            <input
              type="text"
              placeholder="Name"
              value={newProduct.name}
              onChange={(e) => setNewProduct({ ...newProduct, name: e.target.value })}
              className="bg-zinc-700 border border-zinc-600 rounded px-3 py-2"
            />
            <input
              type="text"
              placeholder="Description"
              value={newProduct.description}
              onChange={(e) => setNewProduct({ ...newProduct, description: e.target.value })}
              className="bg-zinc-700 border border-zinc-600 rounded px-3 py-2"
            />
            <input
              type="number"
              placeholder="Price"
              value={newProduct.price}
              onChange={(e) => setNewProduct({ ...newProduct, price: parseFloat(e.target.value) })}
              className="bg-zinc-700 border border-zinc-600 rounded px-3 py-2"
            />
            <input
              type="number"
              placeholder="Stock"
              value={newProduct.stock}
              onChange={(e) => setNewProduct({ ...newProduct, stock: parseInt(e.target.value) })}
              className="bg-zinc-700 border border-zinc-600 rounded px-3 py-2"
            />
            <button
              onClick={handleAddProduct}
              disabled={saving}
              className={`${
                saving ? "bg-gray-600" : "bg-blue-600 hover:bg-blue-500"
              } text-white px-4 py-2 rounded`}
            >
              {saving ? "Saving..." : "Save Product"}
            </button>
          </div>
        </div>
      )}

      <ProductTable products={products} reload={reload} />
    </div>
  )
}

export async function getServerSideProps() {
  const initialProducts = await fetchProductsServerSide()
  return { props: { initialProducts } }
}