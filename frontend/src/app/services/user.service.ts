import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { User } from '../model/user';
import {MealHistory} from '../model/user';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private http = inject(HttpClient);
  private baseUrl = 'http://localhost:8002/api/users';
  private baseUrlCanteen = 'http://localhost:8001/api/canteens';

  getUserById(id: string): Observable<User> {
    return this.http.get<User>(`${this.baseUrl}/${id}`);
  }

  getMealHistory(userId: string) {
    return this.http.get<MealHistory[]>(`${this.baseUrlCanteen}/meal-history/${userId}`);
  }

}
