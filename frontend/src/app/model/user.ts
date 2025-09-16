export interface User {
  id: string;
  firstname: string;
  lastname: string;
  username: string;
  email: string;
  role: string;
  is_active: boolean;
}

export interface MealHistory {
  id?: string;
  menu_id?: string;
  menu_name: string;
  selected_at: string;
  review?: MenuReview;
}

export interface MenuReview {
  id: string;
  menu_id: string;
  user_id: string;
  breakfast_review: number;
  lunch_review: number;
  dinner_review: number;
}
