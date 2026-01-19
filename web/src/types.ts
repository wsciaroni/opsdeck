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
