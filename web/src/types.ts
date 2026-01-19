export interface User {
  id: string;
  email: string;
  name: string;
  role: string;
  avatar_url: string;
  created_at: string;
  updated_at: string;
}

export interface Organization {
  id: string;
  name: string;
  slug: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface Ticket {
  id: string;
  organization_id: string;
  title: string;
  description: string;
  location: string;
  status_id: string;
  priority_id: string;
  reporter_id: string;
  assignee_user_id: string | null;
  created_at: string;
  updated_at: string;
  completed_at: string | null;
}

export interface CreateTicketRequest {
  title: string;
  description: string;
  priority_id: string;
  organization_id: string;
  location?: string;
}
