"use client"
import { createContext, useContext, useReducer, ReactNode } from "react"
import { Product } from "@/services/products"

type CartItem = Product & { quantity: number }

type State = {
  items: CartItem[]
}

type Action =
  | { type: "ADD_ITEM"; product: Product }
  | { type: "REMOVE_ITEM"; id: string }
  | { type: "CLEAR_CART" }

const CartContext = createContext<{
  state: State
  dispatch: React.Dispatch<Action>
}>({ state: { items: [] }, dispatch: () => {} })

function cartReducer(state: State, action: Action): State {
  switch (action.type) {
    case "ADD_ITEM": {
      const existing = state.items.find(i => i.id === action.product.id)
      if (existing) {
        return {
          items: state.items.map(i =>
            i.id === action.product.id ? { ...i, quantity: i.quantity + 1 } : i
          )
        }
      }
      return { items: [...state.items, { ...action.product, quantity: 1 }] }
    }
    case "REMOVE_ITEM":
      return { items: state.items.filter(i => i.id !== action.id) }
    case "CLEAR_CART":
      return { items: [] }
    default:
      return state
  }
}

export function CartProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(cartReducer, { items: [] })
  return (
    <CartContext.Provider value={{ state, dispatch }}>
      {children}
    </CartContext.Provider>
  )
}

export function useCart() {
  return useContext(CartContext)
}