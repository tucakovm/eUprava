export interface User {
  id: string;           // UUID
  firstname: string;
  lastname: string;
  username: string;
  email: string;
  is_active: boolean;
  role: string;
}

export interface MealHistory {
  id: string;
  menu_id: string;
  menu_name: string;
  selected_at: string; 
}

