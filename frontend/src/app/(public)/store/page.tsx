"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { useCart } from "@/context/CartContext";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Modal } from "@/components/ui/Modal";
import {
  ShoppingBag,
  ShoppingCart,
  Check,
  Tag,
  Sparkles,
  ArrowRight,
  Filter,
} from "lucide-react";
import { cn, formatCurrency } from "@/lib/utils";
import { StoreProduct } from "@/types";
import { apiFetch } from "@/lib/api";

type ProductWithSize = StoreProduct & { sizes?: string[]; categoryName: string };

export default function StorePage() {
  const { addItem, itemCount } = useCart();
  const [selectedCategory, setSelectedCategory] = useState<string>("all");
  const [products, setProducts] = useState<ProductWithSize[]>([]);
  const [categories, setCategories] = useState<{ id: string; name: string }[]>([]);
  const [loading, setLoading] = useState(true);

  // Size selection modal state
  const [sizeModalOpen, setSizeModalOpen] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState<ProductWithSize | null>(null);
  const [selectedSize, setSelectedSize] = useState<string>("170cm");
  const [addedToast, setAddedToast] = useState<string | null>(null);

  useEffect(() => {
    async function loadProducts() {
      try {
        const data = await apiFetch<any[]>("/store/products");
        if (data && data.length > 0) {
          const cats: { id: string; name: string }[] = [];
          const allProds: ProductWithSize[] = [];

          data.forEach((group: any) => {
            if (group.category) {
              cats.push({ id: group.category.id, name: group.category.name });
            }
            if (group.products) {
              group.products.forEach((p: any) => {
                let sizes: string[] | undefined;
                if (group.category?.name?.toLowerCase().includes("uniform")) {
                  sizes = ["140cm", "150cm", "160cm", "170cm", "180cm", "190cm"];
                } else if (p.name?.toLowerCase().includes("shoe") || p.name?.toLowerCase().includes("feiyue")) {
                  sizes = ["EU 38", "EU 39", "EU 40", "EU 41", "EU 42", "EU 43", "EU 44"];
                }
                allProds.push({
                  ...p,
                  categoryName: group.category?.name || "General Gear",
                  sizes,
                });
              });
            }
          });

          setCategories(cats);
          setProducts(allProds);
        } else {
          setProducts(fallbackProducts);
        }
      } catch (err) {
        console.error("Failed to load store products:", err);
        setProducts(fallbackProducts);
      } finally {
        setLoading(false);
      }
    }
    loadProducts();
  }, []);

  const filteredProducts = products.filter(
    (p) => selectedCategory === "all" || p.category_id === selectedCategory
  );

  const handleOpenAddToCart = (product: ProductWithSize) => {
    if (product.sizes && product.sizes.length > 0) {
      setSelectedProduct(product);
      setSelectedSize(product.sizes[0]);
      setSizeModalOpen(true);
    } else {
      addItem(product, 1);
      triggerToast(product.name);
    }
  };

  const handleConfirmSizeAdd = () => {
    if (selectedProduct) {
      addItem(selectedProduct, 1, selectedSize);
      triggerToast(`${selectedProduct.name} (${selectedSize})`);
      setSizeModalOpen(false);
    }
  };

  const triggerToast = (itemName: string) => {
    setAddedToast(`Added "${itemName}" to cart!`);
    setTimeout(() => setAddedToast(null), 3000);
  };

  return (
    <div className="py-12 sm:py-16 space-y-12">
      {/* Toast Notification */}
      {addedToast && (
        <div className="fixed bottom-6 right-6 z-50 bg-emerald-950 border border-emerald-700 text-emerald-100 px-4 py-3 rounded-xl shadow-2xl flex items-center gap-2 animate-in fade-in slide-in-from-bottom-5">
          <Check className="h-5 w-5 text-emerald-400 shrink-0" />
          <span className="text-xs font-semibold">{addedToast}</span>
        </div>
      )}

      {/* Header Banner */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center space-y-4">
        <div className="flex items-center justify-center gap-2">
          <Badge variant="gold">Official Martial Arts Academy Store</Badge>
          <Link href="/cart">
            <Button variant="outline" size="sm" className="relative">
              <ShoppingCart className="h-4 w-4" />
              <span>View Cart</span>
              {itemCount > 0 && (
                <span className="absolute -top-2 -right-2 bg-red-600 text-white text-[10px] font-bold rounded-full h-5 w-5 flex items-center justify-center border-2 border-slate-950">
                  {itemCount}
                </span>
              )}
            </Button>
          </Link>
        </div>

        <h1 className="text-4xl sm:text-5xl font-black text-slate-100 uppercase tracking-tight">
          Martial Arts Uniforms, Weapons & Gear
        </h1>
        <p className="text-sm sm:text-base text-slate-400 max-w-2xl mx-auto">
          Traditional Kung Fu uniforms, Karate Gis, Kobudo weapons, Feiyue footwear, and training accessories.
        </p>

        {/* Category Filters */}
        <div className="flex flex-wrap justify-center items-center gap-2 pt-4">
          <button
            onClick={() => setSelectedCategory("all")}
            className={cn(
              "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all cursor-pointer",
              selectedCategory === "all"
                ? "bg-red-600 text-white shadow-md shadow-red-950/60"
                : "bg-slate-900 text-slate-400 hover:text-slate-200 border border-slate-800"
            )}
          >
            All Products
          </button>
          {categories.map((cat) => (
            <button
              key={cat.id}
              onClick={() => setSelectedCategory(cat.id)}
              className={cn(
                "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all cursor-pointer",
                selectedCategory === cat.id
                  ? "bg-red-600 text-white shadow-md shadow-red-950/60"
                  : "bg-slate-900 text-slate-400 hover:text-slate-200 border border-slate-800"
              )}
            >
              {cat.name}
            </button>
          ))}
        </div>
      </div>

      {/* Products Grid */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {loading ? (
          <div className="text-center py-12 text-slate-400 text-sm">Loading official inventory catalog...</div>
        ) : filteredProducts.length === 0 ? (
          <div className="text-center py-12 space-y-3">
            <p className="text-slate-400 text-sm">No products found in this category.</p>
            <Button variant="outline" size="sm" onClick={() => setSelectedCategory("all")}>
              Show All
            </Button>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {filteredProducts.map((p) => {
              const inStock = p.inventory_count > 0;

              return (
                <Card
                  key={p.id}
                  glow="red"
                  className="flex flex-col justify-between overflow-hidden group hover:border-slate-700"
                >
                  <CardHeader className="p-4 space-y-2">
                    <div className="flex items-center justify-between">
                      <Badge variant="primary">{p.categoryName}</Badge>
                      <span className="text-xs font-mono text-slate-400">
                        {inStock ? `${p.inventory_count} in stock` : "Out of Stock"}
                      </span>
                    </div>

                    <CardTitle className="text-base font-bold line-clamp-1 group-hover:text-amber-400 transition-colors">
                      {p.name}
                    </CardTitle>
                    <CardDescription className="text-xs text-slate-400 line-clamp-2">
                      {p.description}
                    </CardDescription>
                  </CardHeader>

                  <CardContent className="p-4 pt-0 space-y-4">
                    <div className="flex items-baseline justify-between pt-2 border-t border-slate-800/80">
                      <div>
                        <span className="text-xl font-black text-slate-100 font-mono">
                          {formatCurrency(p.price_cents)}
                        </span>
                        <span className="text-[10px] text-slate-400 ml-1">+ HST</span>
                      </div>
                    </div>

                    <Button
                      variant={inStock ? "primary" : "outline"}
                      size="sm"
                      className="w-full"
                      disabled={!inStock}
                      onClick={() => handleOpenAddToCart(p)}
                    >
                      <ShoppingBag className="h-4 w-4 mr-1" />
                      <span>{inStock ? (p.sizes ? "Select Size & Add" : "Add to Cart") : "Backorder"}</span>
                    </Button>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        )}
      </div>

      {/* Sizing Modal */}
      <Modal
        isOpen={sizeModalOpen}
        onClose={() => setSizeModalOpen(false)}
        title={`Select Size: ${selectedProduct?.name}`}
        description="Choose your required training size or shoe measurement:"
      >
        <div className="space-y-4 py-2">
          <div className="grid grid-cols-3 gap-2">
            {selectedProduct?.sizes?.map((sz) => (
              <button
                key={sz}
                onClick={() => setSelectedSize(sz)}
                className={cn(
                  "p-3 rounded-xl text-xs font-mono font-bold border transition-all cursor-pointer text-center",
                  selectedSize === sz
                    ? "bg-red-600 text-white border-red-500 shadow-md"
                    : "bg-slate-900 text-slate-300 border-slate-800 hover:border-slate-700"
                )}
              >
                {sz}
              </button>
            ))}
          </div>

          <div className="pt-4 flex justify-end gap-2">
            <Button variant="outline" size="sm" onClick={() => setSizeModalOpen(false)}>
              Cancel
            </Button>
            <Button variant="primary" size="sm" onClick={handleConfirmSizeAdd}>
              <ShoppingCart className="h-4 w-4 mr-1" />
              Add to Cart
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

const fallbackProducts: ProductWithSize[] = [
  {
    id: "prod-uniform-white",
    category_id: "cat-uniforms",
    categoryName: "Uniforms",
    name: "Traditional Kung Fu / Karate Student Uniform",
    description: "Official cotton uniform with durable stitching, pants, and sash for training students.",
    price_cents: 6500,
    image_url: "/products/uniform-white.jpg",
    inventory_count: 35,
    is_active: true,
    sizes: ["140cm", "150cm", "160cm", "170cm", "180cm", "190cm"],
  },
  {
    id: "prod-feiyue-shoes",
    category_id: "cat-shoes",
    categoryName: "Shoes & Gear",
    name: "Original Feiyue Martial Arts Shoes (Black/White)",
    description: "Classic canvas shoes with flexible traction soles engineered for martial arts stances.",
    price_cents: 3500,
    image_url: "/products/feiyue.jpg",
    inventory_count: 50,
    is_active: true,
    sizes: ["EU 38", "EU 39", "EU 40", "EU 41", "EU 42", "EU 43", "EU 44"],
  },
  {
    id: "prod-staff-yinshou",
    category_id: "cat-weapons",
    categoryName: "Weapons",
    name: "Traditional Waxwood Bo Staff (Rokushaku Bo / Gun)",
    description: "Flexible, resilient natural white waxwood staff ideal for traditional staff weapon forms.",
    price_cents: 4500,
    image_url: "/products/staff.jpg",
    inventory_count: 22,
    is_active: true,
    sizes: ["6 Feet (1.8m)", "6.5 Feet (2.0m)"],
  },
];
