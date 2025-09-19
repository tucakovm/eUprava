import {inject, Injectable} from '@angular/core';
import {CanteenDto} from './canteen.service';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Menu} from '../model/menus';
import {Observable} from 'rxjs';
import {AuthService} from './auth.service';

export interface MenuWithCard {
  menu: Menu;
  card?: { id: string; stanje: number; studentID: string };
}

@Injectable({
  providedIn: 'root'
})
export class MenuService {

  private http = inject(HttpClient);
  private authService = inject(AuthService);
  private baseUrl = 'http://localhost:8001/api/menus/';
  private baseUrl2 = 'http://localhost:8001/api/menu/';

  create(menu: Menu) {
    return this.http.post<Menu>(`${this.baseUrl}`, menu);
  }

  getAll(canteenId: string): Observable<Menu[]> {
    return this.http.get<Menu[]>(`${this.baseUrl}/${canteenId}`);
  }

  delete(id: string) {
    return this.http.delete(`${this.baseUrl}${id}`);
  }

  getTopMeals(): Observable<{ menu_name: string, score: number }[]> {
    return this.http.get<{ menu_name: string, score: number }[]>(`${this.baseUrl}top-rated/`);
  }

  getMenu(menuId: string, userId: string): Observable<MenuWithCard> {
    const headers = new HttpHeaders({ 'X-Student-ID': userId });
    return this.http.get<MenuWithCard>(`${this.baseUrl2}${menuId}`, { headers });
  }

  takeMeal(payload: { studentId: string | null; delta: number; menuId: any, studentUsername: string | null}) {
    return this.http.post("http://localhost:8001/api/meal/", payload);
  }

}
