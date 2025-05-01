"use client"

import { createContext, useContext, useReducer, ReactNode } from "react"

interface CartItem {
  productId: string
  productName: string
  quantity: number
}

interface CartState {
  cartItems: CartItem[]
}

type CartAction =
  | { type: "ADD_TO_CART"; payload: CartItem }
  | { type: "REMOVE_FROM_CART"; payload: string }
  | { type: "CLEAR_CART" }

const CartContext = createContext<{
  cartItems: CartItem[]
  dispatch: React.Dispatch<CartAction>
}>({
  cartItems: [],
  dispatch: () => {},
})

function cartReducer(state: CartState, action: CartAction): CartState {
  switch (action.type) {
    case "ADD_TO_CART": {
      const existingItem = state.cartItems.find(item => item.productId === action.payload.productId)
      if (existingItem) {
        return {
          cartItems: state.cartItems.map(item =>
            item.productId === action.payload.productId
              ? { ...item, quantity: item.quantity + action.payload.quantity }
              : item
          ),
        }
      }
      return { cartItems: [...state.cartItems, action.payload] }
    }
    case "REMOVE_FROM_CART":
      return {
        cartItems: state.cartItems.filter(item => item.productId !== action.payload),
      }
    case "CLEAR_CART":
      return { cartItems: [] }
    default:
      return state
  }
}

export function CartProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(cartReducer, { cartItems: [] })

  return (
    <CartContext.Provider value={{ cartItems: state.cartItems, dispatch }}>
      {children}
    </CartContext.Provider>
  )
}

export function useCart() {
  return useContext(CartContext)
}