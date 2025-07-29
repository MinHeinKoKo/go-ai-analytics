import axios from "axios";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export const api = axios.create({
  baseURL: `${API_BASE_URL}/api`,
  headers: {
    "Content-Type": "application/json",
  },
});

// Add auth token to requests
api.interceptors.request.use(
  (config) => {
    // Only access localStorage on client side
    if (typeof window !== "undefined") {
      const token = localStorage.getItem("auth_token");
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token is invalid or expired
      if (typeof window !== "undefined") {
        localStorage.removeItem("auth_token");
        // Redirect to login page
        window.location.href = "/login";
      }
    }
    return Promise.reject(error);
  }
);

// Types
export interface Customer {
  id: string;
  customer_id: string;
  age: number;
  gender: string;
  location: string;
  income_range: string;
  registration_date: string;
  last_purchase_date?: string;
  total_spent: number;
  purchase_frequency: number;
  preferred_category: string;
  created_at: string;
  updated_at: string;
}

export interface Purchase {
  id: string;
  customer_id: string;
  product_id: string;
  category: string;
  amount: number;
  quantity: number;
  purchase_date: string;
  channel: string;
  created_at: string;
}

export interface Campaign {
  id: string;
  campaign_id: string;
  name: string;
  type: string;
  target_segment: string;
  budget: number;
  start_date: string;
  end_date: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface CampaignPerformance {
  id: string;
  campaign_id: string;
  impressions: number;
  clicks: number;
  conversions: number;
  revenue: number;
  cost: number;
  ctr: number;
  cpc: number;
  roas: number;
  date: string;
  created_at: string;
}

export interface CustomerSegment {
  id: string;
  segment_id: string;
  name: string;
  description: string;
  criteria: Record<string, any>;
  size: number;
  created_at: string;
  updated_at: string;
}

export interface PredictionResult {
  id: string;
  customer_id: string;
  prediction_type: string;
  probability: number;
  value: number;
  confidence: number;
  created_at: string;
}

export interface DashboardData {
  total_customers: number;
  total_purchases: number;
  total_revenue: number;
  avg_order_value: number;
  total_campaigns: number;
  active_campaigns: number;
}

// API functions
export const analyticsApi = {
  // Dashboard
  getDashboard: (startDate?: string, endDate?: string) => {
    const params = new URLSearchParams();
    if (startDate) params.append("start_date", startDate);
    if (endDate) params.append("end_date", endDate);
    return api.get<{ dashboard: DashboardData }>(
      `/analytics/dashboard?${params}`
    );
  },

  // Customers
  getCustomers: (limit = 50, offset = 0) =>
    api.get<{ customers: Customer[] }>(
      `/customers?limit=${limit}&offset=${offset}`
    ),

  createCustomer: (
    customer: Omit<Customer, "id" | "created_at" | "updated_at">
  ) => api.post<{ customer: Customer }>("/customers", customer),

  // Purchases
  createPurchase: (purchase: Omit<Purchase, "id" | "created_at">) =>
    api.post<{ purchase: Purchase }>("/purchases", purchase),

  // Campaigns
  getCampaigns: () => api.get<{ campaigns: Campaign[] }>("/campaigns"),

  createCampaign: (
    campaign: Omit<Campaign, "id" | "created_at" | "updated_at">
  ) => api.post<{ campaign: Campaign }>("/campaigns", campaign),

  createCampaignPerformance: (
    performance: Omit<
      CampaignPerformance,
      "id" | "created_at" | "ctr" | "cpc" | "roas"
    >
  ) =>
    api.post<{ performance: CampaignPerformance }>(
      "/campaigns/performance",
      performance
    ),

  // AI Analytics
  performSegmentation: (request: {
    algorithm: string;
    features: string[];
    parameters?: Record<string, any>;
  }) =>
    api.post<{ segments: CustomerSegment[] }>(
      "/analytics/segmentation",
      request
    ),

  predictCustomerBehavior: (request: {
    customer_id: string;
    prediction_type: string;
  }) =>
    api.post<{ prediction: PredictionResult }>(
      "/analytics/prediction",
      request
    ),

  optimizeCampaign: (request: {
    campaign_id: string;
    objective: string;
    parameters?: Record<string, any>;
  }) => api.post<{ optimization: any }>("/analytics/optimization", request),

  // Sample data
  generateSampleData: () =>
    api.post<{ sample_data: any }>("/analytics/sample-data"),

  importTrainingData: (data: {
    customers?: Customer[];
    purchases?: Purchase[];
    campaigns?: Campaign[];
    performance?: CampaignPerformance[];
  }) => api.post<{ import_results: any }>("/analytics/import", data),
};

// Auth API
export const authApi = {
  login: (email: string, password: string) =>
    api.post<{ token: string; user: any }>("/auth/login", { email, password }),

  register: (email: string, password: string, name: string) =>
    api.post<{ token: string; user: any }>("/auth/register", {
      email,
      password,
      name,
    }),

  getMe: () => api.get<{ user: any }>("/auth/me"),
};
