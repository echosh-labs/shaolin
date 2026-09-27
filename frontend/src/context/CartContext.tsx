"use client";

import React, { createContext, useContext, useState, useEffect } from "react";
import { StoreProduct } from "@/types";

export interface CartItem {
  product: StoreProduct;
  quantity: number;
  selectedSize?: string;
}

interface CartContextType {
  items: CartItem[];
  addItem: (product: StoreProduct, quantity?: number, selectedSize?: string) => void;
  removeItem: (productId: string, size?: string) => void;
  updateQuantity: (productId: string, quantity: number, size?: string) => void;
  clearCart: () => void;
  promoCode: string;
  applyPromoCode: (code: string) => boolean;
  discountPercentage: number;
  subtotalCents: number;
  taxCents: number;
  totalCents: number;
  itemCount: number;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export function CartProvider({ children }: { children: React.ReactNode }) {
  const [items, setItems] = useState<CartItem[]>([]);
  const [promoCode, setPromoCode] = useState("");
  const [discountPercentage, setDiscountPercentage] = useState(0);

  // Load from local storage
  useEffect(() => {
    try {
      const saved = localStorage.getItem("martialarts_cart");
      if (saved) {
        setItems(JSON.parse(saved));
      }
    } catch {}
  }, []);

  // Save to local storage
  useEffect(() => {
    try {
      localStorage.setItem("martialarts_cart", JSON.stringify(items));
    } catch {}
  }, [items]);

  const addItem = (product: StoreProduct, quantity = 1, selectedSize?: string) => {
    setItems((prev) => {
      const existingIndex = prev.findIndex(
        (item) => item.product.id === product.id && item.selectedSize === selectedSize
      );

      if (existingIndex > -1) {
        const updated = [...prev];
        updated[existingIndex].quantity += quantity;
        return updated;
      }

      return [...prev, { product, quantity, selectedSize }];
    });
  };

  const removeItem = (productId: string, size?: string) => {
    setItems((prev) =>
      prev.filter(
        (item) => !(item.product.id === productId && item.selectedSize === size)
      )
    );
  };

  const updateQuantity = (productId: string, quantity: number, size?: string) => {
    if (quantity <= 0) {
      removeItem(productId, size);
      return;
    }

    setItems((prev) =>
      prev.map((item) => {
        if (item.product.id === productId && item.selectedSize === size) {
          return { ...item, quantity };
        }
        return item;
      })
    );
  };

  const clearCart = () => {
    setItems([]);
    setPromoCode("");
    setDiscountPercentage(0);
  };

  const applyPromoCode = (code: string) => {
    const trimmed = code.trim().toUpperCase();
    if (trimmed === "SUMMER2026" || trimmed === "MA10") {
      setPromoCode(trimmed);
      setDiscountPercentage(10); // 10% discount
      return true;
    }
    return false;
  };

  const rawSubtotal = items.reduce(
    (sum, item) => sum + item.product.price_cents * item.quantity,
    0
  );

  const discountCents = Math.round(rawSubtotal * (discountPercentage / 100));
  const subtotalCents = Math.max(0, rawSubtotal - discountCents);
  const taxCents = Math.round(subtotalCents * 0.13); // 13% Ontario HST
  const totalCents = subtotalCents + taxCents;
  const itemCount = items.reduce((sum, item) => sum + item.quantity, 0);

  return (
    <CartContext.Provider
      value={{
        items,
        addItem,
        removeItem,
        updateQuantity,
        clearCart,
        promoCode,
        applyPromoCode,
        discountPercentage,
        subtotalCents,
        taxCents,
        totalCents,
        itemCount,
      }}
    >
      {children}
    </CartContext.Provider>
  );
}

export function useCart() {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error("useCart must be used within a CartProvider");
  }
  return context;
}
