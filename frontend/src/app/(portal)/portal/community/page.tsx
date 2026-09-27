"use client";

import React, { useState, useEffect } from "react";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/Button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/Card";
import { Badge } from "@/components/ui/Badge";
import { Modal } from "@/components/ui/Modal";
import { Input, Label } from "@/components/ui/Input";
import { Alert } from "@/components/ui/Alert";
import {
  MessagesSquare,
  Pin,
  MessageCircle,
  Video,
  PlusCircle,
  Sparkles,
  Play,
  User,
  Clock,
  ArrowRight,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { apiFetch } from "@/lib/api";

interface TopicItem {
  id: string;
  board_id: string;
  boardName: string;
  title: string;
  author: string;
  role: string;
  isPinned: boolean;
  repliesCount: number;
  lastActivity: string;
  content: string;
}

interface MediaItem {
  id: string;
  title: string;
  description: string;
  asset_type: string;
  url: string;
  thumbnail_url: string;
  distribution_type: string;
}

export default function CommunityPage() {
  const { user } = useAuth();
  const [activeTab, setActiveTab] = useState<"forums" | "videos">("forums");
  const [boards, setBoards] = useState<any[]>([]);
  const [selectedBoardId, setSelectedBoardId] = useState<string>("all");

  // New Thread Modal
  const [newThreadOpen, setNewThreadOpen] = useState(false);
  const [threadBoardId, setThreadBoardId] = useState<string>("");
  const [threadTitle, setThreadTitle] = useState("");
  const [threadContent, setThreadContent] = useState("");
  const [toastMsg, setToastMsg] = useState<string | null>(null);

  // Active Topic View Modal
  const [activeTopic, setActiveTopic] = useState<TopicItem | null>(null);
  const [posts, setPosts] = useState<any[]>([]);
  const [replyText, setReplyText] = useState("");

  const [topics, setTopics] = useState<TopicItem[]>([]);
  const [mediaAssets, setMediaAssets] = useState<MediaItem[]>([]);
  const [loading, setLoading] = useState(true);

  const loadData = async () => {
    try {
      const [boardsRes, mediaRes] = await Promise.all([
        apiFetch<any[]>("/forums/boards").catch(() => []),
        apiFetch<any[]>("/media").catch(() => []),
      ]);

      if (boardsRes && boardsRes.length > 0) {
        setBoards(boardsRes);
        setThreadBoardId(boardsRes[1]?.id || boardsRes[0]?.id);

        // Fetch topics across boards
        const allTopics: TopicItem[] = [];
        for (const b of boardsRes) {
          try {
            const tList = await apiFetch<any[]>(`/forums/topics/${b.id}`);
            if (tList && tList.length > 0) {
              tList.forEach((t: any) => {
                allTopics.push({
                  id: t.id,
                  board_id: b.id,
                  boardName: b.name,
                  title: t.title,
                  author: "Disciple",
                  role: "student",
                  isPinned: t.is_pinned === 1,
                  repliesCount: 3,
                  lastActivity: "Recent",
                  content: t.content,
                });
              });
            }
          } catch {}
        }

        if (allTopics.length > 0) {
          setTopics(allTopics);
        } else {
          setTopics(fallbackTopics);
        }
      } else {
        setTopics(fallbackTopics);
      }

      if (mediaRes && mediaRes.length > 0) {
        setMediaAssets(mediaRes);
      } else {
        setMediaAssets(fallbackMedia);
      }
    } catch (err) {
      console.error("Failed to load forum/media data:", err);
      setTopics(fallbackTopics);
      setMediaAssets(fallbackMedia);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const filteredTopics = topics.filter(
    (t) => selectedBoardId === "all" || t.board_id === selectedBoardId
  );

  const handleCreateTopic = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!threadTitle.trim() || !threadContent.trim()) return;

    try {
      const targetBoard = threadBoardId || boards[0]?.id || "fb-02";
      const res = await apiFetch<any>("/forums/topics", {
        method: "POST",
        body: JSON.stringify({
          board_id: targetBoard,
          title: threadTitle,
          content: threadContent,
        }),
      });

      const boardObj = boards.find((b) => b.id === targetBoard);
      const newTop: TopicItem = {
        id: res?.id || `top-${Date.now()}`,
        board_id: targetBoard,
        boardName: boardObj?.name || "Community Discussion",
        title: threadTitle,
        author: user?.first_name || "You",
        role: user?.role || "student",
        isPinned: false,
        repliesCount: 0,
        lastActivity: "Just now",
        content: threadContent,
      };

      setTopics((prev) => [newTop, ...prev]);
      setThreadTitle("");
      setThreadContent("");
      setNewThreadOpen(false);

      setToastMsg("Topic posted successfully!");
      setTimeout(() => setToastMsg(null), 3000);
    } catch (err: any) {
      setToastMsg(`Error: ${err.message || "Failed to post topic"}`);
    }
  };

  const handleOpenTopic = async (topic: TopicItem) => {
    setActiveTopic(topic);
    try {
      const postsRes = await apiFetch<any[]>(`/forums/topics/${topic.id}/posts`);
      setPosts(postsRes || []);
    } catch {
      setPosts([]);
    }
  };

  const handleCreatePost = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!activeTopic || !replyText.trim()) return;

    try {
      const res = await apiFetch<any>("/forums/posts", {
        method: "POST",
        body: JSON.stringify({
          topic_id: activeTopic.id,
          content: replyText,
        }),
      });

      setPosts((prev) => [
        ...prev,
        {
          id: res?.id || `post-${Date.now()}`,
          content: replyText,
          author_id: user?.id,
          created_at: new Date().toISOString(),
        },
      ]);
      setReplyText("");
    } catch (err: any) {
      setToastMsg(`Failed to reply: ${err.message}`);
    }
  };

  return (
    <div className="space-y-8">
      {/* Toast Alert */}
      {toastMsg && (
        <div className="fixed bottom-6 right-6 z-50 bg-slate-900 border border-red-700 text-slate-100 px-4 py-3 rounded-xl shadow-2xl animate-in fade-in">
          <span className="text-xs font-semibold">{toastMsg}</span>
        </div>
      )}

      {/* Top Banner */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-6 rounded-2xl bg-slate-900 border border-slate-800 shadow-xl">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <Badge variant="primary">Disciple Community</Badge>
            <Badge variant="gold">Forums & Digital Vault</Badge>
          </div>
          <h1 className="text-2xl font-black text-slate-100 uppercase tracking-tight">
            Martial Arts Message Board & Media
          </h1>
          <p className="text-xs text-slate-400">
            Connect with instructors, ask stance questions, and stream recorded curriculum video tutorials.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <Button variant="primary" size="sm" onClick={() => setNewThreadOpen(true)}>
            <PlusCircle className="h-4 w-4 mr-1.5" />
            <span>New Discussion Topic</span>
          </Button>
        </div>
      </div>

      {/* Tab Navigation */}
      <div className="flex items-center gap-2 border-b border-slate-800 pb-3">
        <button
          onClick={() => setActiveTab("forums")}
          className={cn(
            "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all cursor-pointer flex items-center gap-2",
            activeTab === "forums"
              ? "bg-red-600 text-white shadow-md shadow-red-950/50"
              : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
          )}
        >
          <MessagesSquare className="h-4 w-4" />
          Message Boards
        </button>

        <button
          onClick={() => setActiveTab("videos")}
          className={cn(
            "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all cursor-pointer flex items-center gap-2",
            activeTab === "videos"
              ? "bg-blue-600 text-white shadow-md shadow-blue-950/50"
              : "bg-slate-900 text-slate-400 hover:text-slate-100 border border-slate-800"
          )}
        >
          <Video className="h-4 w-4" />
          Digital Video Vault ({mediaAssets.length})
        </button>
      </div>

      {/* TAB 1: Forums */}
      {activeTab === "forums" && (
        <div className="space-y-6">
          {/* Board Categories */}
          <div className="flex flex-wrap items-center gap-2">
            <button
              onClick={() => setSelectedBoardId("all")}
              className={cn(
                "px-3 py-1.5 rounded-xl text-xs font-bold transition-all cursor-pointer",
                selectedBoardId === "all"
                  ? "bg-slate-100 text-slate-950"
                  : "bg-slate-900 text-slate-400 border border-slate-800 hover:text-slate-200"
              )}
            >
              All Boards
            </button>
            {boards.map((b) => (
              <button
                key={b.id}
                onClick={() => setSelectedBoardId(b.id)}
                className={cn(
                  "px-3 py-1.5 rounded-xl text-xs font-bold transition-all cursor-pointer",
                  selectedBoardId === b.id
                    ? "bg-red-600 text-white"
                    : "bg-slate-900 text-slate-400 border border-slate-800 hover:text-slate-200"
                )}
              >
                {b.name}
              </button>
            ))}
          </div>

          {/* Topics List */}
          <div className="space-y-3">
            {filteredTopics.map((top) => (
              <Card
                key={top.id}
                glow={top.isPinned ? "red" : "none"}
                className="p-5 hover:border-slate-700 transition-all cursor-pointer"
                onClick={() => handleOpenTopic(top)}
              >
                <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2">
                  <div className="space-y-1.5">
                    <div className="flex items-center gap-2">
                      {top.isPinned && (
                        <Badge variant="primary" size="sm" className="flex items-center gap-1">
                          <Pin className="h-3 w-3" />
                          PINNED
                        </Badge>
                      )}
                      <Badge variant="outline" size="sm">
                        {top.boardName}
                      </Badge>
                    </div>

                    <h3 className="text-base font-bold text-slate-100 hover:text-amber-400 transition-colors">
                      {top.title}
                    </h3>
                    <p className="text-xs text-slate-400 line-clamp-1">{top.content}</p>
                  </div>

                  <div className="flex items-center gap-3 text-xs text-slate-400 shrink-0 font-mono">
                    <span className="flex items-center gap-1">
                      <MessageCircle className="h-3.5 w-3.5" />
                      {top.repliesCount}
                    </span>
                    <span>{top.lastActivity}</span>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        </div>
      )}

      {/* TAB 2: Video Vault */}
      {activeTab === "videos" && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {mediaAssets.map((v) => (
            <Card key={v.id} glow="gold" className="overflow-hidden flex flex-col justify-between">
              <div className="relative aspect-video bg-slate-950 flex items-center justify-center group cursor-pointer">
                <div className="p-4 rounded-full bg-red-600/90 text-white shadow-xl group-hover:scale-110 transition-transform">
                  <Play className="h-6 w-6 ml-0.5" />
                </div>
                <div className="absolute top-2 right-2">
                  <Badge variant="gold">
                    {v.distribution_type === "reward" ? "Belt Restricted" : "Public Masterclass"}
                  </Badge>
                </div>
              </div>

              <CardHeader className="p-4 space-y-1">
                <CardTitle className="text-base font-bold">{v.title}</CardTitle>
                <CardDescription className="text-xs text-slate-400">{v.description}</CardDescription>
              </CardHeader>

              <CardContent className="p-4 pt-0">
                <Button variant="secondary" size="sm" className="w-full">
                  <Play className="h-3.5 w-3.5 mr-1" />
                  Stream High-Definition Video
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Create Topic Modal */}
      <Modal
        isOpen={newThreadOpen}
        onClose={() => setNewThreadOpen(false)}
        title="Start New Discussion Thread"
        description="Post a question or milestone to the community board:"
      >
        <form onSubmit={handleCreateTopic} className="space-y-4 py-2">
          <div className="space-y-1.5">
            <Label>Select Message Board</Label>
            <select
              value={threadBoardId}
              onChange={(e) => setThreadBoardId(e.target.value)}
              className="w-full bg-slate-950 border border-slate-800 rounded-xl p-3 text-xs text-slate-200 focus:outline-none focus:border-red-600"
            >
              {boards
                .filter((b) => b.allowed_post_roles === "all" || user?.role === "admin")
                .map((b) => (
                  <option key={b.id} value={b.id}>
                    {b.name}
                  </option>
                ))}
            </select>
          </div>

          <div className="space-y-1.5">
            <Label>Topic Title</Label>
            <Input
              value={threadTitle}
              onChange={(e) => setThreadTitle(e.target.value)}
              placeholder="e.g. Question on 18-form silk reeling footwork"
              required
            />
          </div>

          <div className="space-y-1.5">
            <Label>Content & Description</Label>
            <textarea
              value={threadContent}
              onChange={(e) => setThreadContent(e.target.value)}
              rows={4}
              className="w-full bg-slate-950 border border-slate-800 rounded-xl p-3 text-xs text-slate-200 focus:outline-none focus:border-red-600"
              placeholder="Detail your question, observations, or training notes..."
              required
            />
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" size="sm" onClick={() => setNewThreadOpen(false)}>
              Cancel
            </Button>
            <Button type="submit" variant="primary" size="sm">
              Post Topic
            </Button>
          </div>
        </form>
      </Modal>

      {/* Topic Detail & Replies Modal */}
      <Modal
        isOpen={!!activeTopic}
        onClose={() => setActiveTopic(null)}
        title={activeTopic?.title || "Topic Thread"}
        description={`Posted in ${activeTopic?.boardName}`}
      >
        <div className="space-y-4 py-2 max-h-[65vh] overflow-y-auto pr-2">
          {/* Main post */}
          <div className="p-4 rounded-xl bg-slate-950 border border-slate-800 space-y-2">
            <div className="flex justify-between items-center text-[11px] text-slate-400">
              <span className="font-bold text-slate-200">Author: {activeTopic?.author}</span>
              <Badge variant="primary" size="sm">Original Post</Badge>
            </div>
            <p className="text-xs text-slate-300 leading-relaxed whitespace-pre-wrap">
              {activeTopic?.content}
            </p>
          </div>

          {/* Replies */}
          <div className="space-y-2">
            <h4 className="text-xs font-bold text-slate-300 uppercase tracking-wider">
              Replies & Discussion ({posts.length})
            </h4>
            {posts.length === 0 ? (
              <p className="text-xs text-slate-500 italic">No replies yet. Be the first to chime in!</p>
            ) : (
              posts.map((p, idx) => (
                <div key={idx} className="p-3 rounded-xl bg-slate-900 border border-slate-800 space-y-1">
                  <div className="flex justify-between text-[10px] text-slate-400 font-mono">
                    <span>Reply #{idx + 1}</span>
                    <span>{p.created_at ? p.created_at.split("T")[0] : "Recent"}</span>
                  </div>
                  <p className="text-xs text-slate-300 whitespace-pre-wrap">{p.content}</p>
                </div>
              ))
            )}
          </div>

          {/* Reply Form */}
          <form onSubmit={handleCreatePost} className="space-y-2 pt-2 border-t border-slate-800">
            <textarea
              value={replyText}
              onChange={(e) => setReplyText(e.target.value)}
              rows={3}
              placeholder="Write a respectful reply or instructor advice..."
              className="w-full bg-slate-950 border border-slate-800 rounded-xl p-3 text-xs text-slate-200 focus:outline-none focus:border-red-600"
              required
            />
            <div className="flex justify-end gap-2">
              <Button type="button" variant="outline" size="sm" onClick={() => setActiveTopic(null)}>
                Close
              </Button>
              <Button type="submit" variant="primary" size="sm">
                Post Reply
              </Button>
            </div>
          </form>
        </div>
      </Modal>
    </div>
  );
}

const fallbackTopics: TopicItem[] = [
  {
    id: "top-01",
    board_id: "fb-01",
    boardName: "Academy Announcements & Schedules",
    title: "Welcome to Summer 2026 Term & Belt Exam Schedule",
    author: "Head Master Marcus Vance",
    role: "admin",
    isPinned: true,
    repliesCount: 5,
    lastActivity: "Recent",
    content: "Welcome students! The Summer 2026 Term is now officially active.",
  },
  {
    id: "top-02",
    board_id: "fb-02",
    boardName: "Martial Arts & Stance Technique Q&A",
    title: "Tips for reaching 2-minute Mabu horse stance?",
    author: "Justin",
    role: "student",
    isPinned: false,
    repliesCount: 2,
    lastActivity: "Recent",
    content: "Preparing for the Level 2 exam and struggling past 90 seconds in Ma Bu. Any advice?",
  },
];

const fallbackMedia: MediaItem[] = [
  {
    id: "med-01",
    title: "5 Core Stances Masterclass",
    description: "Detailed breakdown of Ma Bu, Gong Bu, Pu Bu, Xu Bu, and Xie Bu fundamentals.",
    asset_type: "video",
    url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
    thumbnail_url: "/images/media/stances-thumb.jpg",
    distribution_type: "public",
  },
  {
    id: "med-02",
    title: "Xiao Hong Quan (Small Flood Fist) Step-by-Step",
    description: "Section-by-section breakdown of the fundamental fist form.",
    asset_type: "video",
    url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4",
    thumbnail_url: "/images/media/xiao-hong-quan-thumb.jpg",
    distribution_type: "reward",
  },
];
