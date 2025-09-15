export interface Canteen {
  id: string;
  name: string;
  address: string;
  open_at: Date;
  close_at: Date;
}

export interface CanteenDTO {
  id: string;
  name: string;
  address: string;
  open_at: string;
  close_at: string;
}
