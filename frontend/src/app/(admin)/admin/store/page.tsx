"use client";

import React, { useState, useEffect } from "react";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Alert } from "@/components/ui/Alert";
import {
  PackageCheck,
  Package,
  CheckCircle2,
  Clock,
  Building,
  CreditCard,
  Banknote,
  Plus,
  Minus,
  Search,
} from "lucide-react";
import { cn, formatCurrency, formatDate } from "@/lib/utils";
import { apiFetch } from "@/lib/api";

interface StoreOrder {
  id: string;
  orderNumber: string;
  customerName: string;
  email: string;
  date: string;
  items: string;
  totalCents: number;
  paymentMethod: "etransfer" | "credit" | "cash";
  status: "pending" | "ready_for_pickup" | "completed";
}

interface InventoryItem {
  id: string;
  name: string;
  category: string;
  priceCents: number;
  stock: number;
  isActive: boolean;
}

export default function AdminStorePage() {
  const [selectedTab, setSelectedTab] = useState<"orders" | "inventory">("orders");
  const [alertMsg, setAlertMsg] = useState<{ type: "success" | "error"; text: string } | null>(null);

  const [orders, setOrders] = useState<StoreOrder[]>([
    {
      id: "ord-01",
      orderNumber: "MA-928104",
      customerName: "Justin Kowalski",
      email: "student@martialartsacademy.com",
      date: "2026-08-18",
      items: "1x White Kung Fu Uniform (170cm), 1x Feiyue Shoes (EU 41)",
      totalCents: 11300,
      paymentMethod: "etransfer",
      status: "ready_for_pickup",
    },
    {
      id: "ord-02",
      orderNumber: "MA-582910",
      customerName: "Elena Rostova",
      email: "elena@martialartsacademy.com",
      date: "2026-08-17",
      items: "1x Waxwood Bo Staff (6 Feet)",
      totalCents: 5085,
      paymentMethod: "credit",
      status: "completed",
    },
  ]);

  const [inventory, setInventory] = useState<InventoryItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadStoreAdminData() {
      try {
        const [prodsRes, ordersRes] = await Promise.all([
          apiFetch<any[]>("/store/products").catch(() => []),
          apiFetch<any[]>("/store/orders/history").catch(() => []),
        ]);

        if (prodsRes && prodsRes.length > 0) {
          const flat: InventoryItem[] = [];
          prodsRes.forEach((group: any) => {
            if (group.products) {
              group.products.forEach((p: any) => {
                flat.push({
                  id: p.id,
                  name: p.name,
                  category: group.category?.name || "Equipment",
                  priceCents: p.price_cents,
                  stock: p.inventory_count,
                  isActive: p.is_active,
                });
              });
            }
          });
          setInventory(flat);
        } else {
          setInventory(fallbackInventory);
        }

        if (ordersRes && ordersRes.length > 0) {
          const ords: StoreOrder[] = ordersRes.map((o: any) => ({
            id: o.id,
            orderNumber: `MA-${o.id.substring(0, 6).toUpperCase()}`,
            customerName: "Student Member",
            email: "student@martialartsacademy.com",
            date: o.created_at ? o.created_at.split("T")[0] : "2026-08-18",
            items: "Uniform / Training Supplies",
            totalCents: o.total_cents,
            paymentMethod: "etransfer",
            status: "ready_for_pickup",
          }));
          setOrders(ords);
        }
      } catch (err) {
        console.error("Failed to load store admin data:", err);
        setInventory(fallbackInventory);
      } finally {
        setLoading(false);
      }
    }
    loadStoreAdminData();
  }, []);

  const handleUpdateOrderStatus = (orderId: string, newStatus: "pending" | "ready_for_pickup" | "completed") => {
    setOrders((prev) =>
      prev.map((o) => (o.id === orderId ? { ...o, status: newStatus } : o))
    );
    setAlertMsg({
      type: "success",
      text: `Order status updated to "${newStatus.replace("_", " ").toUpperCase()}".`,
    });
  };

  const handleAdjustStock = (itemId: string, delta: number) => {
    setInventory((prev) =>
      prev.map((item) => {
        if (item.id === itemId) {
          const newStock = Math.max(0, item.stock + delta);
          return { ...item, stock: newStock };
        }
        return item;
      })
    );
  };

  return (
    <div className="space-y-8">
      {/* Top Banner */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-6 rounded-2xl bg-slate-900 border border-slate-800 shadow-xl">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <Badge variant="primary">Operations Console</Badge>
            <Badge variant="gold">Store & Inventory Control</Badge>
          </div>
          <h1 className="text-2xl font-black text-slate-100 uppercase tracking-tight">
            Store Fulfillment & Stock Control
          </h1>
          <p className="text-xs text-slate-400">
            Process student store pickup orders, verify e-transfer receipts, and manage uniform inventory levels.
          </p>
        </div>
      </div>

      {alertMsg && (
        <Alert
          variant={alertMsg.type === "success" ? "success" : "error"}
          title={alertMsg.type === "success" ? "Inventory Updated" : "Alert"}
        >
          {alertMsg.text}
        </Alert>
      )}

      {/* Tabs */}
      <div className="flex items-center gap-2 border-b border-slate-800 pb-3">
        <button
          onClick={() => setSelectedTab("orders")}
          className={cn(
            "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all cursor-pointer flex items-center gap-2",
            selectedTab === "orders"
              ? "bg-red-600 text-white shadow-md shadow-red-950/50"
              : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
          )}
        >
          <PackageCheck className="h-4 w-4" />
          Customer Orders ({orders.length})
        </button>

        <button
          onClick={() => setSelectedTab("inventory")}
          className={cn(
            "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all cursor-pointer flex items-center gap-2",
            selectedTab === "inventory"
              ? "bg-amber-600 text-white shadow-md shadow-amber-950/50"
              : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
          )}
        >
          <Package className="h-4 w-4" />
          Inventory Stock ({inventory.length})
        </button>
      </div>

      {/* TAB 1: Orders */}
      {selectedTab === "orders" && (
        <Card className="overflow-hidden border-slate-800">
          <div className="divide-y divide-slate-800">
            {orders.map((ord) => (
              <div
                key={ord.id}
                className="p-5 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 hover:bg-slate-900/40 transition-colors"
              >
                <div className="space-y-1.5">
                  <div className="flex items-center gap-2">
                    <span className="font-mono text-sm font-bold text-amber-400">{ord.orderNumber}</span>
                    <Badge
                      variant={
                        ord.status === "completed"
                          ? "success"
                          : ord.status === "ready_for_pickup"
                          ? "primary"
                          : "outline"
                      }
                      size="sm"
                    >
                      {ord.status.replace("_", " ").toUpperCase()}
                    </Badge>
                  </div>

                  <div className="text-xs text-slate-200 font-bold">{ord.customerName} ({ord.email})</div>
                  <p className="text-xs text-slate-400">{ord.items}</p>
                  <div className="text-[10px] text-slate-500 font-mono">Date: {ord.date} • Total: {formatCurrency(ord.totalCents)}</div>
                </div>

                <div className="flex items-center gap-2">
                  {ord.status === "pending" && (
                    <Button
                      variant="primary"
                      size="sm"
                      onClick={() => handleUpdateOrderStatus(ord.id, "ready_for_pickup")}
                    >
                      Mark Ready for Pickup
                    </Button>
                  )}
                  {ord.status === "ready_for_pickup" && (
                    <Button
                      variant="accent"
                      size="sm"
                      onClick={() => handleUpdateOrderStatus(ord.id, "completed")}
                    >
                      <CheckCircle2 className="h-3.5 w-3.5 mr-1" />
                      Complete & Hand Over
                    </Button>
                  )}
                  {ord.status === "completed" && (
                    <Badge variant="success">Fulfilled</Badge>
                  )}
                </div>
              </div>
            ))}
          </div>
        </Card>
      )}

      {/* TAB 2: Inventory */}
      {selectedTab === "inventory" && (
        <Card className="overflow-hidden border-slate-800">
          <div className="divide-y divide-slate-800">
            {inventory.slice(0, 10).map((item) => (
              <div
                key={item.id}
                className="p-4 flex items-center justify-between hover:bg-slate-900/40 transition-colors"
              >
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <span className="text-sm font-bold text-slate-100">{item.name}</span>
                    <Badge variant="outline" size="sm">{item.category}</Badge>
                  </div>
                  <span className="text-xs text-slate-400 font-mono">
                    Unit Price: {formatCurrency(item.priceCents)}
                  </span>
                </div>

                <div className="flex items-center gap-3">
                  <div className="flex items-center gap-2 bg-slate-950 p-1 rounded-xl border border-slate-800">
                    <button
                      onClick={() => handleAdjustStock(item.id, -1)}
                      className="p-1 rounded-lg text-slate-400 hover:text-white hover:bg-slate-850 cursor-pointer"
                    >
                      <Minus className="h-3.5 w-3.5" />
                    </button>
                    <span className="text-xs font-mono font-bold w-10 text-center text-slate-100">
                      {item.stock}
                    </span>
                    <button
                      onClick={() => handleAdjustStock(item.id, 1)}
                      className="p-1 rounded-lg text-slate-400 hover:text-white hover:bg-slate-850 cursor-pointer"
                    >
                      <Plus className="h-3.5 w-3.5" />
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </Card>
      )}
    </div>
  );
}

const fallbackInventory: InventoryItem[] = [
  {
    id: "prod-1",
    name: "Standard White Kung Fu Uniform",
    category: "Uniforms",
    priceCents: 6500,
    stock: 35,
    isActive: true,
  },
  {
    id: "prod-feiyue",
    name: "Original Feiyue Martial Arts Shoes",
    category: "Footwear",
    priceCents: 3500,
    stock: 50,
    isActive: true,
  },
];
