export interface Meal {
  id: string;          // UUID u string formatu
  name: string;
  description: string;
  price: number;
}

export interface MealDTO {
  name: string;
  description: string;
  price: number;
}

export interface TopMenu {
  id: string;
  name: string;
  avg_rating: number;
}

export interface MenuDTO {
  id: string;
  name: string;
  breakfast: Meal;
  lunch: Meal;
  dinner: Meal;
}



export enum Weekday {
  Monday = 'Monday',
  Tuesday = 'Tuesday',
  Wednesday = 'Wednesday',
  Thursday = 'Thursday',
  Friday = 'Friday',
  Saturday = 'Saturday',
  Sunday = 'Sunday'
}


export interface Menu {
  id: string;           // UUID u string formatu
  name: string;
  canteen_id: string;   // UUID u string formatu
  weekday: Weekday;
  breakfast: MealDTO;
  lunch: MealDTO;
  dinner: MealDTO;
}
