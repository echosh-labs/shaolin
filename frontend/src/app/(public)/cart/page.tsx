"use client";

import React, { useState } from "react";
import Link from "next/link";
import { useCart } from "@/context/CartContext";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Input, Label } from "@/components/ui/Input";
import { Alert } from "@/components/ui/Alert";
import {
  ShoppingCart,
  Trash2,
  Plus,
  Minus,
  Tag,
  ArrowRight,
  ShieldCheck,
  CreditCard,
  Building,
  Banknote,
  CheckCircle2,
  Printer,
} from "lucide-react";
import { cn, formatCurrency } from "@/lib/utils";

export default function CartPage() {
  const {
    items,
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
  } = useCart();

  const [enteredPromo, setEnteredPromo] = useState("");
  const [promoError, setPromoError] = useState("");
  const [promoSuccess, setPromoSuccess] = useState("");

  const [paymentMethod, setPaymentMethod] = useState<"etransfer" | "credit" | "cash">("etransfer");
  const [isCheckingOut, setIsCheckingOut] = useState(false);
  const [completedOrder, setCompletedOrder] = useState<any>(null);

  // Credit card 2.4% surcharge calculation
  const creditCardSurchargeCents =
    paymentMethod === "credit" ? Math.round(totalCents * 0.024) : 0;
  const grandTotalCents = totalCents + creditCardSurchargeCents;

  const handleApplyPromo = (e: React.FormEvent) => {
    e.preventDefault();
    setPromoError("");
    setPromoSuccess("");

    if (!enteredPromo) return;
    const success = applyPromoCode(enteredPromo);
    if (success) {
      setPromoSuccess(`Promo code '${enteredPromo.toUpperCase()}' applied! 10% discount.`);
    } else {
      setPromoError("Invalid promo code. Try 'SUMMER2026' or 'MA10'.");
    }
  };

  const handleCheckout = () => {
    setIsCheckingOut(true);
    setTimeout(() => {
      setIsCheckingOut(false);
      const orderNumber = `MA-${Math.floor(100000 + Math.random() * 900000)}`;
      setCompletedOrder({
        orderNumber,
        items: [...items],
        subtotalCents,
        taxCents,
        creditCardSurchargeCents,
        grandTotalCents,
        paymentMethod,
        date: new Date().toLocaleDateString("en-CA"),
      });
      clearCart();
    }, 800);
  };

  if (completedOrder) {
    return (
      <div className="py-12 sm:py-16 max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        <Card glow="gold" className="p-8 space-y-6">
          <div className="text-center space-y-2 pb-4 border-b border-slate-800">
            <div className="inline-flex p-3 rounded-full bg-emerald-950 text-emerald-400 border border-emerald-800">
              <CheckCircle2 className="h-8 w-8" />
            </div>
            <h1 className="text-2xl font-black text-slate-100 uppercase">
              Order Confirmed & Received!
            </h1>
            <p className="text-xs text-slate-400">
              Order Reference Number: <strong className="text-amber-400 font-mono">{completedOrder.orderNumber}</strong>
            </p>
          </div>

          {/* Receipt Breakdown */}
          <div className="space-y-4 text-xs">
            <div className="flex justify-between text-slate-400">
              <span>Date:</span>
              <span>{completedOrder.date}</span>
            </div>
            <div className="flex justify-between text-slate-400">
              <span>Payment Method:</span>
              <span className="uppercase font-semibold text-slate-200">
                {completedOrder.paymentMethod === "etransfer"
                  ? "Interac e-Transfer"
                  : completedOrder.paymentMethod === "credit"
                  ? "Credit Card (+2.4% surcharge)"
                  : "Cash on Dojo Pickup"}
              </span>
            </div>

            <div className="divide-y divide-slate-800/80 border-y border-slate-800/80 py-3 space-y-2">
              {completedOrder.items.map((item: any, idx: number) => (
                <div key={idx} className="flex justify-between pt-2">
                  <span>
                    {item.quantity}x {item.product.name} {item.selectedSize ? `(${item.selectedSize})` : ""}
                  </span>
                  <span className="font-mono text-slate-200">
                    {formatCurrency(item.product.price_cents * item.quantity)}
                  </span>
                </div>
              ))}
            </div>

            <div className="space-y-1.5 pt-2">
              <div className="flex justify-between text-slate-400">
                <span>Subtotal:</span>
                <span>{formatCurrency(completedOrder.subtotalCents)}</span>
              </div>
              <div className="flex justify-between text-slate-400">
                <span>13% Ontario HST:</span>
                <span>{formatCurrency(completedOrder.taxCents)}</span>
              </div>
              {completedOrder.creditCardSurchargeCents > 0 && (
                <div className="flex justify-between text-slate-400">
                  <span>Credit Card Processing Fee (2.4%):</span>
                  <span>{formatCurrency(completedOrder.creditCardSurchargeCents)}</span>
                </div>
              )}
              <div className="flex justify-between text-sm font-bold text-slate-100 pt-2 border-t border-slate-800">
                <span>Total Paid / Due:</span>
                <span className="text-amber-400">{formatCurrency(completedOrder.grandTotalCents)}</span>
              </div>
            </div>
          </div>

          <div className="p-4 rounded-xl bg-slate-950 border border-slate-800 text-xs text-slate-400 space-y-1">
            <h5 className="font-semibold text-slate-200">Dojo Pickup Instructions:</h5>
            <p>
              Your equipment order will be prepared at the 393 Dundas St W dojo front desk. Please show your order number when picking up your gear before class.
            </p>
          </div>

          <div className="flex justify-between items-center pt-4">
            <Link href="/store">
              <Button variant="secondary" size="sm">
                &larr; Return to Store
              </Button>
            </Link>
            <Button variant="primary" size="sm" onClick={() => window.print()}>
              <Printer className="h-4 w-4" />
              Print Receipt
            </Button>
          </div>
        </Card>
      </div>
    );
  }

  return (
    <div className="py-12 sm:py-16 space-y-12">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center space-y-4">
        <Badge variant="primary">Shopping Cart & Checkout</Badge>
        <h1 className="text-4xl sm:text-5xl font-black text-slate-100 uppercase tracking-tight">
          Review Your Order
        </h1>
        <p className="text-sm text-slate-400 max-w-xl mx-auto">
          Verify your selected uniform sizes, apply school discount promo codes, and finalize payment.
        </p>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {items.length === 0 ? (
          <Card className="p-12 text-center text-slate-400 max-w-xl mx-auto">
            <ShoppingCart className="h-12 w-12 text-slate-600 mx-auto mb-3" />
            <h3 className="text-lg font-bold text-slate-200">Your shopping cart is empty</h3>
            <p className="text-xs text-slate-400 mt-1 mb-6">
              Browse our official store to add uniforms, Feiyue shoes, weapons, and accessories.
            </p>
            <Link href="/store">
              <Button variant="primary" size="md">
                Browse Dojo Store
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
          </Card>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Items Column */}
            <div className="lg:col-span-2 space-y-4">
              <Card className="p-6 space-y-4">
                <div className="flex justify-between items-center border-b border-slate-800 pb-3">
                  <h3 className="font-bold text-slate-100 text-base">Cart Items ({itemCount})</h3>
                  <button
                    onClick={clearCart}
                    className="text-xs text-slate-400 hover:text-red-400 flex items-center gap-1 cursor-pointer"
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                    Clear All
                  </button>
                </div>

                <div className="divide-y divide-slate-800/80">
                  {items.map((item, idx) => (
                    <div key={idx} className="py-4 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                      <div className="space-y-1">
                        <h4 className="font-bold text-slate-100 text-sm">{item.product.name}</h4>
                        {item.selectedSize && (
                          <Badge variant="gold" size="sm">
                            Size: {item.selectedSize}
                          </Badge>
                        )}
                        <p className="text-xs text-slate-400">
                          {formatCurrency(item.product.price_cents)} each
                        </p>
                      </div>

                      <div className="flex items-center gap-4">
                        {/* Quantity Stepper */}
                        <div className="flex items-center rounded-lg bg-slate-950 border border-slate-800">
                          <button
                            onClick={() => updateQuantity(item.product.id, item.quantity - 1, item.selectedSize)}
                            className="p-1.5 text-slate-400 hover:text-white"
                          >
                            <Minus className="h-3.5 w-3.5" />
                          </button>
                          <span className="px-3 text-xs font-bold text-slate-100">{item.quantity}</span>
                          <button
                            onClick={() => updateQuantity(item.product.id, item.quantity + 1, item.selectedSize)}
                            className="p-1.5 text-slate-400 hover:text-white"
                          >
                            <Plus className="h-3.5 w-3.5" />
                          </button>
                        </div>

                        <div className="text-right min-w-[70px]">
                          <span className="font-bold text-slate-100 text-sm font-mono">
                            {formatCurrency(item.product.price_cents * item.quantity)}
                          </span>
                        </div>

                        <button
                          onClick={() => removeItem(item.product.id, item.selectedSize)}
                          className="text-slate-500 hover:text-red-400 p-1"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              </Card>

              {/* Promo Code Form */}
              <Card className="p-6">
                <form onSubmit={handleApplyPromo} className="space-y-3">
                  <Label htmlFor="promo" className="flex items-center gap-1.5">
                    <Tag className="h-3.5 w-3.5 text-red-500" />
                    Discount Promo Code
                  </Label>
                  <div className="flex gap-2">
                    <Input
                      id="promo"
                      placeholder="e.g. SUMMER2026"
                      value={enteredPromo}
                      onChange={(e) => setEnteredPromo(e.target.value)}
                      className="uppercase"
                    />
                    <Button type="submit" variant="secondary" size="md">
                      Apply
                    </Button>
                  </div>
                  {promoSuccess && <p className="text-xs text-emerald-400">{promoSuccess}</p>}
                  {promoError && <p className="text-xs text-red-400">{promoError}</p>}
                </form>
              </Card>
            </div>

            {/* Checkout Summary Column */}
            <div className="space-y-6">
              <Card glow="red" className="p-6 space-y-6">
                <h3 className="font-bold text-slate-100 text-base border-b border-slate-800 pb-3">
                  Order Summary
                </h3>

                <div className="space-y-2.5 text-xs text-slate-300">
                  <div className="flex justify-between">
                    <span className="text-slate-400">Subtotal:</span>
                    <span className="font-mono text-slate-200">{formatCurrency(subtotalCents)}</span>
                  </div>

                  {discountPercentage > 0 && (
                    <div className="flex justify-between text-emerald-400">
                      <span>Promo Discount ({discountPercentage}%):</span>
                      <span className="font-mono">Applied</span>
                    </div>
                  )}

                  <div className="flex justify-between">
                    <span className="text-slate-400">13% Ontario HST:</span>
                    <span className="font-mono text-slate-200">{formatCurrency(taxCents)}</span>
                  </div>

                  {paymentMethod === "credit" && (
                    <div className="flex justify-between text-amber-400">
                      <span>Credit Card Fee (2.4%):</span>
                      <span className="font-mono">{formatCurrency(creditCardSurchargeCents)}</span>
                    </div>
                  )}

                  <div className="border-t border-slate-800 pt-3 flex justify-between text-base font-black text-slate-100">
                    <span>Total (CAD):</span>
                    <span className="text-red-400 font-mono">{formatCurrency(grandTotalCents)}</span>
                  </div>
                </div>

                {/* Payment Method Selector */}
                <div className="space-y-2 pt-2 border-t border-slate-800">
                  <Label>Select Payment Method</Label>
                  <div className="space-y-2">
                    <button
                      type="button"
                      onClick={() => setPaymentMethod("etransfer")}
                      className={cn(
                        "w-full p-3 rounded-xl border text-left flex items-center justify-between text-xs transition-all cursor-pointer",
                        paymentMethod === "etransfer"
                          ? "bg-slate-900 border-red-500 text-slate-100"
                          : "bg-slate-950 border-slate-800 text-slate-400"
                      )}
                    >
                      <div className="flex items-center gap-2">
                        <Building className="h-4 w-4 text-emerald-400" />
                        <div>
                          <span className="font-semibold block">Interac e-Transfer</span>
                          <span className="text-[10px] text-slate-500">0% fee • etransfer@martialartsacademy.com</span>
                        </div>
                      </div>
                      <Badge variant="success" size="sm">
                        Recommended
                      </Badge>
                    </button>

                    <button
                      type="button"
                      onClick={() => setPaymentMethod("credit")}
                      className={cn(
                        "w-full p-3 rounded-xl border text-left flex items-center justify-between text-xs transition-all cursor-pointer",
                        paymentMethod === "credit"
                          ? "bg-slate-900 border-red-500 text-slate-100"
                          : "bg-slate-950 border-slate-800 text-slate-400"
                      )}
                    >
                      <div className="flex items-center gap-2">
                        <CreditCard className="h-4 w-4 text-sky-400" />
                        <div>
                          <span className="font-semibold block">Credit Card</span>
                          <span className="text-[10px] text-slate-500">Visa / Mastercard (+2.4%)</span>
                        </div>
                      </div>
                    </button>

                    <button
                      type="button"
                      onClick={() => setPaymentMethod("cash")}
                      className={cn(
                        "w-full p-3 rounded-xl border text-left flex items-center justify-between text-xs transition-all cursor-pointer",
                        paymentMethod === "cash"
                          ? "bg-slate-900 border-red-500 text-slate-100"
                          : "bg-slate-950 border-slate-800 text-slate-400"
                      )}
                    >
                      <div className="flex items-center gap-2">
                        <Banknote className="h-4 w-4 text-amber-400" />
                        <div>
                          <span className="font-semibold block">Cash on Pickup</span>
                          <span className="text-[10px] text-slate-500">Pay at Toronto front desk</span>
                        </div>
                      </div>
                    </button>
                  </div>
                </div>

                <Button
                  variant="primary"
                  size="lg"
                  className="w-full"
                  isLoading={isCheckingOut}
                  onClick={handleCheckout}
                >
                  Confirm Order & Generate Receipt
                  <ArrowRight className="h-4 w-4" />
                </Button>
              </Card>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
