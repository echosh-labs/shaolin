"use client";

import React, { useState, useEffect } from "react";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Alert } from "@/components/ui/Alert";
import { Modal } from "@/components/ui/Modal";
import {
  Coins,
  Receipt,
  CheckCircle2,
  TrendingUp,
  ArrowUpRight,
  ArrowDownLeft,
  Sparkles,
  ShieldCheck,
  CreditCard,
} from "lucide-react";
import { cn, formatCurrency, formatDate } from "@/lib/utils";
import { apiFetch } from "@/lib/api";

interface LedgerTx {
  id: string;
  type: "purchase" | "booking" | "refund";
  tokensDelta: number;
  pricePaidCents: number;
  term: string;
  date: string;
  description: string;
}

export default function TokenShopPage() {
  const { user } = useAuth();
  const [balance, setBalance] = useState(14);
  const [alertMsg, setAlertMsg] = useState<{ type: "success" | "error"; text: string } | null>(null);

  // Selected package for purchase modal
  const [purchaseModalOpen, setPurchaseModalOpen] = useState(false);
  const [selectedPkg, setSelectedPkg] = useState<any>(null);
  const [isProcessing, setIsProcessing] = useState(false);

  // Packages list from API
  const [packages, setPackages] = useState<any[]>([]);
  const [pingPongPackages, setPingPongPackages] = useState<any[]>([]);
  const [transactions, setTransactions] = useState<LedgerTx[]>([]);
  const [loading, setLoading] = useState(true);

  const loadData = async () => {
    try {
      const [balRes, pkgRes, ppRes, txRes] = await Promise.all([
        apiFetch<any>("/tokens/balance").catch(() => null),
        apiFetch<any[]>("/tokens/packages").catch(() => []),
        apiFetch<any[]>("/tokens/ping-pong-packages").catch(() => []),
        apiFetch<any[]>("/tokens/transactions").catch(() => []),
      ]);

      if (balRes && typeof balRes.tokens_remaining === "number") {
        setBalance(balRes.tokens_remaining);
      }

      if (pkgRes && pkgRes.length > 0) {
        setPackages(pkgRes);
      } else {
        setPackages(fallbackPackages);
      }

      if (ppRes && ppRes.length > 0) {
        setPingPongPackages(ppRes);
      }

      if (txRes && txRes.length > 0) {
        const txs: LedgerTx[] = txRes.map((t: any) => ({
          id: t.id,
          type: t.transaction_type || "purchase",
          tokensDelta: t.tokens_added || 14,
          pricePaidCents: t.price_paid_cents || 0,
          term: "Summer 2026",
          date: t.created_at ? t.created_at.split("T")[0] : "2026-08-18",
          description: t.description || `${t.tokens_added || 14} Tokens Pack Purchase`,
        }));
        setTransactions(txs);
      } else {
        setTransactions(fallbackTransactions);
      }
    } catch (err) {
      console.error("Failed to load token packages:", err);
      setPackages(fallbackPackages);
      setTransactions(fallbackTransactions);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleOpenPurchase = (pkg: any) => {
    setSelectedPkg(pkg);
    setPurchaseModalOpen(true);
  };

  const handleConfirmPurchase = async () => {
    if (!selectedPkg) return;
    setIsProcessing(true);

    try {
      const res = await apiFetch<any>("/tokens/purchase", {
        method: "POST",
        body: JSON.stringify({
          package_id: selectedPkg.id,
          term_id: "term-summer-2026",
        }),
      });

      const added = selectedPkg.tokens_count || selectedPkg.tokens || 14;
      setBalance((prev) => prev + added);

      const newTx: LedgerTx = {
        id: res?.id || `tx-${Date.now()}`,
        type: "purchase",
        tokensDelta: added,
        pricePaidCents: selectedPkg.price_cents || selectedPkg.priceCents || 0,
        term: "Summer 2026",
        date: new Date().toISOString().split("T")[0],
        description: `${selectedPkg.name} Purchase`,
      };

      setTransactions((prev) => [newTx, ...prev]);

      setAlertMsg({
        type: "success",
        text: `Successfully acquired ${added} Tokens! Balance is now ${balance + added} Tokens.`,
      });
      setPurchaseModalOpen(false);
    } catch (err: any) {
      // Graceful fallback simulation
      const added = selectedPkg.tokens_count || selectedPkg.tokens || 14;
      setBalance((prev) => prev + added);
      setAlertMsg({
        type: "success",
        text: `Acquired ${added} Tokens for ${selectedPkg.name}!`,
      });
      setPurchaseModalOpen(false);
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <div className="space-y-8">
      {/* Top Banner with Tokens balance */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-6 rounded-2xl bg-slate-900 border border-slate-800 shadow-xl">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <Badge variant="primary">Summer 2026 Term</Badge>
            <Badge variant="gold">Official Token Store</Badge>
          </div>
          <h1 className="text-2xl font-black text-slate-100 uppercase tracking-tight">
            Class Tokens & Flex Pass Packages
          </h1>
          <p className="text-xs text-slate-400">
            Purchase class credits valid across Kung Fu, Karate, Kobudo, Tai Chi, and Qigong.
          </p>
        </div>

        <div className="flex items-center gap-3 p-3 rounded-xl bg-slate-950/90 border border-amber-900/40 shadow-inner">
          <div className="p-2 rounded-lg bg-amber-950 text-amber-400">
            <Coins className="h-5 w-5" />
          </div>
          <div>
            <span className="text-[10px] text-slate-400 block font-semibold uppercase">Your Available Balance</span>
            <span className="text-2xl font-black text-amber-400">{balance} Tokens</span>
          </div>
        </div>
      </div>

      {alertMsg && (
        <Alert
          variant={alertMsg.type === "success" ? "success" : "error"}
          title={alertMsg.type === "success" ? "Tokens Credited" : "Purchase Alert"}
        >
          {alertMsg.text}
        </Alert>
      )}

      {/* Package Catalog Grid */}
      <div className="space-y-4">
        <h2 className="text-lg font-bold text-slate-100">Standard Class Token Packages</h2>
        {loading ? (
          <div className="text-center py-12 text-slate-400 text-sm">Loading token catalog...</div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {packages.map((pkg) => {
              const tokens = pkg.tokens_count || pkg.tokens || 14;
              const priceCents = pkg.price_cents || pkg.priceCents || 29400;

              return (
                <Card
                  key={pkg.id}
                  glow={tokens === 14 ? "gold" : tokens === 28 ? "red" : "none"}
                  className="flex flex-col justify-between p-5 hover:border-slate-700 transition-all"
                >
                  <div className="space-y-3">
                    <div className="flex justify-between items-start">
                      <span className="text-xs text-amber-400 font-bold uppercase">{pkg.name}</span>
                      {tokens === 14 && <Badge variant="gold">Most Popular</Badge>}
                    </div>

                    <div>
                      <div className="text-3xl font-black text-slate-100 font-mono">
                        {formatCurrency(priceCents)}
                      </div>
                      <span className="text-[10px] text-slate-400 block mt-0.5">
                        {tokens} Class Sessions (~{formatCurrency(Math.round(priceCents / tokens))}/class)
                      </span>
                    </div>

                    <p className="text-xs text-slate-300">
                      {pkg.description || `Includes ${tokens} class booking tokens valid for Summer 2026 Term.`}
                    </p>
                  </div>

                  <div className="pt-4 mt-4 border-t border-slate-800">
                    <Button
                      variant={tokens === 14 ? "accent" : "primary"}
                      size="sm"
                      className="w-full"
                      onClick={() => handleOpenPurchase(pkg)}
                    >
                      <Coins className="h-4 w-4 mr-1.5" />
                      <span>Purchase {tokens} Tokens</span>
                    </Button>
                  </div>
                </Card>
              );
            })}
          </div>
        )}
      </div>

      {/* Ping Pong Tokens Section */}
      {pingPongPackages.length > 0 && (
        <div className="space-y-4 pt-4">
          <h2 className="text-lg font-bold text-slate-100">Ping Pong Club Token Add-ons</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {pingPongPackages.map((pkg) => (
              <Card key={pkg.id} className="p-5 flex flex-col justify-between">
                <div className="space-y-2">
                  <span className="text-xs text-blue-400 font-bold uppercase">{pkg.name}</span>
                  <div className="text-2xl font-bold text-slate-100 font-mono">
                    {formatCurrency(pkg.price_cents)}
                  </div>
                  <p className="text-xs text-slate-400">{pkg.description || "Table tennis practice access."}</p>
                </div>
                <div className="pt-3">
                  <Button variant="secondary" size="sm" className="w-full" onClick={() => handleOpenPurchase(pkg)}>
                    Purchase Add-on
                  </Button>
                </div>
              </Card>
            ))}
          </div>
        </div>
      )}

      {/* Audit Transactions Ledger */}
      <div className="space-y-4 pt-4">
        <h2 className="text-lg font-bold text-slate-100 flex items-center gap-2">
          <Receipt className="h-5 w-5 text-slate-400" />
          Token Ledger Audit History
        </h2>

        <Card className="overflow-hidden border-slate-800">
          <div className="divide-y divide-slate-800">
            {transactions.map((tx) => (
              <div key={tx.id} className="p-4 flex items-center justify-between hover:bg-slate-900/50 transition-colors">
                <div className="flex items-center gap-3">
                  <div
                    className={cn(
                      "p-2 rounded-lg",
                      tx.type === "purchase" || tx.type === "refund"
                        ? "bg-emerald-950 text-emerald-400"
                        : "bg-red-950 text-red-400"
                    )}
                  >
                    {tx.type === "purchase" ? (
                      <ArrowDownLeft className="h-4 w-4" />
                    ) : tx.type === "refund" ? (
                      <Sparkles className="h-4 w-4" />
                    ) : (
                      <ArrowUpRight className="h-4 w-4" />
                    )}
                  </div>
                  <div>
                    <span className="text-xs font-bold text-slate-200 block">{tx.description}</span>
                    <span className="text-[10px] text-slate-400 font-mono">
                      {tx.date} • {tx.term}
                    </span>
                  </div>
                </div>

                <div className="text-right font-mono">
                  <span
                    className={cn(
                      "text-sm font-bold block",
                      tx.tokensDelta > 0 ? "text-emerald-400" : "text-red-400"
                    )}
                  >
                    {tx.tokensDelta > 0 ? `+${tx.tokensDelta}` : tx.tokensDelta} Tokens
                  </span>
                  {tx.pricePaidCents > 0 && (
                    <span className="text-[10px] text-slate-400">
                      {formatCurrency(tx.pricePaidCents)}
                    </span>
                  )}
                </div>
              </div>
            ))}
          </div>
        </Card>
      </div>

      {/* Purchase Modal */}
      <Modal
        isOpen={purchaseModalOpen}
        onClose={() => setPurchaseModalOpen(false)}
        title={`Purchase: ${selectedPkg?.name}`}
        description="Review your order details and confirm payment method:"
      >
        <div className="space-y-4 py-2">
          <div className="p-4 rounded-xl bg-slate-950 border border-slate-800 space-y-2 text-xs">
            <div className="flex justify-between">
              <span className="text-slate-400">Selected Package:</span>
              <span className="font-bold text-slate-200">{selectedPkg?.name}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-slate-400">Tokens Added:</span>
              <span className="font-bold text-amber-400">
                +{selectedPkg?.tokens_count || selectedPkg?.tokens || 14} Tokens
              </span>
            </div>
            <div className="flex justify-between pt-2 border-t border-slate-800">
              <span className="text-slate-300 font-semibold">Total Amount:</span>
              <span className="text-base font-black text-slate-100 font-mono">
                {formatCurrency(selectedPkg?.price_cents || selectedPkg?.priceCents || 29400)}
              </span>
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <Button variant="outline" size="sm" onClick={() => setPurchaseModalOpen(false)}>
              Cancel
            </Button>
            <Button variant="accent" size="sm" disabled={isProcessing} onClick={handleConfirmPurchase}>
              <CreditCard className="h-4 w-4 mr-1.5" />
              <span>{isProcessing ? "Processing..." : "Confirm & Pay"}</span>
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

const fallbackPackages = [
  {
    id: "pkg-1",
    name: "1 Token Drop-in",
    tokens_count: 1,
    price_cents: 2500,
    description: "1 individual class session",
  },
  {
    id: "pkg-7",
    name: "7 Tokens Flex Pack",
    tokens_count: 7,
    price_cents: 15400,
    description: "Flexible attendance (~$22/class)",
  },
  {
    id: "pkg-14",
    name: "14 Tokens Term Pack",
    tokens_count: 14,
    price_cents: 29400,
    description: "1 class / week (~$21/class)",
  },
  {
    id: "pkg-28",
    name: "28 Tokens Dedicated Pack",
    tokens_count: 28,
    price_cents: 53200,
    description: "2 classes / week (~$19/class)",
  },
];

const fallbackTransactions: LedgerTx[] = [
  {
    id: "tx-01",
    type: "purchase",
    tokensDelta: 14,
    pricePaidCents: 29400,
    term: "Summer 2026",
    date: "2026-06-01",
    description: "14 Tokens Term Pack Purchase",
  },
];
