import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of, catchError } from 'rxjs';
import { User, MealHistory, MenuReview } from '../model/user';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private http = inject(HttpClient);
  private baseUrl = 'http://localhost:8002/api/users';
  private baseUrlCanteen = 'http://localhost:8001/api/canteens';
  private baseUrlMenu = 'http://localhost:8001/api/menus';

  getUserById(id: string): Observable<User> {
    return this.http.get<User>(`${this.baseUrl}/${id}`);
  }

  getMealHistory(userId: string): Observable<MealHistory[]> {
    return this.http.get<MealHistory[]>(`${this.baseUrlCanteen}/meal-history/${userId}`);
  }

  // Jednostavnija metoda - koristi novi backend endpoint
  getMealHistoryWithReviews(userId: string): Observable<MealHistory[]> {
    return this.http.get<MealHistory[]>(`${this.baseUrlMenu}/reviews/${userId}`);
  }

  // Review-related methods
  createReview(review: MenuReview): Observable<MenuReview> {
    return this.http.post<MenuReview>(`${this.baseUrlMenu}/reviews/`, review);
  }

  updateReview(review: MenuReview): Observable<MenuReview> {
    return this.http.put<MenuReview>(`${this.baseUrlMenu}/reviews/`, review);
  }

  getReview(reviewId: string): Observable<MenuReview> {
    return this.http.get<MenuReview>(`${this.baseUrlMenu}/review/{reviewId}`);
  }

}
