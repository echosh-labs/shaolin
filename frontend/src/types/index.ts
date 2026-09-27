export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  date_of_birth: string;
  role: "student" | "instructor" | "admin";
  current_rank: string;
  created_at: string;
}

export interface StudentWeeklyStats {
  user_id: string;
  year_week: string;
  tokens_used: number;
  kungfu_attended: number;
  taichi_attended: number;
  qigong_attended: number;
  total_attended: number;
  weekly_points: number;
}

export interface StudentOverallStats {
  user_id: string;
  total_tokens_used: number;
  kungfu_attended: number;
  taichi_attended: number;
  qigong_attended: number;
  total_attended: number;
  overall_points: number;
  level_tier: string;
  last_updated: string;
}

export interface Term {
  id: string;
  name: string;
  start_date: string;
  end_date: string;
  is_active: boolean;
}

export interface TermBreak {
  id: string;
  term_id: string;
  start_date: string;
  end_date: string;
  notes: string;
}

export interface EventType {
  id: number;
  name: string;
  bg_color: string;
  fg_color: string;
}

export interface Hall {
  id: number;
  name: string;
  location_id?: string;
}

export interface ClassItem {
  id: string;
  term_id: string;
  event_type_id: number;
  name: string;
  instructor_id?: string;
  day_of_week: number;
  start_time: string;
  end_time: string;
  capacity: number;
  zoom_link?: string;
  hall_ids?: number[];
}

export interface ClassOccurrence {
  id: string;
  class_id: string;
  date: string;
  class_week?: number;
  booked_count: number;
  status: "scheduled" | "cancelled" | "completed";
  notes?: string;
  image_url?: string;
  zoom_link?: string;
}

export interface Booking {
  id: string;
  user_id: string;
  occurrence_id: string;
  status: "confirmed" | "cancelled" | "waitlisted";
  attendance_mode: "in_person" | "live_stream";
  attendance_status: "confirmed" | "attended" | "no_show" | "excused";
  booked_at: string;
  cancelled_at?: string;
}

export interface TokenPackage {
  id: string;
  name: string;
  tokens_count: number;
  adult_price_cents: number;
  youth_price_cents: number;
  is_unlimited: boolean;
  terms_span?: number;
}

export interface UserTermTokens {
  user_id: string;
  term_id: string;
  tokens_remaining: number;
}

export interface TokenTransaction {
  id: string;
  user_id: string;
  term_id: string;
  tokens_added: number;
  price_paid_cents: number;
  transaction_type: string;
  created_at: string;
}

export interface StoreProduct {
  id: string;
  category_id: string;
  name: string;
  description: string;
  price_cents: number;
  image_url: string;
  inventory_count: number;
  is_active: boolean;
}

export interface CartItem {
  id: string;
  user_id: string;
  product_id: string;
  quantity: number;
  product?: StoreProduct;
}

export interface GradingTrack {
  id: string;
  name: string;
  description: string;
}

export interface ForumBoard {
  id: string;
  title: string;
  description: string;
  category: string;
  allowed_post_roles: string;
  display_order: number;
}

export interface ForumTopic {
  id: string;
  board_id: string;
  user_id: string;
  title: string;
  is_pinned: boolean;
  is_locked: boolean;
  post_count: number;
  last_post_at: string;
  created_at: string;
}
