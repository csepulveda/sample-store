export interface Product {
    id: string
    name: string
    description: string
    price: number
    stock: number
}

function getBaseUrl() {
if (typeof window !== "undefined") {
    return "" // Client-side uses relative path
}
return process.env.NEXT_PUBLIC_APP_URL || "http://ui-backoffice:3000"
}

export async function fetchProductsServerSide(): Promise<Product[]> {
const baseUrl = getBaseUrl()
const res = await fetch(`${baseUrl}/api/proxy/products`, {
    method: "GET",
    headers: { "Content-Type": "application/json" },
})
if (!res.ok) {
    throw new Error("Failed to fetch products")
}
return res.json()
}

export async function createProduct(product: Partial<Product>) {
const baseUrl = getBaseUrl()
const res = await fetch(`${baseUrl}/api/proxy/products`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(product),
})
if (!res.ok) {
    throw new Error("Failed to create product")
}
}

export async function updateProduct(id: string, product: Partial<Product>) {
const baseUrl = getBaseUrl()
const res = await fetch(`${baseUrl}/api/proxy/products/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(product),
})
if (!res.ok) {
    throw new Error("Failed to update product")
}
}

export async function deleteProduct(id: string) {
const baseUrl = getBaseUrl()
const res = await fetch(`${baseUrl}/api/proxy/products/${id}`, {
    method: "DELETE",
})
if (!res.ok) {
    throw new Error("Failed to delete product")
}
}